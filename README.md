# Abstract JSON 

[![Build Status](https://travis-ci.com/spyzhov/ajson.svg?token=swf7VyTzTWuHyiC9QzT4&branch=master)](https://travis-ci.com/spyzhov/ajson)
[![Go Report Card](https://goreportcard.com/badge/github.com/spyzhov/ajson)](https://goreportcard.com/report/github.com/spyzhov/ajson)
[![GoDoc](https://godoc.org/github.com/spyzhov/ajson?status.svg)](https://godoc.org/github.com/spyzhov/ajson)

Abstract JSON is a small golang package that provide a parser for JSON, in case when you are not sure in it's structure.

Method `Unmarshal` will scan all the byte slice to create a root node of JSON structure, with all it behaviors.

Each `Node` has it's own type and calculated value, which will be calculated on demand. 
Calculated value saves in `atomic.Value`, so it's thread safe.

## Example

Calculating `AVG(price)` when object is heterogeneous.

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

```
$ go test -bench=. -cpu=1 -benchmem
goos: linux
goarch: amd64
pkg: github.com/spyzhov/ajson
BenchmarkUnmarshal_AJSON          200000              6139 ns/op            5592 B/op         96 allocs/op
BenchmarkUnmarshal_JSON           200000             10264 ns/op             840 B/op         28 allocs/op
```

# TODO

- Functions 
- [ ] `func JsonPath(data [] byte, path string) ([]*Node, error)` 
- [ ] `func (n *Node) JsonPath(path string) ([]*Node, error)`
- [ ] `func Validate(data [] byte, path string) error`
- node
- [ ] add `atomic.Value` for `Path()`
- [ ] add `atomic.Value` for `Key()`, remove preparse key value
- buffer
- [ ] add tests
- errors
- [ ] expected error: `wrong symbol '%s' expected %s, on %d`
- [ ] add buffer in error: detect column and line from index
- future
- [ ] use io.Reader instead of []byte
- refactoring
- [ ] try to remove node.borders
- [ ] remove reflection in node.inheritors
- FixMe:
- [ ] backslash system symbols in JsonPath
- [ ] evaluate... Result as key or filter
- [ ] array index before zero
- [x] slices like [:3], [-1:]
- [ ] exists `‌$..book[?(@.isbn)]`
- [ ] ‌rpn: only operations in stack (validation)
- [x] array length as calculated node
- [ ] ‌math: round, ceil, floor, exp, log, ln, sin, cos, tan, ctg,... Const: pi, e
- [ ] check, what to do with an argument functions like `round(value, n)`?
- [ ] const as *Node & function as func(*Node) (*Node, error)
- [ ] ‌eval, Dereference : deref(Node, path):[] *Node , eval (node, cmd): [] *Node. 
	Where: deref - parse path, eval each route and process all subresults. 
	Eval - take a *Node and command and process it.