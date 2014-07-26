package main

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func ParseTree(raw []byte) (tree Tree, err error) {
	var data interface{}

	buf := bytes.NewBuffer(raw)
	jsonDecoder := json.NewDecoder(buf)
	jsonDecoder.UseNumber()

	err = jsonDecoder.Decode(&data)

	if err != nil {
		return
	}

	tree.Populate(data, "")
	tree.Normalize()

	return
}

func TestNil(t *testing.T) {
	testCases := [][]byte{
		[]byte(`{"Foo":null}`),
	}

	for _, testCase := range testCases {
		tree, err := ParseTree(testCase)

		if err != nil {
			t.Fatal(err)
		}

		expected := Tree{Key: "", Kind: Struct, Children: []Tree{
			Tree{Key: "Foo", Kind: Primitive, Type: Nil},
		}}
		if !reflect.DeepEqual(tree, expected) {
			t.Fatalf("Expected: %+v Got: %#v", expected, tree)
		}
	}
}

func TestBool(t *testing.T) {
	testCases := [][]byte{
		[]byte(`{"Foo":true}`),
		[]byte(`{"Foo":false}`),
	}

	for _, testCase := range testCases {
		tree, err := ParseTree(testCase)

		if err != nil {
			t.Fatal(err)
		}

		expected := Tree{Key: "", Kind: Struct, Children: []Tree{
			Tree{Key: "Foo", Kind: Primitive, Type: Bool},
		}}
		if !reflect.DeepEqual(tree, expected) {
			t.Fatalf("Expected: %+v Got: %#v", expected, tree)
		}
	}
}

func TestInt(t *testing.T) {
	testCases := [][]byte{
		[]byte(`{"Foo":0}`),
		[]byte(`{"Foo":1}`),
		[]byte(`{"Foo":-1}`),
	}

	for _, testCase := range testCases {
		tree, err := ParseTree(testCase)

		if err != nil {
			t.Fatal(err)
		}

		expected := Tree{Key: "", Kind: Struct, Children: []Tree{
			Tree{Key: "Foo", Kind: Primitive, Type: Int},
		}}
		if !reflect.DeepEqual(tree, expected) {
			t.Fatalf("Expected: %+v Got: %#v", expected, tree)
		}
	}
}

func TestFloat(t *testing.T) {
	testCases := [][]byte{
		[]byte(`{"Foo":0.0}`),
		[]byte(`{"Foo":1.0}`),
		[]byte(`{"Foo":-1.0}`),
	}

	for _, testCase := range testCases {
		tree, err := ParseTree(testCase)

		if err != nil {
			t.Fatal(err)
		}

		expected := Tree{Key: "", Kind: Struct, Children: []Tree{
			Tree{Key: "Foo", Kind: Primitive, Type: Float},
		}}
		if !reflect.DeepEqual(tree, expected) {
			t.Fatalf("Expected: %+v Got: %#v", expected, tree)
		}
	}
}

func TestString(t *testing.T) {
	testCases := [][]byte{
		[]byte(`{"Foo":""}`),
		[]byte(`{"Foo":" "}`),
		[]byte(`{"Foo":"a"}`),
		[]byte(`{"Foo":"ab"}`),
		[]byte(`{"Foo":"abc"}`),
	}

	for _, testCase := range testCases {
		tree, err := ParseTree(testCase)

		if err != nil {
			t.Fatal(err)
		}

		expected := Tree{Key: "", Kind: Struct, Children: []Tree{
			Tree{Key: "Foo", Kind: Primitive, Type: String},
		}}
		if !reflect.DeepEqual(tree, expected) {
			t.Fatalf("Expected: %+v Got: %#v", expected, tree)
		}
	}
}

func TestBoolList(t *testing.T) {
	testCases := [][]byte{
		[]byte(`{"Foo":[true]}`),
		[]byte(`{"Foo":[false]}`),
		[]byte(`{"Foo":[true,true,true,true]}`),
		[]byte(`{"Foo":[false,false,false,false]}`),
		[]byte(`{"Foo":[true,false,true,false]}`),
	}

	for _, testCase := range testCases {
		tree, err := ParseTree(testCase)

		if err != nil {
			t.Fatal(err)
		}

		expected := Tree{Key: "", Kind: Struct, Children: []Tree{
			Tree{Key: "Foo", Kind: Array, Type: Bool},
		}}

		if !reflect.DeepEqual(tree, expected) {
			t.Fatalf("Expected: %+v Got: %#v", expected, tree)
		}
	}
}

func TestIntList(t *testing.T) {
	testCases := [][]byte{
		[]byte(`{"Foo":[0]}`),
		[]byte(`{"Foo":[1]}`),
		[]byte(`{"Foo":[-1]}`),

		[]byte(`{"Foo":[0,1]}`),
		[]byte(`{"Foo":[1,-1]}`),
		[]byte(`{"Foo":[-1,0]}`),
	}

	for _, testCase := range testCases {
		tree, err := ParseTree(testCase)

		if err != nil {
			t.Fatal(err)
		}

		expected := Tree{Key: "", Kind: Struct, Children: []Tree{
			Tree{Key: "Foo", Kind: Array, Type: Int},
		}}

		if !reflect.DeepEqual(tree, expected) {
			t.Fatalf("Expected: %+v Got: %#v", expected, tree)
		}
	}
}

func TestFloatList(t *testing.T) {
	testCases := [][]byte{
		[]byte(`{"Foo":[0.0]}`),
		[]byte(`{"Foo":[1.0]}`),
		[]byte(`{"Foo":[-1.0]}`),

		[]byte(`{"Foo":[0.0,1.0]}`),
		[]byte(`{"Foo":[1.0,-1.0]}`),
		[]byte(`{"Foo":[-1.0,0.0]}`),
	}

	for _, testCase := range testCases {
		tree, err := ParseTree(testCase)

		if err != nil {
			t.Fatal(err)
		}

		expected := Tree{Key: "", Kind: Struct, Children: []Tree{
			Tree{Key: "Foo", Kind: Array, Type: Float},
		}}

		if !reflect.DeepEqual(tree, expected) {
			t.Fatalf("Expected: %+v Got: %#v", expected, tree)
		}
	}
}

func TestStringList(t *testing.T) {
	testCases := [][]byte{
		[]byte(`{"Foo":[""]}`),
		[]byte(`{"Foo":[" "]}`),
		[]byte(`{"Foo":["a"]}`),
		[]byte(`{"Foo":["ab"]}`),
		[]byte(`{"Foo":["abc"]}`),

		[]byte(`{"Foo":["", " "]}`),
		[]byte(`{"Foo":[" ", "a"]}`),
		[]byte(`{"Foo":["a", "ab"]}`),
		[]byte(`{"Foo":["ab", "abc"]}`),
		[]byte(`{"Foo":["abc", ""]}`),
	}

	for _, testCase := range testCases {
		tree, err := ParseTree(testCase)

		if err != nil {
			t.Fatal(err)
		}

		expected := Tree{Key: "", Kind: Struct, Children: []Tree{
			Tree{Key: "Foo", Kind: Array, Type: String},
		}}

		if !reflect.DeepEqual(tree, expected) {
			t.Fatalf("Expected: %+v Got: %#v", expected, tree)
		}
	}
}

func TestCompound(t *testing.T) {
	testCases := [][]byte{
		[]byte(`{
			"Foo": true,
			"BarI": 1,
			"BarF": 1.0,
			"Baz": "a",
			"FooList": [true, false, true, false, true],
			"BarIList": [0, 1, 2, 3, 4],
			"BarFList": [0.0, 1.0, 2.0, 3.0, 4.0],
			"BazList": ["0", "1", "2", "3", "4"]
		}`),
	}

	for _, testCase := range testCases {
		tree, err := ParseTree(testCase)

		if err != nil {
			t.Fatal(err)
		}

		expected := Tree{Key: "", Kind: Struct, Children: []Tree{
			Tree{Key: "BarF", Kind: Primitive, Type: Float},
			Tree{Key: "BarFList", Kind: Array, Type: Float},
			Tree{Key: "BarI", Kind: Primitive, Type: Int},
			Tree{Key: "BarIList", Kind: Array, Type: Int},
			Tree{Key: "Baz", Kind: Primitive, Type: String},
			Tree{Key: "BazList", Kind: Array, Type: String},
			Tree{Key: "Foo", Kind: Primitive, Type: Bool},
			Tree{Key: "FooList", Kind: Array, Type: Bool},
		}}

		if !reflect.DeepEqual(tree, expected) {
			t.Fatalf("Expected: %+v Got: %#v", expected, tree)
		}
	}
}

func TestCompoundList(t *testing.T) {
	testCases := [][]byte{
		[]byte(`[
			{
				"Foo": true,
				"BarI": 1,
				"BarF": 1.0,
				"Baz": "a"
			},
			{
				"Foo": true,
				"BarI": 1,
				"BarF": 1.0,
				"Baz": "a",
				"FooList": [true, false, true, false, true],
				"BarIList": [0, 1, 2, 3, 4],
				"BarFList": [0.0, 1.0, 2.0, 3.0, 4.0],
				"BazList": ["0", "1", "2", "3", "4"]
			},
			{
				"Baz": "a",
				"BarI": 1,
				"BarF": 1.0,
				"Foo": true,
				"FooList": [true, false, true, false, true],
				"BarIList": [0, 1, 2, 3, 4],
				"BarFList": [0.0, 1.0, 2.0, 3.0, 4.0],
				"BazList": ["0", "1", "2", "3", "4"]
			},
			{
				"Baz": "a",
				"BarI": 1,
				"BarF": 1.0,
				"Foo": true
			}
		]`),
	}

	for _, testCase := range testCases {
		tree, err := ParseTree(testCase)

		if err != nil {
			t.Fatal(err)
		}

		expected := Tree{Key: "", Kind: ArrayOfStruct, Children: []Tree{
			Tree{Key: "BarF", Kind: Primitive, Type: Float},
			Tree{Key: "BarFList", Kind: Array, Type: Float},
			Tree{Key: "BarI", Kind: Primitive, Type: Int},
			Tree{Key: "BarIList", Kind: Array, Type: Int},
			Tree{Key: "Baz", Kind: Primitive, Type: String},
			Tree{Key: "BazList", Kind: Array, Type: String},
			Tree{Key: "Foo", Kind: Primitive, Type: Bool},
			Tree{Key: "FooList", Kind: Array, Type: Bool},
		}}

		if !reflect.DeepEqual(tree, expected) {
			t.Fatalf("Expected: %+v Got: %#v", expected, tree)
		}
	}
}

func TestNonHomogeneousList(t *testing.T) {
	testCases := [][]byte{
		[]byte(`[null, true, false, 0, 1, -1, 0.0, 1.0, -1.0, "", "a", "b", "c"]`),
	}

	for _, testCase := range testCases {
		tree, err := ParseTree(testCase)

		if err != nil {
			t.Fatal(err)
		}

		expected := Tree{Kind: Array, Type: Interface}

		if !reflect.DeepEqual(tree, expected) {
			t.Fatalf("Expected: %+v Got: %#v", expected, tree)
		}
	}
}

func TestCompoundFieldConflict(t *testing.T) {
	testCases := [][]byte{
		[]byte(`[
			{
				"Foo": true,
				"Bar": 1,
				"Baz": "a"
			},
			{
				"Foo": "a",
				"Bar": true,
				"Baz": 1
			}
		]`),
	}

	for _, testCase := range testCases {
		tree, err := ParseTree(testCase)

		if err != nil {
			t.Fatal(err)
		}

		expected := Tree{Kind: ArrayOfStruct, Children: []Tree{
			Tree{Key: "Bar", Kind: Primitive, Type: Interface},
			Tree{Key: "Baz", Kind: Primitive, Type: Interface},
			Tree{Key: "Foo", Kind: Primitive, Type: Interface},
		}}

		if !reflect.DeepEqual(tree, expected) {
			t.Fatalf("Expected: %+v Got: %#v", expected, tree)
		}
	}
}

func TestBoolFormat(t *testing.T) {
	tree := Tree{Key: "", Kind: Struct, Children: []Tree{
		Tree{Key: "Foo", Kind: Primitive, Type: Bool},
	}}

	expected := "type _ struct {\n\tFoo bool\n} "
	got := tree.Format()
	if got != expected {
		t.Fatalf("Expected: %q Got: %q", expected, got)
	}
}

func TestIntFormat(t *testing.T) {
	tree := Tree{Key: "", Kind: Struct, Children: []Tree{
		Tree{Key: "Foo", Kind: Primitive, Type: Int},
	}}

	expected := "type _ struct {\n\tFoo int64\n} "
	got := tree.Format()
	if got != expected {
		t.Fatalf("Expected: %q Got: %q", expected, got)
	}
}

func TestFloatFormat(t *testing.T) {
	tree := Tree{Key: "", Kind: Struct, Children: []Tree{
		Tree{Key: "Foo", Kind: Primitive, Type: Float},
	}}

	expected := "type _ struct {\n\tFoo float64\n} "
	got := tree.Format()
	if got != expected {
		t.Fatalf("Expected: %q Got: %q", expected, got)
	}
}

func TestStringFormat(t *testing.T) {
	tree := Tree{Key: "", Kind: Struct, Children: []Tree{
		Tree{Key: "Foo", Kind: Primitive, Type: String},
	}}

	expected := "type _ struct {\n\tFoo string\n} "
	got := tree.Format()
	if got != expected {
		t.Fatalf("Expected: %q Got: %q", expected, got)
	}
}

func TestBoolListFormat(t *testing.T) {
	tree := Tree{Key: "", Kind: Struct, Children: []Tree{
		Tree{Key: "Foo", Kind: Array, Type: Bool},
	}}

	expected := "type _ struct {\n\tFoo []bool\n} "
	got := tree.Format()
	if got != expected {
		t.Fatalf("Expected: %q Got: %q", expected, got)
	}
}

func TestIntListFormat(t *testing.T) {
	tree := Tree{Key: "", Kind: Struct, Children: []Tree{
		Tree{Key: "Foo", Kind: Array, Type: Int},
	}}

	expected := "type _ struct {\n\tFoo []int64\n} "
	got := tree.Format()
	if got != expected {
		t.Fatalf("Expected: %q Got: %q", expected, got)
	}
}

func TestFloatListFormat(t *testing.T) {
	tree := Tree{Key: "", Kind: Struct, Children: []Tree{
		Tree{Key: "Foo", Kind: Array, Type: Float},
	}}

	expected := "type _ struct {\n\tFoo []float64\n} "
	got := tree.Format()
	if got != expected {
		t.Fatalf("Expected: %q Got: %q", expected, got)
	}
}

func TestStringListFormat(t *testing.T) {
	tree := Tree{Key: "", Kind: Struct, Children: []Tree{
		Tree{Key: "Foo", Kind: Array, Type: String},
	}}

	expected := "type _ struct {\n\tFoo []string\n} "
	got := tree.Format()
	if got != expected {
		t.Fatalf("Expected: %q Got: %q", expected, got)
	}
}

func TestCompoundFormat(t *testing.T) {
	tree := Tree{Key: "", Kind: Struct, Children: []Tree{
		Tree{Key: "Foo", Kind: Primitive, Type: Bool},
		Tree{Key: "BarI", Kind: Primitive, Type: Int},
		Tree{Key: "BarF", Kind: Primitive, Type: Float},
		Tree{Key: "Baz", Kind: Primitive, Type: String},
		Tree{Key: "FooList", Kind: Array, Type: Bool},
		Tree{Key: "BarIList", Kind: Array, Type: Int},
		Tree{Key: "BarFList", Kind: Array, Type: Float},
		Tree{Key: "BazList", Kind: Array, Type: String},
	}}

	expected := "type _ struct {\n\tFoo      bool\n\tBarI     int64\n\tBarF     float64\n\tBaz      string\n\tFooList  []bool\n\tBarIList []int64\n\tBarFList []float64\n\tBazList  []string\n} "
	got := tree.Format()
	if got != expected {
		t.Fatalf("Expected: %q Got: %q", expected, got)
	}
}

func TestCompoundListFormat(t *testing.T) {
	tree := Tree{Key: "", Kind: ArrayOfStruct, Children: []Tree{
		Tree{Key: "Foo", Kind: Primitive, Type: Bool},
		Tree{Key: "BarI", Kind: Primitive, Type: Int},
		Tree{Key: "BarF", Kind: Primitive, Type: Float},
		Tree{Key: "Baz", Kind: Primitive, Type: String},
		Tree{Key: "FooList", Kind: Array, Type: Bool},
		Tree{Key: "BarIList", Kind: Array, Type: Int},
		Tree{Key: "BarFList", Kind: Array, Type: Float},
		Tree{Key: "BazList", Kind: Array, Type: String},
	}}

	expected := "type _ []struct {\n\tFoo      bool\n\tBarI     int64\n\tBarF     float64\n\tBaz      string\n\tFooList  []bool\n\tBarIList []int64\n\tBarFList []float64\n\tBazList  []string\n} "
	got := tree.Format()
	if got != expected {
		t.Fatalf("Expected: %q Got: %q", expected, got)
	}
}

func TestIdentifierSanitizer(t *testing.T) {
	sanitary := []string{"Sanitary", "_Sanitary", "Sanitary0"}
	titleReplacer := strings.NewReplacer("_", "", "-", "")

	for _, id := range sanitary {
		expected := id

		config.titleCase = false
		sanitized := Name(id).String()

		if expected != sanitized {
			t.Fatalf("Expected: %q Got: %q\n", expected, sanitized)
		}

		config.titleCase = true
		expected = titleReplacer.Replace(expected)
		sanitized = Name(id).String()

		if expected != sanitized {
			t.Fatalf("Expected: %q Got: %q\n", expected, sanitized)
		}
	}

	unsanitary := []string{"0Unsanitary", "123"}
	for _, id := range unsanitary {
		if id == Name(id).String() {
			t.Fail()
		}
	}
}

func TestCompoundTagFormat(t *testing.T) {
	tree := Tree{Key: "", Kind: Struct, Children: []Tree{
		Tree{Key: "foo", Kind: Primitive, Type: Bool},
		Tree{Key: "bari", Kind: Primitive, Type: Int},
		Tree{Key: "barf", Kind: Primitive, Type: Float},
		Tree{Key: "baz", Kind: Primitive, Type: String},
		Tree{Key: "foolist", Kind: Array, Type: Bool},
		Tree{Key: "barilist", Kind: Array, Type: Int},
		Tree{Key: "barflist", Kind: Array, Type: Float},
		Tree{Key: "bazlist", Kind: Array, Type: String},
	}}

	expected := "type _ struct {\n\tFoo      bool      `json:\"foo\"`\n\tBari     int64     `json:\"bari\"`\n\tBarf     float64   `json:\"barf\"`\n\tBaz      string    `json:\"baz\"`\n\tFoolist  []bool    `json:\"foolist\"`\n\tBarilist []int64   `json:\"barilist\"`\n\tBarflist []float64 `json:\"barflist\"`\n\tBazlist  []string  `json:\"bazlist\"`\n} "
	got := tree.Format()
	if got != expected {
		t.Fatalf("Expected: %q Got: %q", expected, got)
	}
}
