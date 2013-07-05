package jsongen

import (
	"fmt"
	"go/format"
	"log"
	"strings"
)

type Name string

// Returns sanitized field names
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

type Kind string

func (k Kind) String() string {
	switch k {
	case "<nil>":
		return "string"
	}
	return string(k)
}

// Type is used for describing the structure of decoded JSON data.
type Type struct {
	Name Name
	Kind Kind
	Tag  string

	IsList     bool
	IsCompound bool
	Fields     map[string]Type
}

func (t Type) String() string {
	return fmt.Sprintf("{Name:%s Type:%s IsList:%t IsCompound:%t}", t.Name, t.Kind, t.IsList, t.IsCompound)
}

func (t Type) Format() string {
	str := t.formatHelper(0)

	formatted, err := format.Source([]byte(str))
	if err != nil {
		log.Fatal("Error formatting type:", err)
	}

	return string(formatted)
}

func (t Type) formatHelper(depth int) (r string) {
	if t.Name.String() != string(t.Name) {
		t.Tag = fmt.Sprintf("`json:\"%s\"`", string(t.Name))
	}

	indent := strings.Repeat("\t", depth)
	r += indent
	if depth == 0 {
		r += "type "
		t.IsCompound = true
	}

	if t.IsList {
		r += fmt.Sprintf("%s []%s", t.Name, t.Kind)
		if !t.IsCompound {
			r += t.Tag
			return
		}
	}

	if t.IsCompound {
		if !t.IsList {
			r += fmt.Sprintf("%s struct", t.Name)
		}
		r += " {\n"
		defer func() {
			r += indent + "}" + t.Tag
		}()

		for _, f := range t.Fields {
			r += f.formatHelper(depth+1) + "\n"
		}
		return
	}

	r += fmt.Sprintf("%s %s %s", t.Name, t.Kind, t.Tag)

	return
}

// Parse takes a string defining the name of the type and an interface{} to populate a type.
func Parse(name string, data interface{}) (t Type) {
	t.Name = Name(name)

	switch T := data.(type) {
	case map[string]interface{}:
		t.IsCompound = true
		t.Kind = Kind("struct")
		t.Fields = make(map[string]Type)
		for k, v := range T {
			t.Fields[k] = Parse(k, v)
		}
	case []interface{}:
		t.IsList = true
		listTypes := make(map[Kind]bool)
		for _, i := range T {
			k := Kind(fmt.Sprintf("%T", i))
			listTypes[k] = true
			t.Kind = k

			// If the list has more than one type, halt parsing
			if len(listTypes) > 1 {
				t.Kind = Kind("interface{}")
				return
			}
		}

		if t.Kind == Kind("map[string]interface {}") {
			t.IsCompound = true
			t.Kind = Kind("struct")
			t.Fields = make(map[string]Type)
			for _, i := range T {
				tmp := Parse(name, i)
				for k, v := range tmp.Fields {
					if _, exists := t.Fields[k]; exists && t.Fields[k].Kind != v.Kind {
						t.IsCompound = false
						t.Kind = Kind("interface{}")
						t.Fields = nil
						return
					}
					t.Fields[k] = v
				}
			}
		}

	default:
		t.Name = Name(name)
		t.Kind = Kind(fmt.Sprintf("%T", T))
	}

	return
}

func init() {
	log.SetFlags(log.Lshortfile)
}
