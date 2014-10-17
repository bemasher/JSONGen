package main

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

func Parse(raw string) (tree Tree, err error) {
	var data interface{}

	buf := bytes.NewBufferString(raw)
	jsonDecoder := json.NewDecoder(buf)
	jsonDecoder.UseNumber()

	err = jsonDecoder.Decode(&data)

	if err != nil {
		return
	}

	tree.Populate(data)
	tree.Normalize()

	return
}

type TreeTestCase struct {
	Source string
	Tree   Tree
}

func (tc TreeTestCase) TestTree(t *testing.T) {
	tree, err := Parse(tc.Source)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(tree, tc.Tree) {
		t.Errorf("Expected: %+v Got: %#v", tc.Tree, tree)
	}

	formatted, err := tree.Format()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s\n", tc.Source)
	t.Logf("%s\n", string(formatted))
}

func TestNil(t *testing.T) {
	testCases := []TreeTestCase{{`null`, Tree{Type: Interface}}}

	for _, testCase := range testCases {
		testCase.TestTree(t)
	}
}

func TestBool(t *testing.T) {
	testCases := []TreeTestCase{
		{`true`, Tree{Type: Bool}},
		{`false`, Tree{Type: Bool}},
	}

	for _, testCase := range testCases {
		testCase.TestTree(t)
	}
}

func TestInt(t *testing.T) {
	testCases := []TreeTestCase{
		{`-1`, Tree{Type: Int}},
		{`0`, Tree{Type: Int}},
		{`1`, Tree{Type: Int}},
		{`42`, Tree{Type: Int}},
	}

	for _, testCase := range testCases {
		testCase.TestTree(t)
	}
}

func TestFloat(t *testing.T) {
	testCases := []TreeTestCase{
		{`-1.0`, Tree{Type: Float}},
		{`0.0`, Tree{Type: Float}},
		{`1.0`, Tree{Type: Float}},
		{`42.0`, Tree{Type: Float}},
	}

	for _, testCase := range testCases {
		testCase.TestTree(t)
	}
}

func TestString(t *testing.T) {
	testCases := []TreeTestCase{
		{`""`, Tree{Type: String}},
		{`"foo"`, Tree{Type: String}},
	}

	for _, testCase := range testCases {
		testCase.TestTree(t)
	}
}

func TestBoolList(t *testing.T) {
	testCases := []TreeTestCase{
		{`[true]`, Tree{Type: Bool, List: true}},
		{`[false]`, Tree{Type: Bool, List: true}},
		{`[true, false]`, Tree{Type: Bool, List: true}},
	}

	for _, testCase := range testCases {
		testCase.TestTree(t)
	}
}

func TestIntList(t *testing.T) {
	testCases := []TreeTestCase{
		{`[-1]`, Tree{Type: Int, List: true}},
		{`[0]`, Tree{Type: Int, List: true}},
		{`[1]`, Tree{Type: Int, List: true}},
		{`[42]`, Tree{Type: Int, List: true}},
		{`[-1, 0]`, Tree{Type: Int, List: true}},
		{`[-1, 0, 1]`, Tree{Type: Int, List: true}},
		{`[-1, 0, 1, 42]`, Tree{Type: Int, List: true}},
	}

	for _, testCase := range testCases {
		testCase.TestTree(t)
	}
}

func TestFloatList(t *testing.T) {
	testCases := []TreeTestCase{
		{`[-1.0]`, Tree{Type: Float, List: true}},
		{`[0.0]`, Tree{Type: Float, List: true}},
		{`[1.0]`, Tree{Type: Float, List: true}},
		{`[42.0]`, Tree{Type: Float, List: true}},
		{`[-1.0, 0.0]`, Tree{Type: Float, List: true}},
		{`[-1.0, 0.0, 1.0]`, Tree{Type: Float, List: true}},
		{`[-1.0, 0.0, 1.0, 42.0]`, Tree{Type: Float, List: true}},
		{`[-1.0, 0.0, 1, 42]`, Tree{Type: Float, List: true}},
		{`[-1, 0, 1.0, 42.0]`, Tree{Type: Float, List: true}},
	}

	for _, testCase := range testCases {
		testCase.TestTree(t)
	}
}

func TestStringList(t *testing.T) {
	testCases := []TreeTestCase{
		{`[""]`, Tree{Type: String, List: true}},
		{`["foo"]`, Tree{Type: String, List: true}},
		{`["", "foo"]`, Tree{Type: String, List: true}},
		{`["foo", ""]`, Tree{Type: String, List: true}},
		{`["foo", "bar"]`, Tree{Type: String, List: true}},
	}

	for _, testCase := range testCases {
		testCase.TestTree(t)
	}
}

func TestHeterogeneousList(t *testing.T) {
	testCases := []TreeTestCase{
		{`[true, false, 0, 1, 0.0, 1.0, "", "foo"]`, Tree{Type: Interface, List: true}},
	}

	for _, testCase := range testCases {
		testCase.TestTree(t)
	}
}

func TestStruct(t *testing.T) {
	testCases := []TreeTestCase{
		{`{}`, Tree{Type: Struct}},
		{`{"nil":null}`, Tree{Type: Struct, Children: []*Tree{{Name: "nil", Type: Interface}}}},
		{`{"bool":true}`, Tree{Type: Struct, Children: []*Tree{{Name: "bool", Type: Bool}}}},
		{`{"int":1}`, Tree{Type: Struct, Children: []*Tree{{Name: "int", Type: Int}}}},
		{`{"float":1.0}`, Tree{Type: Struct, Children: []*Tree{{Name: "float", Type: Float}}}},
		{`{"string":"foo"}`, Tree{Type: Struct, Children: []*Tree{{Name: "string", Type: String}}}},
		{`{"bool":[true,false]}`, Tree{Type: Struct, Children: []*Tree{{Name: "bool", Type: Bool, List: true}}}},
		{`{"int":[0,1]}`, Tree{Type: Struct, Children: []*Tree{{Name: "int", Type: Int, List: true}}}},
		{`{"float":[0.0,1.0]}`, Tree{Type: Struct, Children: []*Tree{{Name: "float", Type: Float, List: true}}}},
		{`{"string":["","foo"]}`, Tree{Type: Struct, Children: []*Tree{{Name: "string", Type: String, List: true}}}},
		{`{"heterogeneous":[true, false, 0, 1, 0.0, 1.0, "", "foo"]}`,
			Tree{Type: Struct, Children: []*Tree{{Name: "heterogeneous", Type: Interface, List: true}}},
		},
	}

	for _, testCase := range testCases {
		testCase.TestTree(t)
	}
}

func TestStructNormalize(t *testing.T) {
	testCases := []TreeTestCase{
		{`[
			{
				"bool":true,
				"int":1
			},
			{
				"float":0.0,
				"string":"foo"
			}
		]`,
			Tree{Type: Struct, List: true, Children: []*Tree{
				{Name: "bool", Type: Bool},
				{Name: "float", Type: Float},
				{Name: "int", Type: Int},
				{Name: "string", Type: String},
			}},
		},
		{`[
			{
				"bool":true,
				"int":1
			},
			{
				"bool":true,
				"int":1,
				"float":0.0,
				"string":"foo"
			}
		]`,
			Tree{Type: Struct, List: true, Children: []*Tree{
				{Name: "bool", Type: Bool},
				{Name: "float", Type: Float},
				{Name: "int", Type: Int},
				{Name: "string", Type: String},
			}},
		},
	}

	for _, testCase := range testCases {
		testCase.TestTree(t)
	}
}

func TestStructConflict(t *testing.T) {
	testCases := []TreeTestCase{
		{`[
			{
				"bool":true,
				"int":1
			},
			{
				"float":0.0,
				"string":"foo"
			},
			{
				"float":true,
				"string":1
			}
		]`,
			Tree{Type: Struct, List: true, Children: []*Tree{
				{Name: "bool", Type: Bool},
				{Name: "float", Type: Interface},
				{Name: "int", Type: Int},
				{Name: "string", Type: Interface},
			}},
		},
	}

	for _, testCase := range testCases {
		testCase.TestTree(t)
	}
}

func (tc TreeTestCase) TestFormat(t *testing.T) {
	source, err := tc.Tree.Format()

	if err != nil {
		t.Fatal(err)
	}

	if string(source) != tc.Source {
		t.Errorf("Expected: %q Got: %q", tc.Source, source)
	}
}

func TestInterfaceFormat(t *testing.T) {
	TreeTestCase{"type _ interface{}\n", Tree{Type: Interface}}.TestFormat(t)
	TreeTestCase{"type _ []interface{}\n", Tree{Type: Interface, List: true}}.TestFormat(t)
}

func TestBoolFormat(t *testing.T) {
	TreeTestCase{"type _ bool\n", Tree{Type: Bool}}.TestFormat(t)
	TreeTestCase{"type _ []bool\n", Tree{Type: Bool, List: true}}.TestFormat(t)
}

func TestIntFormat(t *testing.T) {
	TreeTestCase{"type _ int64\n", Tree{Type: Int}}.TestFormat(t)
	TreeTestCase{"type _ []int64\n", Tree{Type: Int, List: true}}.TestFormat(t)
}

func TestFloatFormat(t *testing.T) {
	TreeTestCase{"type _ float64\n", Tree{Type: Float}}.TestFormat(t)
	TreeTestCase{"type _ []float64\n", Tree{Type: Float, List: true}}.TestFormat(t)
}

func TestStringFormat(t *testing.T) {
	TreeTestCase{"type _ string\n", Tree{Type: String}}.TestFormat(t)
	TreeTestCase{"type _ []string\n", Tree{Type: String, List: true}}.TestFormat(t)
}

func TestStructFormat(t *testing.T) {
	testCases := []TreeTestCase{
		{"type _ struct {\n}\n", Tree{Type: Struct}},
		{"type _ struct {\n\tInterface interface{} `json:\"interface\"`\n}\n", Tree{Type: Struct, Children: []*Tree{{Name: "interface", Type: Interface}}}},
		{"type _ struct {\n\tBool bool `json:\"bool\"`\n}\n", Tree{Type: Struct, Children: []*Tree{{Name: "bool", Type: Bool}}}},
		{"type _ struct {\n\tInt int64 `json:\"int\"`\n}\n", Tree{Type: Struct, Children: []*Tree{{Name: "int", Type: Int}}}},
		{"type _ struct {\n\tFloat float64 `json:\"float\"`\n}\n", Tree{Type: Struct, Children: []*Tree{{Name: "float", Type: Float}}}},
		{"type _ struct {\n\tString string `json:\"string\"`\n}\n", Tree{Type: Struct, Children: []*Tree{{Name: "string", Type: String}}}},
		{"type _ struct {\n\tInterface []interface{} `json:\"interface\"`\n}\n", Tree{Type: Struct, Children: []*Tree{{Name: "interface", Type: Interface, List: true}}}},
		{"type _ struct {\n\tBool []bool `json:\"bool\"`\n}\n", Tree{Type: Struct, Children: []*Tree{{Name: "bool", Type: Bool, List: true}}}},
		{"type _ struct {\n\tInt []int64 `json:\"int\"`\n}\n", Tree{Type: Struct, Children: []*Tree{{Name: "int", Type: Int, List: true}}}},
		{"type _ struct {\n\tFloat []float64 `json:\"float\"`\n}\n", Tree{Type: Struct, Children: []*Tree{{Name: "float", Type: Float, List: true}}}},
		{"type _ struct {\n\tString []string `json:\"string\"`\n}\n", Tree{Type: Struct, Children: []*Tree{{Name: "string", Type: String, List: true}}}},
	}

	for _, testCase := range testCases {
		testCase.TestFormat(t)
	}
}

type SanitizerTestCase struct {
	Source, Sanitized string
	TitleCase         bool
}

func TestSanitizier(t *testing.T) {
	testCases := []SanitizerTestCase{
		{"Sanitary", "Sanitary", true},
		{"sanitary", "Sanitary", true},
		{"_Sanitary", "Sanitary", true},
		{"_sanitary", "Sanitary", true},
		{"Sanitary", "Sanitary", false},
		{"sanitary", "Sanitary", false},
		{"_Sanitary", "Sanitary", false},
		{"_sanitary", "Sanitary", false},

		{"titlecase", "Titlecase", true},
		{"title-case", "TitleCase", true},
		{"title_case", "TitleCase", true},
		{"title case", "TitleCase", true},
		{"titlecase", "Titlecase", false},
		{"title-case", "Title_case", false},
		{"title_case", "Title_case", false},
		{"title case", "TitleCase", false},

		{"123", "_", true},
	}

	for _, testCase := range testCases {
		sanitized := Ident(testCase.Source)
		config.titleCase = testCase.TitleCase
		if testCase.Sanitized != sanitized.String() {
			t.Fatalf("Expected: %q Got: %q\n", testCase.Sanitized, sanitized.String())
		}
	}
}
