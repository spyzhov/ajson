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
	return sliceString(result)
}

func sliceString(array []string) string {
	return "[" + strings.Join(array, ", ") + "]"
}

func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestJsonPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{name: "root", path: "$", expected: "[$]"},
		{name: "roots", path: "$.", expected: "[$]"},
		{name: "all objects", path: "$..", expected: "[$, $['store'], $['store']['bicycle'], $['store']['book'], $['store']['book'][0], $['store']['book'][1], $['store']['book'][2], $['store']['book'][3]]"},
		{name: "only children", path: "$.*", expected: "[$['store']]"},

		{name: "by key", path: "$.store.bicycle", expected: "[$['store']['bicycle']]"},
		{name: "all key 1", path: "$..bicycle", expected: "[$['store']['bicycle']]"},
		{name: "all key 2", path: "$..price", expected: "[$['store']['bicycle']['price'], $['store']['book'][0]['price'], $['store']['book'][1]['price'], $['store']['book'][2]['price'], $['store']['book'][3]['price']]"},
		{name: "all key bracket", path: "$..['price']", expected: "[$['store']['bicycle']['price'], $['store']['book'][0]['price'], $['store']['book'][1]['price'], $['store']['book'][2]['price'], $['store']['book'][3]['price']]"},
		{name: "all fields", path: "$['store']['book'][1].*", expected: "[$['store']['book'][1]['author'], $['store']['book'][1]['category'], $['store']['book'][1]['price'], $['store']['book'][1]['title']]"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := JSONPath(jsonpathTestData, test.path)
			if err != nil {
				t.Errorf("Error on JsonPath(json, %s) as %s: %s", test.path, test.name, err.Error())
			} else if fullPath(result) != test.expected {
				t.Errorf("Error on JsonPath(json, %s) as %s: path doesn't match\nExpected: %s\nActual:   %s", test.path, test.name, test.expected, fullPath(result))
			}
		})
	}
}

func TestParseJSONPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected []string
	}{
		{name: "root", path: "$", expected: []string{"$"}},
		{name: "roots", path: "$.", expected: []string{"$"}},
		{name: "all objects", path: "$..", expected: []string{"$", ".."}},
		{name: "only children", path: "$.*", expected: []string{"$", "*"}},
		{name: "all objects children", path: "$..*", expected: []string{"$", "..", "*"}},
		{name: "path dot:simple", path: "$.root.element", expected: []string{"$", "root", "element"}},
		{name: "path dot:combined", path: "$.root.*.element", expected: []string{"$", "root", "*", "element"}},
		{name: "path bracket:simple", path: "$['root']['element']", expected: []string{"$", "root", "element"}},
		{name: "path bracket:combined", path: "$['root'][*]['element']", expected: []string{"$", "root", "*", "element"}},
		{name: "path bracket:int", path: "$['store']['book'][0]['title']", expected: []string{"$", "store", "book", "0", "title"}},
		{name: "path combined:simple", path: "$['root'].*['element']", expected: []string{"$", "root", "*", "element"}},
		{name: "path combined:dotted", path: "$.['root'].*.['element']", expected: []string{"$", "root", "*", "element"}},
		{name: "path combined:dotted small", path: "$['root'].*.['element']", expected: []string{"$", "root", "*", "element"}},
		{name: "phoneNumbers", path: "$.phoneNumbers[*].type", expected: []string{"$", "phoneNumbers", "*", "type"}},
		{name: "filtered", path: "$.store.book[?(@.price < 10)].title", expected: []string{"$", "store", "book", "?(@.price < 10)", "title"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := ParseJSONPath(test.path)
			if err != nil {
				t.Errorf("Error on parseJsonPath(json, %s) as %s: %s", test.path, test.name, err.Error())
			} else if !sliceEqual(result, test.expected) {
				t.Errorf("Error on parseJsonPath(%s) as %s: path doesn't match\nExpected: %s\nActual: %s", test.path, test.name, sliceString(test.expected), sliceString(result))
			}
		})
	}
}
