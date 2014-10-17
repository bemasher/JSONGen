// JSONGen - A tool for generating native Golang types from JSON objects.
// Copyright (C) 2014 Douglas Hall
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"sort"
	"strings"
	"unicode"
)

var config Config

type Config struct {
	dumpFilename string

	dumpFile  *os.File
	inputFile *os.File

	titleCase bool
	normalize bool
}

func (c *Config) Parse() (err error) {
	flag.StringVar(&config.dumpFilename, "dump", os.DevNull, "Dump tree structure to file.")
	flag.BoolVar(&config.normalize, "normalize", true, "Squash arrays of struct and determine primitive array type.")
	flag.BoolVar(&config.titleCase, "title", true, "Convert identifiers to title case, treating '_' and '-' as word boundaries.")

	flag.Parse()

	if flag.NArg() == 0 {
		config.inputFile = os.Stdin
	} else {
		config.inputFile, err = os.Open(flag.Arg(0))
		if err != nil {
			return
		}
	}

	c.dumpFile, err = os.Create(c.dumpFilename)
	if err != nil {
		return
	}

	return
}

func (c Config) Close() {
	c.dumpFile.Close()
	c.inputFile.Close()
}

// Field name sanitizer.
type Ident string

// Golang identifiers must begin with a letter and may contain letters, digits
// and _'s. If config.titleCase is true, -, _ and spaces are treated as word
// boundaries, otherwise only spaces are treated as word boundaries.
func (id Ident) String() (s string) {
	s = strings.TrimLeftFunc(string(id), func(r rune) bool {
		return !unicode.IsLetter(r)
	})

	s = strings.Map(func(r rune) rune {
		if r == ' ' {
			return ' '
		}

		if r == '-' || r == '_' {
			if config.titleCase {
				return ' '
			}
			return '_'
		}

		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}

		return -1
	}, s)

	s = strings.Title(s)
	s = strings.Map(func(r rune) rune {
		if r == ' ' {
			return -1
		}
		return r
	}, s)

	if len(s) == 0 {
		s = "_"
	}

	return
}

// Returns a field tag for the original field name.
func (id Ident) Tag() string {
	return "`json:\"" + string(id) + "\"`"
}

// JSON values are translated to go types as follows:
// null   -> interface{}
// bool   -> bool
// int    -> int64
// float  -> float64
// string -> string
// object -> struct
type Type int

const (
	Interface Type = iota + 1
	Bool
	Int
	Float
	String
	Struct
)

func (t Type) String() string {
	switch t {
	case Bool:
		return "bool"
	case Int:
		return "int64"
	case Float:
		return "float64"
	case String:
		return "string"
	case Struct:
		return "struct"
	case Interface:
		return "interface{}"
	}
	return "unset"
}

func (t Type) MarshalText() (text []byte, err error) {
	return []byte(t.String()), nil
}

type Tree struct {
	Name     Ident `json:",omitempty"`
	List     bool  `json:",omitempty"`
	Type     Type
	Children []*Tree `json:",omitempty"`
}

// A tree implements the sort interface on it's children's sanitized names.
func (t Tree) Len() int {
	return len(t.Children)
}

func (t Tree) Less(i, j int) bool {
	return t.Children[i].Name.String() < t.Children[j].Name.String()
}

func (t Tree) Swap(i, j int) {
	t.Children[i], t.Children[j] = t.Children[j], t.Children[i]
}

// Returns canonical golang of the type structure.
func (t *Tree) Format() (formatted []byte, err error) {
	unformatted := []byte("type " + t.formatHelper(0))
	formatted, err = format.Source(unformatted)

	if err != nil {
		formatted = unformatted
	}
	return
}

func (t *Tree) formatHelper(depth int) (r string) {
	indent := strings.Repeat("\t", depth)

	r += indent + t.Name.String() + " "

	defer func() {
		if depth != 0 && string(t.Name) != t.Name.String() {
			r += " " + t.Name.Tag()
		}
		r += "\n"
	}()

	if t.List {
		r += "[]"
	}

	if t.Type == Struct {
		r += "struct {\n"
		defer func() {
			r += indent + "}"
		}()

		for _, child := range t.Children {
			r += child.formatHelper(depth + 1)
		}
	} else {
		r += t.Type.String()
	}

	return
}

// Given a value which JSON has been parsed into, populates the tree.
func (t *Tree) Populate(v interface{}) {
	if v == nil {
		t.Type = Interface
	}

	switch i := v.(type) {
	case bool:
		t.Type = Bool
	case string:
		t.Type = String
	case json.Number:
		if _, err := i.Int64(); err == nil {
			t.Type = Int
		} else {
			if _, err := i.Float64(); err == nil {
				t.Type = Float
			}
		}
	case []interface{}:
		t.List = true
		t.Type = Interface
		for _, v := range i {
			child := &Tree{}
			child.Populate(v)
			t.Children = append(t.Children, child)
		}
	case map[string]interface{}:
		t.Type = Struct
		for k, v := range i {
			child := &Tree{Name: Ident(k)}
			child.Populate(v)
			t.Children = append(t.Children, child)
		}
		sort.Sort(t)
	}
}

type ElementType struct {
	List bool
	Type Type
}

// Flattens homogeneous lists of primitive types and squashes lists of struct
// into one struct. If fields have conflicting types while squashing a
// list of struct, the offending field is converted to the empty interface.
func (t *Tree) Normalize() {
	for idx := range t.Children {
		t.Children[idx].Normalize()
	}

	if !t.List || len(t.Children) == 0 {
		return
	}

	types := make(map[Type]bool)
	for _, child := range t.Children {
		types[child.Type] = true
	}
	switch len(types) {
	case 1:
		// Get first key out of 1-element map.
		for typ := range types {
			t.Type = typ
		}

		if t.List && t.Type == Struct {
			fields := make(map[Ident]*Tree)

			for _, element := range t.Children {
				for _, child := range element.Children {
					if _, exists := fields[child.Name]; !exists {
						fields[child.Name] = child
					} else {
						if !Compare(fields[child.Name], child) {
							fields[child.Name].Type = Interface
							fields[child.Name].Children = nil
						}
					}
				}
			}

			t.Children = nil
			for _, child := range fields {
				t.Children = append(t.Children, child)
			}
			sort.Sort(t)
		} else {
			t.Children = nil
		}
	case 2:
		if types[Int] && types[Float] {
			t.Type = Float
			t.Children = nil
		}
	default:
		t.Type = Interface
		t.Children = nil
	}
}

// Used for comparing fields between structs while squashing a list of struct.
type FieldType struct {
	Name Ident
	List bool
	Type Type
}

// Recursively compares field names and types of two structs.
func Compare(t1, t2 *Tree) bool {
	c1, c2 := Walker(t1), Walker(t2)
	for {
		v1, ok1 := <-c1
		v2, ok2 := <-c2
		if !ok1 || !ok2 {
			return ok1 == ok2
		}
		if v1 != v2 {
			break
		}
	}
	return false
}

func Walker(t *Tree) <-chan FieldType {
	ch := make(chan FieldType)
	go func() {
		Walk(t, ch)
		close(ch)
	}()
	return ch
}

func Walk(t *Tree, ch chan FieldType) {
	if t == nil {
		return
	}

	ch <- FieldType{t.Name, t.List, t.Type}
	for _, child := range t.Children {
		Walk(child, ch)
	}
}

func init() {
	log.SetFlags(log.Lshortfile)

	if err := config.Parse(); err != nil {
		log.Fatal("Error parsing flags:", err)
	}
}

func main() {
	defer config.Close()

	jsonDecoder := json.NewDecoder(config.inputFile)
	jsonDecoder.UseNumber()
	var data interface{}
	err := jsonDecoder.Decode(&data)
	if err != nil {
		log.Fatal("Error decoding input: ", err)
	}

	var tree Tree
	tree.Populate(data)
	if config.normalize {
		tree.Normalize()
	}

	indented, err := json.MarshalIndent(tree, "", "\t")
	if err != nil {
		log.Fatal("Error encoding tree:", err)
	}

	_, err = config.dumpFile.Write(indented)
	if err != nil {
		log.Fatal("Error dumping tree:", err)
	}

	source, err := tree.Format()
	fmt.Println(string(source))
	if err != nil {
		log.Fatal("Error formatting source:", err)
	}
}
