# Abstract JSON 

[![Build Status](https://travis-ci.com/spyzhov/ajson.svg?branch=master)](https://travis-ci.com/spyzhov/ajson)
[![Go Report Card](https://goreportcard.com/badge/github.com/spyzhov/ajson)](https://goreportcard.com/report/github.com/spyzhov/ajson)
[![GoDoc](https://godoc.org/github.com/spyzhov/ajson?status.svg)](https://godoc.org/github.com/spyzhov/ajson)
[![Coverage Status](https://coveralls.io/repos/github/spyzhov/ajson/badge.svg?branch=master)](https://coveralls.io/github/spyzhov/ajson?branch=master)

Abstract [JSON](https://www.json.org/) is a small golang package that provide a parser for JSON with support of JSONPath, in case when you are not sure in it's structure.

Method `Unmarshal` will scan all the byte slice to create a root node of JSON structure, with all it behaviors.

Method `Marshal` will serialize current `Node` object to JSON structure.

Each `Node` has it's own type and calculated value, which will be calculated on demand. 
Calculated value saves in `atomic.Value`, so it's thread safe.

Method `JSONPath` will returns slice of founded elements in current JSON data, by it's JSONPath.

## Example

Calculating `AVG(price)` when object is heterogeneous.

### Unmarshal

```go
package main

import (
	"fmt"
	"github.com/spyzhov/ajson"
)

func main() {
	data := []byte(`{ 
      "store": {
        "book": [ 
          { "category": "reference",
            "author": "Nigel Rees",
            "title": "Sayings of the Century",
            "price": 8.95
          },
          { "category": "fiction",
            "author": "Evelyn Waugh",
            "title": "Sword of Honour",
            "price": 12.99
          },
          { "category": "fiction",
            "author": "Herman Melville",
            "title": "Moby Dick",
            "isbn": "0-553-21311-3",
            "price": 8.99
          },
          { "category": "fiction",
            "author": "J. R. R. Tolkien",
            "title": "The Lord of the Rings",
            "isbn": "0-395-19395-8",
            "price": 22.99
          }
        ],
        "bicycle": { "color": "red",
          "price": 19.95
        },
        "tools": null
      }
    }`)

	root, err := ajson.Unmarshal(data)
	if err != nil {
		panic(err)
	}

	store := root.MustKey("store").MustObject()

	var prices float64
	size := 0
	for _, objects := range store {
		if objects.IsArray() && objects.Size() > 0 {
			size += objects.Size()
			for _, object := range objects.MustArray() {
				prices += object.MustKey("price").MustNumeric()
			}
		} else if objects.IsObject() && objects.HasKey("price") {
			size++
			prices += objects.MustKey("price").MustNumeric()
		}
	}

	if size > 0 {
		fmt.Println("AVG price:", prices/float64(size))
	} else {
		fmt.Println("AVG price:", 0)
	}
}
```

### JSONPath:

```go
package main

import (
	"fmt"
	"github.com/spyzhov/ajson"
)

func main() {
	data := []byte(`{ 
      "store": {
        "book": [ 
          { "category": "reference",
            "author": "Nigel Rees",
            "title": "Sayings of the Century",
            "price": 8.95
          },
          { "category": "fiction",
            "author": "Evelyn Waugh",
            "title": "Sword of Honour",
            "price": 12.99
          },
          { "category": "fiction",
            "author": "Herman Melville",
            "title": "Moby Dick",
            "isbn": "0-553-21311-3",
            "price": 8.99
          },
          { "category": "fiction",
            "author": "J. R. R. Tolkien",
            "title": "The Lord of the Rings",
            "isbn": "0-395-19395-8",
            "price": 22.99
          }
        ],
        "bicycle": { "color": "red",
          "price": 19.95
        },
        "tools": null
      }
    }`)

	nodes, err := ajson.JSONPath(data, "$..price")
	if err != nil {
		panic(err)
	}

	var prices float64
	size := len(nodes)
	for _, node := range nodes {
		prices += node.MustNumeric()
	}

	if size > 0 {
		fmt.Println("AVG price:", prices/float64(size))
	} else {
		fmt.Println("AVG price:", 0)
	}
}
```

### Eval

```go
package main

import (
	"fmt"
	"github.com/spyzhov/ajson"
)

func main() {
	json := []byte(`{ "store": {
    "book": [ 
      { "category": "reference",
        "author": "Nigel Rees",
        "title": "Sayings of the Century",
        "price": 8.95
      },
      { "category": "fiction",
        "author": "Evelyn Waugh",
        "title": "Sword of Honour",
        "price": 12.99
      },
      { "category": "fiction",
        "author": "Herman Melville",
        "title": "Moby Dick",
        "isbn": "0-553-21311-3",
        "price": 8.99
      },
      { "category": "fiction",
        "author": "J. R. R. Tolkien",
        "title": "The Lord of the Rings",
        "isbn": "0-395-19395-8",
        "price": 22.99
      }
    ],
    "bicycle": [
      {
        "color": "red",
        "price": 19.95
      }
    ]
  }
}`)
	root, err := ajson.Unmarshal(json)
	if err != nil {
		panic(err)
	}
	result, err := ajson.Eval(root, "avg($..price)")
	if err != nil {
		panic(err)
	}
	fmt.Println("AVG price:", result.MustNumeric())
}
```

### Marshal

```go
package main

import (
	"fmt"
	"github.com/spyzhov/ajson"
)

func main() {
	json := []byte(`{ "store": {
    "book": [ 
      { "category": "reference",
        "author": "Nigel Rees",
        "title": "Sayings of the Century",
        "price": 8.95
      },
      { "category": "fiction",
        "author": "Evelyn Waugh",
        "title": "Sword of Honour",
        "price": 12.99
      },
      { "category": "fiction",
        "author": "Herman Melville",
        "title": "Moby Dick",
        "isbn": "0-553-21311-3",
        "price": 8.99
      },
      { "category": "fiction",
        "author": "J. R. R. Tolkien",
        "title": "The Lord of the Rings",
        "isbn": "0-395-19395-8",
        "price": 22.99
      }
    ],
    "bicycle": [
      {
        "color": "red",
        "price": 19.95
      }
    ]
  }
}`)
	root := ajson.Must(ajson.Unmarshal(json))
	result := ajson.Must(ajson.Eval(root, "avg($..price)"))
	err := root.AppendObject("price(avg)", result)
	if err != nil {
		panic(err)
	}
	marshalled, err := ajson.Marshal(root)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", marshalled)
	// Output:
	// {"store":{
	//    "book": [
	//      { "category": "reference",
	//        "author": "Nigel Rees",
	//        "title": "Sayings of the Century",
	//        "price": 8.95
	//      },
	//      { "category": "fiction",
	//        "author": "Evelyn Waugh",
	//        "title": "Sword of Honour",
	//        "price": 12.99
	//      },
	//      { "category": "fiction",
	//        "author": "Herman Melville",
	//        "title": "Moby Dick",
	//        "isbn": "0-553-21311-3",
	//        "price": 8.99
	//      },
	//      { "category": "fiction",
	//        "author": "J. R. R. Tolkien",
	//        "title": "The Lord of the Rings",
	//        "isbn": "0-395-19395-8",
	//        "price": 22.99
	//      }
	//    ],
	//    "bicycle": [
	//      {
	//        "color": "red",
	//        "price": 19.95
	//      }
	//    ]
	//  },"price(avg)":14.774000000000001}
}
```

# JSONPath

Current package supports JSONPath selection described at [http://goessner.net/articles/JsonPath/](http://goessner.net/articles/JsonPath/).

JSONPath expressions always refer to a JSON structure in the same way as XPath expression are used in combination with an XML document. Since a JSON structure is usually anonymous and doesn't necessarily have a "root member object" JSONPath assumes the abstract name $ assigned to the outer level object.

JSONPath expressions can use the dot–notation

`$.store.book[0].title`

or the bracket–notation

`$['store']['book'][0]['title']`

for input pathes. Internal or output pathes will always be converted to the more general bracket–notation.

JSONPath allows the wildcard symbol `*` for member names and array indices. It borrows the descendant operator `..` from E4X and the array slice syntax proposal `[start:end:step]` from ECMASCRIPT 4.

Expressions of the underlying scripting language `(<expr>)` can be used as an alternative to explicit names or indices as in

`$.store.book[(@.length-1)].title`

using the symbol `@` for the current object. Filter expressions are supported via the syntax `?(<boolean expr>)` as in

`$.store.book[?(@.price < 10)].title`

Here is a complete overview and a side by side comparison of the JSONPath syntax elements with its XPath counterparts.

| JSONPath | Description |
|----------|---|
| `$`      | the root object/element |
| `@`      | the current object/element |
| `.` or `[]` | child operator |
| `..`     | recursive descent. JSONPath borrows this syntax from E4X. |
| `*`      | wildcard. All objects/elements regardless their names. |
| `[]`     | subscript operator. XPath uses it to iterate over element collections and for predicates. In Javascript and JSON it is the native array operator. |
| `[,]`    | Union operator in XPath results in a combination of node sets. JSONPath allows alternate names or array indices as a set. |
| `[start:end:step]` | array slice operator borrowed from ES4. |
| `?()`    | applies a filter (script) expression. |
| `()`     | script expression, using the underlying script engine. |

## Script engine

### Predefined constant

Package has several predefined constants. You are free to add new one with `AddConstant`

     e       math.E     float64
     pi      math.Pi    float64
     phi     math.Phi   float64
     
     sqrt2     math.Sqrt2   float64
     sqrte     math.SqrtE   float64
     sqrtpi    math.SqrtPi  float64
     sqrtphi   math.SqrtPhi float64
     
     ln2     math.Ln2    float64
     log2e   math.Log2E  float64
     ln10    math.Ln10   float64
     log10e  math.Log10E float64
          
     true    true       bool
     false   false      bool
     null    nil        interface{}

### Supported operations

Package has several predefined operators. You are free to add new one with `AddOperator`

[Operator precedence](https://golang.org/ref/spec#Operator_precedence)

	Precedence    Operator
	    6	    	  **
	    5             *  /  %  <<  >>  &  &^
	    4             +  -  |  ^
	    3             ==  !=  <  <=  >  >=
	    2             &&
	    1             ||

[Arithmetic operators](https://golang.org/ref/spec#Arithmetic_operators)

	**   power                  integers, floats
	+    sum                    integers, floats, strings
	-    difference             integers, floats
	*    product                integers, floats
	/    quotient               integers, floats
	%    remainder              integers

	&    bitwise AND            integers
	|    bitwise OR             integers
	^    bitwise XOR            integers
	&^   bit clear (AND NOT)    integers

	<<   left shift             integer << unsigned integer
	>>   right shift            integer >> unsigned integer

### Supported functions

Package has several predefined functions. You are free to add new one with `AddFunction`

    abs          math.Abs          integers, floats
    acos         math.Acos         integers, floats
    acosh        math.Acosh        integers, floats
    asin         math.Asin         integers, floats
    asinh        math.Asinh        integers, floats
    atan         math.Atan         integers, floats
    atanh        math.Atanh        integers, floats
    avg          Average           array of integers or floats
    cbrt         math.Cbrt         integers, floats
    ceil         math.Ceil         integers, floats
    cos          math.Cos          integers, floats
    cosh         math.Cosh         integers, floats
    erf          math.Erf          integers, floats
    erfc         math.Erfc         integers, floats
    erfcinv      math.Erfcinv      integers, floats
    erfinv       math.Erfinv       integers, floats
    exp          math.Exp          integers, floats
    exp2         math.Exp2         integers, floats
    expm1        math.Expm1        integers, floats
    factorial    N!                unsigned integer
    floor        math.Floor        integers, floats
    gamma        math.Gamma        integers, floats
    j0           math.J0           integers, floats
    j1           math.J1           integers, floats
    length       len               array
    log          math.Log          integers, floats
    log10        math.Log10        integers, floats
    log1p        math.Log1p        integers, floats
    log2         math.Log2         integers, floats
    logb         math.Logb         integers, floats
    pow10        math.Pow10        integer
    round        math.Round        integers, floats
    roundtoeven  math.RoundToEven  integers, floats
    sin          math.Sin          integers, floats
    sinh         math.Sinh         integers, floats
    sqrt         math.Sqrt         integers, floats
    tan          math.Tan          integers, floats
    tanh         math.Tanh         integers, floats
    trunc        math.Trunc        integers, floats
    y0           math.Y0           integers, floats
    y1           math.Y1           integers, floats

# Benchmarks

Current package is comparable with `encoding/json` package. 

Test data:
```json
{ "store": {
    "book": [ 
      { "category": "reference",
        "author": "Nigel Rees",
        "title": "Sayings of the Century",
        "price": 8.95
      },
      { "category": "fiction",
        "author": "Evelyn Waugh",
        "title": "Sword of Honour",
        "price": 12.99
      },
      { "category": "fiction",
        "author": "Herman Melville",
        "title": "Moby Dick",
        "isbn": "0-553-21311-3",
        "price": 8.99
      },
      { "category": "fiction",
        "author": "J. R. R. Tolkien",
        "title": "The Lord of the Rings",
        "isbn": "0-395-19395-8",
        "price": 22.99
      }
    ],
    "bicycle": {
      "color": "red",
      "price": 19.95
    }
  }
}
```
JSONPath: `$.store..price`

```
$ go test -bench=. -cpu=1 -benchmem
goos: linux
goarch: amd64
pkg: github.com/spyzhov/ajson
BenchmarkUnmarshal_AJSON           87807             13668 ns/op            5344 B/op         95 allocs/op
BenchmarkUnmarshal_JSON            64012             17560 ns/op             968 B/op         31 allocs/op
BenchmarkJSONPath_all_prices       56743             26049 ns/op            7368 B/op        161 allocs/op
```
