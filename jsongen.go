package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"strings"
)

const (
	PrettyPrint = true
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

type Tree struct {
	Key      Name
	Kind     Kind
	Type     Type
	Children []Tree
}

// Returns canonical form golang of the type structure.
func (t Tree) Format() string {
	str := "type " + t.formatHelper(0)

	formatted, err := format.Source([]byte(str))
	if err != nil {
		fmt.Println(str)
		log.Fatal("Error formatting type:", err)
	}

	return string(formatted)
}

func (t Tree) formatHelper(depth int) (r string) {
	var tag string

	if depth != 0 && t.Key.String() != string(t.Key) {
		tag = fmt.Sprintf("`json:\"%s\"`", string(t.Key))
	}

	indent := strings.Repeat("\t", depth)
	r += indent

	r += t.Key.String() + " "

	if t.Kind == Array || t.Kind == ArrayOfStruct {
		r += fmt.Sprintf("[]")
	}

	if t.Kind == Struct || t.Kind == ArrayOfStruct {
		r += "struct {\n"
		defer func() {
			r += indent + "} " + tag
		}()

		for _, f := range t.Children {
			r += f.formatHelper(depth+1) + "\n"
		}
		return
	}

	r += fmt.Sprintf("%s %s", t.Type.Repr(), tag)

	return
}

// Sanitizes field names.
type Name string

func (n Name) String() (s string) {
	s = strings.TrimLeft(string(n), "0123456789")

	valid := func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		if r >= 'a' && r <= 'z' {
			return r
		}
		if r >= 'A' && r <= 'Z' {
			return r
		}
		if r == '_' || r == '-' {
			if config.titleCase {
				return ' '
			} else {
				return '_'
			}
		}
		return -1
	}
	s = strings.Map(valid, s)

	s = strings.Title(s)
	s = strings.Replace(s, " ", "", -1)

	if len(s) == 0 {
		return "_" + strings.Map(valid, string(n))
	}
	return strings.Title(s)
}

// Type kinds of nodes.
type Kind byte

const (
	Primitive Kind = iota
	Struct
	Array
	ArrayOfStruct
)

func (k Kind) String() string {
	return []string{"Primitive", "Struct", "Array", "ArrayOfStruct"}[k]
}

func (k Kind) MarshalText() (text []byte, err error) {
	text = []byte(k.String())
	return
}

// Types of nodes.
type Type byte

const (
	Nil Type = iota
	Bool
	IntNumber
	FloatNumber
	String
	Interface
)

func (t Type) String() string {
	return []string{"Null", "Bool", "Intnumber", "FloatNumber", "String", "Unknown"}[t]
}

func (t Type) Repr() string {
	return []string{"interface{}", "bool", "int", "float64", "string", "interface{}"}[t]
}

func (t Type) MarshalText() (text []byte, err error) {
	text = []byte(t.String())
	return
}

// Given an empty interface which json has been parsed into, populates the tree.
func (t *Tree) Populate(data interface{}, key string) {
	t.Key = Name(key)

	switch typ := data.(type) {
	case map[string]interface{}:
		t.Kind = Struct
		for k, v := range typ {
			var child Tree
			child.Populate(v, k)
			t.Children = append(t.Children, child)
		}
	case []interface{}:
		t.Kind = Array
		for _, v := range typ {
			var child Tree
			child.Populate(v, "")
			t.Children = append(t.Children, child)
		}
	case bool:
		t.Kind = Primitive
		t.Type = Bool
	case float64:
		// unused branch
		t.Kind = Primitive
		t.Type = FloatNumber
	case string:
		t.Kind = Primitive
		t.Type = String
	case json.Number:
		_, err := data.(json.Number).Int64()
		if err != nil {
                    t.Kind = Primitive
		    t.Type = FloatNumber
		} else {
		    t.Kind = Primitive
		    t.Type = IntNumber
		}
	default:
		t.Kind = Primitive
		t.Type = Nil
	}
}

// Flattens arrays of both primitive and compound types.
func (t *Tree) Normalize() {
	t.normalizePrimitiveArray()
	t.normalizeCompoundArray()
}

func (t *Tree) normalizePrimitiveArray() {
	// Traverse in depth-first order.
	for idx := range t.Children {
		t.Children[idx].normalizePrimitiveArray()
	}

	if t.Kind == Array {
		var typ Type
		isPrimitive := true
		isHomogeneous := true
		for _, c := range t.Children {
			if c.Kind != Primitive {
				isPrimitive = false
				isHomogeneous = false
				break
			}
			if typ == Nil {
				typ = c.Type
			}
			if typ != c.Type {
				isHomogeneous = false
				break
			}
		}

		if isPrimitive {
			if isHomogeneous {
				t.Type = typ
			}
			if !isHomogeneous || len(t.Children) == 0 {
				t.Type = Interface
			}
			t.Children = nil
		}
	}
}

func (t *Tree) normalizeCompoundArray() {
	// Traverse in depth-first order.
	for idx := range t.Children {
		t.Children[idx].normalizeCompoundArray()
	}

	if t.Kind == Array {
		fields := make(map[Name]Tree)
		for _, listChild := range t.Children {
			for _, structChild := range listChild.Children {
				if _, exists := fields[structChild.Key]; exists {
					if fields[structChild.Key].Type != structChild.Type {
						field := fields[structChild.Key]
						field.Type = Interface
						fields[structChild.Key] = field
					}
				} else {
					fields[structChild.Key] = structChild
				}
			}
		}

		if len(fields) > 0 && t.Type == Nil {
			t.Kind = ArrayOfStruct
			t.Children = nil
			for _, val := range fields {
				t.Children = append(t.Children, val)
			}
		}
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
	tree.Populate(data, "")
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

	fmt.Println(tree.Format())
}
