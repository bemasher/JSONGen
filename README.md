## Purpose
JSONGen is a tool for generating native Golang types from JSON objects. This automates what is otherwise a very tedious and error prone task when working with JSON.

[![Build Status](http://img.shields.io/travis/bemasher/JSONGen.svg?style=flat)](https://travis-ci.org/bemasher/JSONGen)
[![GPLv3 License](http://img.shields.io/badge/license-GPLv3-blue.svg?style=flat)](http://choosealicense.com/licenses/gpl-3.0/)

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
	Bool              bool          `json:"bool"`
	Boollist          []bool        `json:"boollist"`
	Float             float64       `json:"float"`
	Floatlist         []float64     `json:"floatlist"`
	Heterogeneouslist []interface{} `json:"heterogeneouslist"`
	Int               int64         `json:"int"`
	Intlist           []int64       `json:"intlist"`
	Nil               interface{}   `json:"nil"`
	Nillist           []interface{} `json:"nillist"`
	Sanitary          string        `json:"_Sanitary"`
	Sanitary          string
	Sanitary0         string
	String            string   `json:"string"`
	Stringlist        []string `json:"stringlist"`
	Struct            struct {
		Bool   bool        `json:"bool"`
		Float  float64     `json:"float"`
		Int    int64       `json:"int"`
		Nil    interface{} `json:"nil"`
		String string      `json:"string"`
	} `json:"struct"`
	Structlist []struct {
		Bool   bool    `json:"bool"`
		Float  float64 `json:"float"`
		Int    int64   `json:"int"`
		String string  `json:"string"`
	} `json:"structlist"`
	Structlistsquash []struct {
		Bool   bool    `json:"bool"`
		Float  float64 `json:"float"`
		Int    int64   `json:"int"`
		String string  `json:"string"`
	} `json:"structlistsquash"`
	Structlistsquashconflict []struct {
		Bool     bool        `json:"bool"`
		Conflict interface{} `json:"conflict"`
		Float    float64     `json:"float"`
		Int      int64       `json:"int"`
		String   string      `json:"string"`
	} `json:"structlistsquashconflict"`
	TitleCase  string `json:"title case"`
	TitleCase  string `json:"title_case"`
	TitleCase  string `json:"title-case"`
	Titlecase  string `json:"titlecase"`
	Unsanitary string `json:"0Unsanitary"`
	_          string `json:"123"`
}
```

## Parsing
### Field Names
  * Field names are sanitized and written as exported fields of the generated type.
  * If sanitizing produces an empty string the identifier is changed to `_`, this will need to be set by hand in order to properly decode the type.
  * If sanitizing produces a field name different from the original value a JSON tag is added to the field.
  * Spaces and `-` are converted to `_`.
  * Field names are converted to title case treating `_` and `-` as word boundaries along with spaces. This can be disabled using `-title=false`.

## Types
### Primitive
  * Primitive types are parsed and stored as-is.
  * Valid types are bool, int64, float64 and string.
  * The JSON value `null` is translated to the empty interface.

### Object
  * Object types are treated as structs.
  * Fields of structures are sorted lexicographically by sanitized field name.
  * If a structure contains duplicate fields of different types, the type will be chosen at random since Golang's map iteration order is undefined. This shouldn't be a problem since this is not permitted in JSON specification, but this is the expected behavior should it happen.

### Lists
  * A homogeneous list of primitive values are treated as a list of the primitive type e.g.: `[]float64`
  * Lists of heterogeneous types are treated as a list of the empty interface: `[]interface{}`
  * Lists containing both integers and floating point values are interpreted as `[]float64`.
  * Lists with object elements are treated as a list of structs.
    * Fields of each element are "squashed" into a single struct. The result is an array of a struct containing all encountered fields.   
    * If a field in one element has a different type in another of the same list, the offending field is treated as an empty interface.

Examples of all of the above can be found in [test.json](test.json).

## Caveats
  * Currently sibling field names are not guaranteed to be unique.

### License
The source of this project is licensed under GNU GPL v3.0, according to [http://choosealicense.com/licenses/gpl-3.0/](http://choosealicense.com/licenses/gpl-3.0/):

#### Required:

 * **Disclose Source**: Source code must be made available when distributing the software. In the case of LGPL, the source for the library (and not the entire program) must be made available.
 * **License and copyright notice**: Include a copy of the license and copyright notice with the code.
 * **State Changes**: Indicate significant changes made to the code.

#### Permitted:

 * **Commercial Use**: This software and derivatives may be used for commercial purposes.
 * **Distribution**: You may distribute this software.
 * **Modification**: This software may be modified.
 * **Patent Grant**: This license provides an express grant of patent rights from the contributor to the recipient.
 * **Private Use**: You may use and modify the software without distributing it.

#### Forbidden:

 * **Hold Liable**: Software is provided without warranty and the software author/license owner cannot be held liable for damages.
 * **Sublicensing**: You may not grant a sublicense to modify and distribute this software to third parties not included in the license.

## Feedback
If you find a case that produces incorrect results or you have a feature suggestion, let me know: submit an issue.
