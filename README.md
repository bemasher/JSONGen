## Purpose
JSONGen is a tool for generating native Golang types from JSON objects. This automates what is otherwise a very tedious and error prone task when working with JSON.

## Usage

```
$ jsongen -h
Usage of jsongen:
  -dump="NUL": Dump tree structure to file.
  -normalize=true: Squash arrays of struct and determine primitive array type.
  -title=true: Convert identifiers to title case, treating '_' and '-' as word boundaries.
```

Reading from stdin can be done as follows:
```
$ cat test.json | jsongen
```

Or a filename can be passed:
```
$ jsongen test.json
```

Using [test.json](test.json) as input the example will produce:
```go
type _ struct {
	Foo          string    `json:"foo"`
	Intlist      []int     `json:"intlist"`
	Floatlist    []float64 `json:"floatlist"`
	Compoundlist []struct {
		Foo        string   `json:"foo"`
		Bar        int      `json:"bar"`
		Baz        bool     `json:"baz"`
		Intlist    []int    `json:"intlist"`
		Stringlist []string `json:"stringlist"`
		Boollist   []bool   `json:"boollist"`
	} `json:"compoundlist"`
	FieldConflict []struct {
		Foo        interface{} `json:"foo"`
		Bar        int         `json:"bar"`
		Baz        bool        `json:"baz"`
		Intlist    []int       `json:"intlist"`
		Stringlist []string    `json:"stringlist"`
		Boollist   []bool      `json:"boollist"`
	} `json:"field-conflict"`
	Stringlist []string    `json:"stringlist"`
	Sanitary   string      `json:"_Sanitary"`
	Nil        interface{} `json:"nil"`
	Bar        int         `json:"bar"`
	Boollist   []bool      `json:"boollist"`
	Compound   struct {
		Foo        string   `json:"foo"`
		Bar        int      `json:"bar"`
		Baz        bool     `json:"baz"`
		Intlist    []int    `json:"intlist"`
		Stringlist []string `json:"stringlist"`
		Boollist   []bool   `json:"boollist"`
	} `json:"compound"`
	Sanitary0      string
	Baz            bool `json:"baz"`
	Sanitary       string
	Unsanitary     string        `json:"0Unsanitary"`
	NonHomogeneous []interface{} `json:"non-homogeneous"`
} 
```

## Parsing
### Field Names
  * Field names are sanitized and written as exported fields of the generated type.
  * If sanitizing produces an empty string the original field name is prefixed with an underscore and only invalid identifier characters are removed.
    * The initial sanitizing method trims digits from the left of the identifier. This step performed on a field name of "12345" would produce an empty string. At this point the field name is instead stripped of only invalid characters like punctuation and prefixed with an underscore.
  * If sanitizing produces a field name different from the original value a JSON tag is added to the field allowing parsing after the field name has been modified.
  * Field names are converted to title case treating '_' and '-' as word boundaries along with spaces.

## Types
### Primitive
  * Primitive types are parsed and stored as-is.
  * Valid types are bool, int64, float64 and string.
  * The JSON value `null` is translated to the empty interface.

### Object
  * Object types are treated as structs.
  * The top-level object must be either an object or list.
  * Fields of object structures have no guaranteed order.
  * If a object structure contains duplicate fields of different types, one of the fields is chosen at random. This is due to golang's unordered iteration over map entries. This should never occur since it is not permitted in the JSON specification, but this is the expected behavior should it happen.

### Lists
  * A homogeneous list of primitive  values are treated as a list of the primitive type e.g.: `[]float64`
  * Lists of heterogeneous types are treated as a list of the empty interface: `[]interface{}`
  * Lists with object elements are treated as an array of structs.
    * Fields of each element are "squashed" into a single struct. The result is an array of a struct containing all encountered fields.
    * If a field in one element has a different type in another of the same list, the offending field is treated as an empty interface.

Examples of all of the above can be found in [test.json](test.json).

## Caveats
  * Currently field names within a struct are considered unique based on their unsanitized form. This could be troublesome if sanitizing produces non-unique field names of siblings. This also complicates the handling of field tags in the case of unique unsanitized names which sanitize to non-unique names.
  * Lists containing both integers and floating point values are interpretted as a list of the empty interface. This functionality will eventually be implemented so that lists of mixed numbers are stored as floats.