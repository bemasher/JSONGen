package jsongen

import (
	"encoding/json"
	"reflect"
	"testing"
)

var root Type

func TestingParser(input []byte) (t Type, err error) {
	var raw interface{}
	err = json.Unmarshal(input, &raw)

	t = Parse("", raw)
	return
}

func TestBool(t *testing.T) {
	gtype, err := TestingParser([]byte(`{"foo": true}`))
	if err != nil {
		t.FailNow()
	}

	expected := root
	expected.Fields = make(map[string]Type)
	expected.Fields["foo"] = Type{
		Name:   "foo",
		Kind:   "bool",
		Fields: map[string]Type(nil),
	}

	if !reflect.DeepEqual(gtype, expected) {
		t.Fail()
	}
}

func TestFloat64(t *testing.T) {
	gtype, err := TestingParser([]byte(`{"foo": 1.0}`))
	if err != nil {
		t.FailNow()
	}

	expected := Type{
		Name:       "",
		Kind:       "struct",
		IsCompound: true,
		Fields: map[string]Type{
			"foo": Type{
				Name:   "foo",
				Kind:   "float64",
				Fields: map[string]Type(nil),
			},
		},
	}

	if !reflect.DeepEqual(gtype, expected) {
		t.Fail()
	}
}

func TestString(t *testing.T) {
	gtype, err := TestingParser([]byte(`{"foo": "bar"}`))
	if err != nil {
		t.FailNow()
	}

	expected := Type{
		Name:       "",
		Kind:       "struct",
		IsCompound: true,
		Fields: map[string]Type{
			"foo": Type{
				Name: "foo",
				Kind: "string",
			},
		},
	}

	if !reflect.DeepEqual(gtype, expected) {
		t.Fail()
	}
}

func TestNull(t *testing.T) {
	gtype, err := TestingParser([]byte(`{"foo": null}`))
	if err != nil {
		t.FailNow()
	}

	expected := root
	expected.Fields = make(map[string]Type)
	expected.Fields["foo"] = Type{
		Name: "foo",
		Kind: "<nil>",
	}

	if !reflect.DeepEqual(gtype, expected) {
		t.Fail()
	}

	if expected.Fields["foo"].Kind.String() != "string" {
		t.Fail()
	}
}

func TestCompound(t *testing.T) {
	gtype, err := TestingParser([]byte(`{"compound": {
		"foo": "stuff",
		"bar": 1,
		"baz": true,
		"intlist": [0, 1, 2, 3, 4],
		"stringlist": ["0", "1", "2", "3", "4"],
		"boollist": [true, false, true, false, true]
	}}`))
	if err != nil {
		t.FailNow()
	}

	expected := root
	expected.Fields = make(map[string]Type)
	expected.Fields["compound"] = Type{
		Name:       "compound",
		Kind:       "struct",
		IsCompound: true,
		Fields: map[string]Type{
			"foo": Type{
				Name: "foo",
				Kind: "string",
			},
			"bar": Type{
				Name: "bar",
				Kind: "float64",
			},
			"baz": Type{
				Name: "baz",
				Kind: "bool",
			},
			"intlist": Type{
				Name:   "intlist",
				Kind:   "float64",
				IsList: true,
			},
			"stringlist": Type{
				Name:   "stringlist",
				Kind:   "string",
				IsList: true,
			},
			"boollist": Type{
				Name:   "boollist",
				Kind:   "bool",
				IsList: true,
			},
		},
	}

	if !reflect.DeepEqual(gtype, expected) {
		t.Fail()
	}
}

func TestBoolList(t *testing.T) {
	gtype, err := TestingParser([]byte(`{"boollist": [true, false, true, false, true]}`))
	if err != nil {
		t.FailNow()
	}

	expected := root
	expected.Fields = make(map[string]Type)

	expected.Fields["boollist"] = Type{
		Name:   "boollist",
		Kind:   "bool",
		IsList: true,
	}

	if !reflect.DeepEqual(gtype, expected) {
		t.Fail()
	}
}

func TestNumberList(t *testing.T) {
	gtype, err := TestingParser([]byte(`{"intlist": [0, 1, 2, 3, 4]}`))
	if err != nil {
		t.FailNow()
	}

	expected := root
	expected.Fields = make(map[string]Type)

	expected.Fields["intlist"] = Type{
		Name:   "intlist",
		Kind:   "float64",
		IsList: true,
	}

	if !reflect.DeepEqual(gtype, expected) {
		t.Fail()
	}
}

func TestStringList(t *testing.T) {
	gtype, err := TestingParser([]byte(`{"stringlist": ["0", "1", "2", "3", "4"]}`))
	if err != nil {
		t.FailNow()
	}

	expected := root
	expected.Fields = make(map[string]Type)

	expected.Fields["stringlist"] = Type{
		Name:   "stringlist",
		Kind:   "string",
		IsList: true,
	}

	if !reflect.DeepEqual(gtype, expected) {
		t.Fail()
	}
}

func TestCompoundList(t *testing.T) {
	gtype, err := TestingParser([]byte(`{
	"compoundlist": [
		{
			"foo": "stuff",
			"bar": 1,
			"baz": true
		},
		{
			"foo": "stuff",
			"bar": 1,
			"baz": true,
			"intlist": [0, 1, 2, 3, 4],
			"stringlist": ["0", "1", "2", "3", "4"],
			"boollist": [true, false, true, false, true]
		},
		{
			"foo": "stuff",
			"bar": 1,
			"baz": true,
			"intlist": [0, 1, 2, 3, 4],
			"stringlist": ["0", "1", "2", "3", "4"],
			"boollist": [true, false, true, false, true]
		},
		{
			"foo": "stuff",
			"bar": 1,
			"baz": true
		}
	]}`))
	if err != nil {
		t.FailNow()
	}

	expected := root
	expected.Fields = make(map[string]Type)

	expected.Fields["compoundlist"] = Type{
		Name:       "compoundlist",
		Kind:       "struct",
		IsList:     true,
		IsCompound: true,
		Fields: map[string]Type{
			"foo": Type{
				Name: "foo",
				Kind: "string",
			},
			"bar": Type{
				Name: "bar",
				Kind: "float64",
			},
			"baz": Type{
				Name: "baz",
				Kind: "bool",
			},
			"intlist": Type{
				Name:   "intlist",
				Kind:   "float64",
				IsList: true,
			},
			"stringlist": Type{
				Name:   "stringlist",
				Kind:   "string",
				IsList: true,
			},
			"boollist": Type{
				Name:   "boollist",
				Kind:   "bool",
				IsList: true,
			},
		},
	}

	if !reflect.DeepEqual(gtype, expected) {
		t.Fail()
	}
}

func TestListFieldConflict(t *testing.T) {
	gtype, err := TestingParser([]byte(`{
	"field-conflict": [
		{
			"foo": "stuff",
			"bar": 1,
			"baz": true
		},
		{
			"foo": 1,
			"bar": 1,
			"baz": true,
			"intlist": [0, 1, 2, 3, 4],
			"stringlist": ["0", "1", "2", "3", "4"],
			"boollist": [true, false, true, false, true]
		}
	]}`))
	if err != nil {
		t.FailNow()
	}

	expected := root
	expected.Fields = make(map[string]Type)
	expected.Fields["field-conflict"] = Type{
		Name:       "field-conflict",
		Kind:       "struct",
		IsList:     true,
		IsCompound: true,
		Fields: map[string]Type{
			"foo": Type{
				Name: "foo",
				Kind: "interface{}",
			},
			"bar": Type{
				Name: "bar",
				Kind: "float64",
			},
			"baz": Type{
				Name: "baz",
				Kind: "bool",
			},
			"intlist": Type{
				Name:   "intlist",
				Kind:   "float64",
				IsList: true,
			},
			"stringlist": Type{
				Name:   "stringlist",
				Kind:   "string",
				IsList: true,
			},
			"boollist": Type{
				Name:   "boollist",
				Kind:   "bool",
				IsList: true,
			},
		},
	}

	if !reflect.DeepEqual(gtype, expected) {
		t.Fail()
	}
}

func TestHeterogeneousList(t *testing.T) {
	gtype, err := TestingParser([]byte(`{"non-homogeneous": [0, 1, true, false, "stuff", "and", "things"]}`))
	if err != nil {
		t.FailNow()
	}

	expected := root
	expected.Fields = make(map[string]Type)
	expected.Fields["non-homogeneous"] = Type{
		Name:   "non-homogeneous",
		Kind:   "interface{}",
		IsList: true,
	}

	if !reflect.DeepEqual(gtype, expected) {
		t.Fail()
	}
}

func TestIdentifierSanitize(t *testing.T) {
	sanitary := []string{"Sanitary", "_Sanitary", "Sanitary0"}
	for _, id := range sanitary {
		if id != Name(id).String() {
			t.Fail()
		}
	}

	unsanitary := []string{"0Unsanitary", "123"}
	for _, id := range unsanitary {
		if id == Name(id).String() {
			t.Fail()
		}
	}
}

func init() {
	root.Kind = Kind("struct")
	root.IsCompound = true
}
