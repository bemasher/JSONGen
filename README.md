## Purpose
JSONGen is a tool for generating native Golang types from JSON objects. This automates what is otherwise a very tedious and error prone task when working with JSON.

## Usage
Using [test.json](example/test.json) as input the example will produce:
```go
type _ struct {
	Foo            string    `json:"foo"`
	Intlist        []float64 `json:"intlist"`
	Field_conflict []struct {
		Foo        interface{} `json:"foo"`
		Bar        float64     `json:"bar"`
		Baz        bool        `json:"baz"`
		Intlist    []float64   `json:"intlist"`
		Stringlist []string    `json:"stringlist"`
		Boollist   []bool      `json:"boollist"`
	} `json:"field-conflict"`
	Non_homogeneous []interface{} `json:"non-homogeneous"`
	Nil             nil           `json:"nil"`
	Baz             bool          `json:"baz"`
	Stringlist      []string      `json:"stringlist"`
	_Sanitary       string
	Sanitary0       string
	Unsanitary      string  `json:"0Unsanitary"`
	Bar             float64 `json:"bar"`
	Boollist        []bool  `json:"boollist"`
	Compound        struct {
		Foo        string    `json:"foo"`
		Bar        float64   `json:"bar"`
		Baz        bool      `json:"baz"`
		Intlist    []float64 `json:"intlist"`
		Stringlist []string  `json:"stringlist"`
		Boollist   []bool    `json:"boollist"`
	} `json:"compound"`
	Compoundlist []struct {
		Foo        string    `json:"foo"`
		Bar        float64   `json:"bar"`
		Baz        bool      `json:"baz"`
		Intlist    []float64 `json:"intlist"`
		Stringlist []string  `json:"stringlist"`
		Boollist   []bool    `json:"boollist"`
	} `json:"compoundlist"`
	Sanitary string
}
```

## Usage

```
$ jsongen -h
Usage of jsongen:
  -input="/dev/stdin": Filename to parse and generate type from, or omit for stdin.
```

Reading from stdin can be done as follows:
```
	jsongen < test.json
```

Reading directly from a file is then:
```
	jsongen -input=test.json
```

## Parsing
### Field Names
  * Field names are sanitized and written as exported fields of the generated type.
  * If sanitizing produces an empty string the original field name is prefixed with an underscore and only invalid identifier characters are removed.
    * The initial sanitizing method trims digits from the left of the identifier. This step performed on a field name of "12345" would produce an empty string. At this point the field name is instead stripped of only invalid characters like punctuation and prefixed with an underscore.
  * If sanitizing produces a field name different from the original value a JSON tag is added to the field allowing parsing after the field name has been modified.

## Types
### Primitive
  * Primitive types are parsed and stored as-is.
  * Valid types are bool, float64 and string.
  * The JSON value `null` is treated as the empty interface.

### Compound
  * Compound types are treated as structs.
  * The top-level object must be a compound type.
  * Fields of compound structures have no guaranteed order.
  * If a compound structure contains duplicate fields of different types, one of the fields is chosen at random. This is due to golang's unordered iteration over map entries. This should never occur since it is not permitted in the JSON specification, but this is the expected behavior should it happen.

### Lists
  * A homogeneous list of primitive  values are treated as a list of the primitive type e.g.: `[]float64`
  * Lists of heterogeneous types are treated as a list of the empty interface: `[]interface{}`
  * Lists with compound elements are treated as an array of structs.
    * Fields of each element are "squashed" into a single struct. The result is an array of a struct containing all encountered fields.
    * If a field in one element has a different type in another of the same list, the offending field is treated as an empty interface.

Examples of all of the above can be found in [test.json](test.json).

## Caveats
  * Currently field names within a struct are considered unique based on their unsanitized form. This could be troublesome if sanitizing produces non-unique field names of siblings. This also complicates the handling of field tags in the case of unique unsanitized names which sanitize to non-unique names.