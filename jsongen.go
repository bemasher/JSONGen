package jsongen

import (
	"fmt"
	"go/format"
	"log"
	"strings"
)

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

// Maps some special cases to valid golang types.
type Kind string

func (k Kind) String() string {
	switch k {
	case "<nil>":
		return "string"
	}
	return string(k)
}

// Used for describing the structure of decoded JSON data.
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

// Returns canonical form golang of the type structure.
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

	// There are three base cases of types: concrete, compound and list.
	switch T := data.(type) {

	// A compound json object.
	case map[string]interface{}:
		t.IsCompound = true
		t.Kind = Kind("struct")
		t.Fields = make(map[string]Type)

		// Recurse on the compound object's fields.
		for k, v := range T {
			t.Fields[k] = Parse(k, v)
		}

	// A list of json objects
	case []interface{}:
		t.IsList = true

		// Determine if the list is homogeneous or not.
		listTypes := make(map[Kind]bool)
		for _, i := range T {
			k := Kind(fmt.Sprintf("%T", i))
			listTypes[k] = true
			t.Kind = k

			// If the list has more than one type, it is heterogeneous
			// and is represented as a list of empty interfaces.
			if len(listTypes) > 1 {
				t.Kind = Kind("interface{}")
				return
			}
		}

		// If this is a list of compound tags, recurse on each item
		// and merge fields together from all elements of the list.
		if t.Kind == Kind("map[string]interface {}") {
			t.IsCompound = true
			t.Kind = Kind("struct")
			t.Fields = make(map[string]Type)
			for _, i := range T {
				tmp := Parse(name, i)
				for k, v := range tmp.Fields {
					// If we've previously seen this field name and it has a different
					// type than the last field, then stop parsing and treat list as a
					// list of empty interfaces.
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

	// This field must be a concrete type.
	default:
		t.Name = Name(name)
		t.Kind = Kind(fmt.Sprintf("%T", T))
	}

	return
}

func init() {
	log.SetFlags(log.Lshortfile)
}
