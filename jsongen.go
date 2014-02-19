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

var filename string

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

	if t.Type == Nil {
		t.Type = Interface
	}

	r += fmt.Sprintf("%s %s", t.Type, tag)

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
			return '_'
		}
		return -1
	}
	s = strings.Map(valid, s)

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

// Types of nodes.
type Type byte

const (
	Nil Type = iota
	Bool
	Number
	String
	Interface
)

func (t Type) String() string {
	return []string{"nil", "bool", "float64", "string", "interface{}"}[t]
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
		t.Kind = Primitive
		t.Type = Number
	case string:
		t.Kind = Primitive
		t.Type = String
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

	flag.StringVar(&filename, "input", "/dev/stdin", "Filename to parse and generate type from, or omit for stdin.")
	flag.Parse()
}

func main() {
	var (
		inputFile *os.File
		err       error
	)
	if filename == "/dev/stdin" {
		inputFile = os.Stdin
	} else {
		inputFile, err = os.Open(filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer inputFile.Close()
	}

	jsonDecoder := json.NewDecoder(inputFile)
	var data interface{}
	err = jsonDecoder.Decode(&data)
	if err != nil {
		log.Fatal("error decoding input: ", err)
	}

	var tree Tree
	tree.Populate(data, "")
	tree.Normalize()

	fmt.Println(tree.Format())
}
