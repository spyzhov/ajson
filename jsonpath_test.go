package ajson

import (
	"strings"
	"testing"
)

// JSON from example https://goessner.net/articles/JsonPath/index.html#e3
var jsonpathTestData = []byte(`{ "store": {
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
}`)

func fullPath(array []*Node) string {
	result := make([]string, 0, len(array))
	for _, element := range array {
		result = append(result, element.Path())
	}
	return "[" + strings.Join(result, ", \n") + "]"
}

func TestJsonPath_simple(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		count int
	}{
		{name: "root", path: "$", count: 1},
		{name: "roots", path: "$.", count: 1},
		{name: "all objects", path: "$..", count: 8},
		{name: "only children", path: "$.*", count: 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := JsonPath(jsonpathTestData, test.path)
			if err != nil {
				t.Errorf("Error on JsonPath(json, %s) as %s: %s", test.path, test.name, err.Error())
			} else if len(result) != test.count {
				t.Errorf("Error on JsonPath(json, %s) as %s: length must be %d found %d\n%s", test.path, test.name, test.count, len(result), fullPath(result))
			}
		})
	}
}
