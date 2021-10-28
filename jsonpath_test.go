package ajson

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// JSON from example https://goessner.net/articles/JsonPath/index.html#e3
var jsonPathTestData = []byte(`{ "store": {
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

var jpStubs = map[string]string{
	// JSON from API https://randomuser.me/api/?results=1
	"random_user": `{
  "results": [
    {
      "gender": "female",
      "name": {
        "title": "Miss",
        "first": "Annette",
        "last": "Jennings"
      },
      "location": {
        "street": {
          "number": 855,
          "name": "Dane St"
        },
        "city": "Abilene",
        "state": "New York",
        "country": "United States",
        "postcode": 90538,
        "coordinates": {
          "latitude": "88.1096",
          "longitude": "22.9540"
        },
        "timezone": {
          "offset": "-1:00",
          "description": "Azores, Cape Verde Islands"
        }
      },
      "email": "annette.jennings@example.com",
      "login": {
        "uuid": "0726120a-e330-42f6-821c-fcec841b797a",
        "username": "smallwolf778",
        "password": "banker",
        "salt": "modk0zGi",
        "md5": "5b3a522f1e66625d0e76b92f031ffe80",
        "sha1": "d56c12fa9110585956392523a90b06aca99a7fc9",
        "sha256": "1b5a15c21683337d60d53dfb5853c8357fa218a96e280f2adf5aefdf2157fd76"
      },
      "dob": {
        "date": "1994-09-14T00:34:29.287Z",
        "age": 26
      },
      "registered": {
        "date": "2015-11-18T14:09:34.475Z",
        "age": 5
      },
      "phone": "(217)-464-6621",
      "cell": "(445)-236-5456",
      "id": {
        "name": "SSN",
        "value": "071-39-5493"
      },
      "picture": {
        "large": "https://randomuser.me/api/portraits/women/48.jpg",
        "medium": "https://randomuser.me/api/portraits/med/women/48.jpg",
        "thumbnail": "https://randomuser.me/api/portraits/thumb/women/48.jpg"
      },
      "nat": "US"
    }
  ],
  "info": {
    "seed": "eabbc063d314ff54",
    "results": 1,
    "page": 1,
    "version": "1.3"
  }
}`,
}

func fullPath(array []*Node) string {
	return sliceString(Paths(array))
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
		wantErr  bool
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

		{name: "union fields", path: "$['store']['book'][2]['author','price','title']", expected: "[$['store']['book'][2]['author'], $['store']['book'][2]['price'], $['store']['book'][2]['title']]"},
		{name: "union indexes", path: "$['store']['book'][1,2]", expected: "[$['store']['book'][1], $['store']['book'][2]]"},
		{name: "union indexes calculate", path: "$['store']['book'][-2,(@.length-1)]", expected: "[$['store']['book'][2], $['store']['book'][3]]"},
		{name: "union indexes position", path: "$['store']['book'][-1,-3]", expected: "[$['store']['book'][3], $['store']['book'][1]]"},

		{name: "slices 1", path: "$..[1:4]", expected: "[$['store']['book'][1], $['store']['book'][2], $['store']['book'][3]]"},
		{name: "slices 2", path: "$..[1:4:]", expected: "[$['store']['book'][1], $['store']['book'][2], $['store']['book'][3]]"},
		{name: "slices 3", path: "$..[1:4:1]", expected: "[$['store']['book'][1], $['store']['book'][2], $['store']['book'][3]]"},
		{name: "slices 4", path: "$..[1:]", expected: "[$['store']['book'][1], $['store']['book'][2], $['store']['book'][3]]"},
		{name: "slices 5", path: "$..[:2]", expected: "[$['store']['book'][0], $['store']['book'][1]]"},
		{name: "slices 6", path: "$..[:4:2]", expected: "[$['store']['book'][0], $['store']['book'][2]]"},
		{name: "slices 7", path: "$..[:4:]", expected: "[$['store']['book'][0], $['store']['book'][1], $['store']['book'][2], $['store']['book'][3]]"},
		{name: "slices 8", path: "$..[::]", expected: "[$['store']['book'][0], $['store']['book'][1], $['store']['book'][2], $['store']['book'][3]]"},
		{name: "slices 9", path: "$['store']['book'][1:4:2]", expected: "[$['store']['book'][1], $['store']['book'][3]]"},
		{name: "slices 10", path: "$['store']['book'][1:4:3]", expected: "[$['store']['book'][1]]"},
		{name: "slices 11", path: "$['store']['book'][:-1]", expected: "[$['store']['book'][0], $['store']['book'][1], $['store']['book'][2]]"},
		{name: "slices 12", path: "$['store']['book'][-1:]", expected: "[$['store']['book'][3]]"},
		{name: "slices 13", path: "$..[::-1]", expected: "[$['store']['book'][3], $['store']['book'][2], $['store']['book'][1], $['store']['book'][0]]"},
		{name: "slices 14", path: "$..[::-2]", expected: "[$['store']['book'][3], $['store']['book'][1]]"},
		{name: "slices 15", path: "$..[::2]", expected: "[$['store']['book'][0], $['store']['book'][2]]"},
		{name: "slices 16", path: "$..[-3:(@.length)]", expected: "[$['store']['book'][1], $['store']['book'][2], $['store']['book'][3]]"},
		{name: "slices 17", path: "$..[1:(@.length - 1)]", expected: "[$['store']['book'][1], $['store']['book'][2]]"},
		{name: "slices 18", path: "$..[(foobar(@.length))::]", wantErr: true},
		{name: "slices 19", path: "$..[::0]", wantErr: true},
		{name: "slices 20", path: "$..[:(1/0):]", wantErr: true},
		{name: "slices 21", path: "$..[:(1/2):]", wantErr: true},
		{name: "slices 22", path: "$..[:0.5:]", wantErr: true},

		{name: "calculated 1", path: "$['store']['book'][(@.length-1)]", expected: "[$['store']['book'][3]]"},
		{name: "calculated 2", path: "$['store']['book'][(3.5 - 3/2)]", expected: "[$['store']['book'][2]]"},
		{name: "calculated 3", path: "$..book[?(@.isbn)]", expected: "[$['store']['book'][2], $['store']['book'][3]]"},
		{name: "calculated 4", path: "$..[?(@.price < factorial(3) + 3)]", expected: "[$['store']['book'][0], $['store']['book'][2]]"},
		{name: "calculated 5", path: "$..[(1/0)]", wantErr: true},
		{name: "calculated 6", path: "$[('store')][('bo'+'ok')][(@.length - 1)]", expected: "[$['store']['book'][3]]"},
		{name: "calculated 7", path: "$[('store'+'')][('bo'+'ok')][(true)]", expected: "[$['store']['book'][0], $['store']['book'][1], $['store']['book'][2], $['store']['book'][3]]"},
		{name: "calculated 8", path: "$.store.book[(@.length / 0)]", wantErr: true},
		{name: "calculated 9", path: "$.store.book[?(@.price / 0 > 0)]", wantErr: true},
		{name: "calculated 10", path: "$.store.bicycle.price[(@.length-1)]", expected: `[]`},
		{name: "calculated 11", path: "$.store.bicycle.price[?(@ > 0)]", expected: `[]`},
		{name: "calculated 12", path: "$.store.book[?(@.price * 0 = 0)]", wantErr: true},

		{name: "$.store.book[*].author", path: "$.store.book[*].author", expected: "[$['store']['book'][0]['author'], $['store']['book'][1]['author'], $['store']['book'][2]['author'], $['store']['book'][3]['author']]"},
		{name: "$..author", path: "$..author", expected: "[$['store']['book'][0]['author'], $['store']['book'][1]['author'], $['store']['book'][2]['author'], $['store']['book'][3]['author']]"},
		{name: "$.store..price", path: "$.store..price", expected: "[$['store']['bicycle']['price'], $['store']['book'][0]['price'], $['store']['book'][1]['price'], $['store']['book'][2]['price'], $['store']['book'][3]['price']]"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := JSONPath(jsonPathTestData, test.path)
			if (err != nil) != test.wantErr {
				t.Errorf("JSONPath() error = %v, wantErr %v. got = %v", err, test.wantErr, result)
				return
			}
			if test.wantErr {
				return
			}
			if fullPath(result) != test.expected {
				t.Errorf("Error on JsonPath(json, %s) as %s: path doesn't match\nExpected: %s\nActual:   %s", test.path, test.name, test.expected, fullPath(result))
			}
		})
	}
}

func TestJsonPath_value(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected interface{}
	}{
		{name: "length", path: "$['store']['book'].length", expected: float64(4)},
		{name: "price", path: "$['store']['book'][?(@.price + 0.05 == 9)].price", expected: float64(8.95)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := JSONPath(jsonPathTestData, test.path)
			if err != nil {
				t.Errorf("Error on JsonPath(json, %s) as %s: %s", test.path, test.name, err.Error())
			} else if len(result) != 1 {
				t.Errorf("Error on JsonPath(json, %s) as %s: path to long, expected only value\nActual: %s", test.path, test.name, fullPath(result))
			} else {
				val, err := result[0].Value()
				if err != nil {
					t.Errorf("Error on JsonPath(json, %s): error %s", test.path, err.Error())
				} else {
					switch {
					case result[0].IsNumeric():
						if val.(float64) != test.expected.(float64) {
							t.Errorf("Error on JsonPath(json, %s): value doesn't match\nExpected: %v\nActual:   %v", test.path, test.expected, val)
						}
					case result[0].IsString():
						if val.(string) != test.expected.(string) {
							t.Errorf("Error on JsonPath(json, %s): value doesn't match\nExpected: %v\nActual:   %v", test.path, test.expected, val)
						}
					case result[0].IsBool():
						if val.(bool) != test.expected.(bool) {
							t.Errorf("Error on JsonPath(json, %s): value doesn't match\nExpected: %v\nActual:   %v", test.path, test.expected, val)
						}
					default:
						t.Errorf("Error on JsonPath(json, %s): unsupported type found", test.path)
					}
				}
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
		{name: "path bracket:simple", path: "$['root']['element']", expected: []string{"$", "'root'", "'element'"}},
		{name: "path bracket:combined", path: "$['root'][*]['element']", expected: []string{"$", "'root'", "*", "'element'"}},
		{name: "path bracket:int", path: "$['store']['book'][0]['title']", expected: []string{"$", "'store'", "'book'", "0", "'title'"}},
		{name: "path combined:simple", path: "$['root'].*['element']", expected: []string{"$", "'root'", "*", "'element'"}},
		{name: "path combined:dotted", path: "$.['root'].*.['element']", expected: []string{"$", "'root'", "*", "'element'"}},
		{name: "path combined:dotted small", path: "$['root'].*.['element']", expected: []string{"$", "'root'", "*", "'element'"}},
		{name: "phoneNumbers", path: "$.phoneNumbers[*].type", expected: []string{"$", "phoneNumbers", "*", "type"}},
		{name: "filtered", path: "$.store.book[?(@.price < 10)].title", expected: []string{"$", "store", "book", "?(@.price < 10)", "title"}},
		{name: "formula", path: "$..phoneNumbers..('ty' + 'pe')", expected: []string{"$", "..", "phoneNumbers", "..", "('ty' + 'pe')"}},
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

// Test suites from cburgmer/json-path-comparison
func TestJSONPath_suite(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		path     string
		expected []interface{}
		wantErr  bool
	}{
		{
			name:     "Bracket notation with double quotes",
			input:    `{"key": "value"}`,
			path:     `$["key"]`,
			expected: []interface{}{"value"}, // ["value"]
		},
		{
			name:     "Filter expression with bracket notation",
			input:    `[{"key": 0}, {"key": 42}, {"key": -1}, {"key": 41}, {"key": 43}, {"key": 42.0001}, {"key": 41.9999}, {"key": 100}, {"some": "value"}]`,
			path:     `$[?(@['key']==42)]`,
			expected: []interface{}{map[string]interface{}{"key": float64(42)}}, // [{"key": 42}]
		},
		{
			name:     "Filter expression with equals string with dot literal",
			input:    `[{"key": "some"}, {"key": "value"}, {"key": "some.value"}]`,
			path:     `$[?(@.key=="some.value")]`,
			expected: []interface{}{map[string]interface{}{"key": "some.value"}}, // [{"key": "some.value"}]
		},
		{
			name:     "Array slice with negative step only",
			input:    `["first", "second", "third", "forth", "fifth"]`,
			path:     `$[::-2]`,
			expected: []interface{}{"fifth", "third", "first"}, // ["fifth", "third", "first"]
		},
		{
			name:     "Filter expression with bracket notation with -1",
			input:    `[[2, 3], ["a"], [0, 2], [2]]`,
			path:     `$[?(@[-1]==2)]`,
			expected: []interface{}{[]interface{}{float64(0), float64(2)}, []interface{}{float64(2)}}, // [[0, 2], [2]]
		},
		{
			name:     "Filter expression with bracket notation with number",
			input:    `[["a", "b"], ["x", "y"]]`,
			path:     `$[?(@[1]=='b')]`,
			expected: []interface{}{[]interface{}{"a", "b"}}, // [["a", "b"]]
		},
		{
			name:     "Filter expression with equals string with current object literal",
			input:    `[{"key": "some"}, {"key": "value"}, {"key": "hi@example.com"}]`,
			path:     `$[?(@.key=="hi@example.com")]`,
			expected: []interface{}{map[string]interface{}{"key": "hi@example.com"}}, // [{"key": "hi@example.com"}]
		},
		// 		{
		// 			name:     "Filter expression with negation and equals",
		// 			input:    `[
		//     {"key": 0},
		//     {"key": 42},
		//     {"key": -1},
		//     {"key": 41},
		//     {"key": 43},
		//     {"key": 42.0001},
		//     {"key": 41.9999},
		//     {"key": 100},
		//     {"key": "43"},
		//     {"key": "42"},
		//     {"key": "41"},
		//     {"key": "value"},
		//     {"some": "value"}
		// ]`,
		// 			path:     `$[?(!(@.key==42))]`,
		// 			expected: []interface{}{
		// 				map[string]interface{}{"key": float64(0)},
		// 				map[string]interface{}{"key": float64(-1)},
		// 				map[string]interface{}{"key": float64(41)},
		// 				map[string]interface{}{"key": float64(43)},
		// 				map[string]interface{}{"key": float64(42.0001)},
		// 				map[string]interface{}{"key": float64(41.9999)},
		// 				map[string]interface{}{"key": float64(100)},
		// 				map[string]interface{}{"key": "43"},
		// 				map[string]interface{}{"key": "42"},
		// 				map[string]interface{}{"key": "41"},
		// 				map[string]interface{}{"key": "value"},
		// 				map[string]interface{}{"some": "value"},
		// 			},
		// 		},
		{
			name:     "Filter expression with bracket notation with number on object",
			input:    `{"1": ["a", "b"], "2": ["x", "y"]}`,
			path:     `$[?(@[1]=='b')]`,
			expected: []interface{}{[]interface{}{"a", "b"}}, // [["a", "b"]]
		},
		// {
		// 	name:     "Dot notation with single quotes and dot",
		// 	input:    `{"some.key": 42, "some": {"key": "value"}}`,
		// 	path:     `$.'some.key'`,
		// 	expected: []interface{}{float64(42)}, // [42]
		// },
		{
			name:    "Array slice with step 0",
			input:   `["first", "second", "third", "forth", "fifth"]`,
			path:    `$[0:3:0]`,
			wantErr: true,
		},
		{
			name:     "$[2:1]",
			input:    `["first", "second", "third", "forth"]`,
			path:     `$[2:1]`,
			expected: []interface{}{}, // []
		},
		{
			name:     "$[-4:]",
			input:    `["first", "second", "third"]`,
			path:     `$[-4:]`,
			expected: []interface{}{"first", "second", "third"}, // ["first", "second", "third"]
		},
		{
			name:     "$[']']",
			input:    `{"]": 42}`,
			path:     `$[']']`,
			expected: []interface{}{float64(42)}, // [42]
		},
		{
			name:     `$['"']`,
			input:    `{"\"": "value", "another": "entry"}`,
			path:     `$['"']`,
			expected: []interface{}{"value"}, // ["value"]
		},
		{
			name:     "$[2:113667776004]",
			input:    `["first", "second", "third", "forth", "fifth"]`,
			path:     `$[2:113667776004]`,
			expected: []interface{}{"third", "forth", "fifth"}, // ["third", "forth", "fifth"]
		},
		{
			name:     "$[2:-113667776004:-1]",
			input:    `["first", "second", "third", "forth", "fifth"]`,
			path:     `$[2:-113667776004:-1]`,
			expected: []interface{}{"third", "second", "first"}, // ["third", "second", "first"]
		},
		{
			name:     "$[-113667776004:2]",
			input:    `["first", "second", "third", "forth", "fifth"]`,
			path:     `$[-113667776004:2]`,
			expected: []interface{}{"first", "second"}, // ["first", "second"]
		},
		{
			name:     "$[113667776004:2:-1]",
			input:    `["first", "second", "third", "forth", "fifth"]`,
			path:     `$[113667776004:2:-1]`,
			expected: []interface{}{"fifth", "forth"}, // ["fifth", "forth"]
		},
		{
			name:     "$.length",
			input:    `[4, 5, 6]`,
			path:     `$.length`,
			expected: []interface{}{float64(3)}, // [3]
		},
		{
			name:    "$[?()]",
			input:   `[1, {"key": 42}, "value", null]`,
			path:    `$[?()]`,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nodes, err := JSONPath([]byte(test.input), test.path)
			if (err != nil) != test.wantErr {
				t.Errorf("JSONPath() error = %v, wantErr %v. got = %v", err, test.wantErr, nodes)
				return
			}
			if test.wantErr {
				return
			}

			results := make([]interface{}, 0)
			for _, node := range nodes {
				value, err := node.Unpack()
				if err != nil {
					t.Errorf("node.Unpack(): unexpected error: %v", err)
					return
				}
				results = append(results, value)
			}

			if !reflect.DeepEqual(results, test.expected) {
				t.Errorf("JSONPath(): wrong result:\nExpected: %#+v\nActual:   %#+v", test.expected, results)
			}
		})
	}
}

func ExampleJSONPath() {
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
    "bicycle": {
      "color": "red",
      "price": 19.95
    }
  }
}`)
	authors, err := JSONPath(json, "$.store.book[*].author")
	if err != nil {
		panic(err)
	}
	for _, author := range authors {
		fmt.Println(author.MustString())
	}
	// Output:
	// Nigel Rees
	// Evelyn Waugh
	// Herman Melville
	// J. R. R. Tolkien
}

func ExampleJSONPath_array() {
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
    "bicycle": {
      "color": "red",
      "price": 19.95
    }
  }
}`)
	authors, err := JSONPath(json, "$.store.book[*].author")
	if err != nil {
		panic(err)
	}
	result, err := Marshal(ArrayNode("", authors))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(result))
	// Output:
	// ["Nigel Rees","Evelyn Waugh","Herman Melville","J. R. R. Tolkien"]
}

func ExampleEval() {
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
	root, err := Unmarshal(json)
	if err != nil {
		panic(err)
	}
	result, err := Eval(root, "avg($..price)")
	if err != nil {
		panic(err)
	}
	fmt.Print(result.MustNumeric())
	// Output:
	// 14.774000000000001
}

func TestEval(t *testing.T) {
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
	tests := []struct {
		name     string
		root     *Node
		eval     string
		expected *Node
		wantErr  bool
	}{
		{
			name:     "avg($..price)",
			root:     Must(Unmarshal(json)),
			eval:     "avg($..price)",
			expected: NumericNode("", 14.774000000000001),
			wantErr:  false,
		},
		{
			name:     "avg($..price)",
			root:     Must(Unmarshal(json)),
			eval:     "avg()",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "avg($..price)",
			root:     Must(Unmarshal(json)),
			eval:     "($..price+)",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "round(avg($..price)+pi)",
			root:     Must(Unmarshal(json)),
			eval:     "round(avg($..price)+pi)",
			expected: NumericNode("", 18),
			wantErr:  false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Eval(test.root, test.eval)
			if (err != nil) != test.wantErr {
				t.Errorf("Eval() error = %v, wantErr %v. got = %v", err, test.wantErr, result)
				return
			}
			if test.wantErr {
				return
			}
			if result == nil {
				t.Errorf("Eval() result in nil")
				return
			}

			if ok, err := result.Eq(test.expected); !ok {
				t.Errorf("result.Eq(): wrong result:\nExpected: %#+v\nActual: %#+v", test.expected, result.value.Load())
			} else if err != nil {
				t.Errorf("result.Eq() error = %v", err)
			}

		})
	}
}

func BenchmarkJSONPath_all_prices(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		_, err = JSONPath(jsonPathTestData, "$.store..price")
		if err != nil {
			b.Error()
		}
	}
}

// https://github.com/cburgmer/json-path-comparison/blob/master/regression_suite/regression_suite.yaml
func TestJSONPath_comparison_consensus(t *testing.T) {
	tests := []struct {
		name      string
		selector  string
		document  string
		consensus string
	}{
		{
			name:      `array_slice`,
			selector:  `$[1:3]`,
			document:  `["first", "second", "third", "forth", "fifth"]`,
			consensus: `["second", "third"]`,
		},
		{
			name:      `array_slice_on_exact_match`,
			selector:  `$[0:5]`,
			document:  `["first", "second", "third", "forth", "fifth"]`,
			consensus: `["first", "second", "third", "forth", "fifth"]`,
		},
		{
			name:      `array_slice_on_non_overlapping_array`,
			selector:  `$[7:10]`,
			document:  `["first", "second", "third"]`,
			consensus: `[]`,
		},
		{
			name:      `array_slice_on_object`,
			selector:  `$[1:3]`,
			document:  `{":": 42, "more": "string", "a": 1, "b": 2, "c": 3}`,
			consensus: `[]`,
		},
		{
			name:      `array_slice_on_partially_overlapping_array`,
			selector:  `$[1:10]`,
			document:  `["first", "second", "third"]`,
			consensus: `["second", "third"]`,
		},
		{
			name:      `array_slice_with_open_end`,
			selector:  `$[1:]`,
			document:  `["first", "second", "third", "forth", "fifth"]`,
			consensus: `["second", "third", "forth", "fifth"]`,
		},
		{
			name:      `array_slice_with_open_start`,
			selector:  `$[:2]`,
			document:  `["first", "second", "third", "forth", "fifth"]`,
			consensus: `["first", "second"]`,
		},
		{
			name:      `array_slice_with_open_start_and_end`,
			selector:  `$[:]`,
			document:  `["first", "second"]`,
			consensus: `["first", "second"]`,
		},
		{
			name:      `array_slice_with_open_start_and_end_and_step_empty`,
			selector:  `$[::]`,
			document:  `["first", "second"]`,
			consensus: `["first", "second"]`,
		},
		{
			name:      `array_slice_with_range_of_-1`,
			selector:  `$[2:1]`,
			document:  `["first", "second", "third", "forth"]`,
			consensus: `[]`,
		},
		{
			name:      `array_slice_with_range_of_0`,
			selector:  `$[0:0]`,
			document:  `["first", "second"]`,
			consensus: `[]`,
		},
		{
			name:      `array_slice_with_range_of_1`,
			selector:  `$[0:1]`,
			document:  `["first", "second"]`,
			consensus: `["first"]`,
		},
		{
			name:      `array_slice_with_start_-1_and_open_end`,
			selector:  `$[-1:]`,
			document:  `["first", "second", "third"]`,
			consensus: `["third"]`,
		},
		{
			name:      `array_slice_with_start_-2_and_open_end`,
			selector:  `$[-2:]`,
			document:  `["first", "second", "third"]`,
			consensus: `["second", "third"]`,
		},
		{
			name:      `array_slice_with_start_large_negative_number_and_open_end_on_short_array`,
			selector:  `$[-4:]`,
			document:  `["first", "second", "third"]`,
			consensus: `["first", "second", "third"]`,
		},
		{
			name:      `array_slice_with_step`,
			selector:  `$[0:3:2]`,
			document:  `["first", "second", "third", "forth", "fifth"]`,
			consensus: `["first", "third"]`,
		},
		{
			name:      `array_slice_with_step_1`,
			selector:  `$[0:3:1]`,
			document:  `["first", "second", "third", "forth", "fifth"]`,
			consensus: `["first", "second", "third"]`,
		},
		{
			name:      `array_slice_with_step_and_leading_zeros`,
			selector:  `$[010:024:010]`,
			document:  `[0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25]`,
			consensus: `[10, 20]`,
		},
		{
			name:      `array_slice_with_step_but_end_not_aligned`,
			selector:  `$[0:4:2]`,
			document:  `["first", "second", "third", "forth", "fifth"]`,
			consensus: `["first", "third"]`,
		},
		{
			name:      `array_slice_with_step_empty`,
			selector:  `$[1:3:]`,
			document:  `["first", "second", "third", "forth", "fifth"]`,
			consensus: `["second", "third"]`,
		},
		{
			name:      `array_slice_with_step_only`,
			selector:  `$[::2]`,
			document:  `["first", "second", "third", "forth", "fifth"]`,
			consensus: `["first", "third", "fifth"]`,
		},
		{
			name:      `bracket_notation`,
			selector:  `$['key']`,
			document:  `{"key": "value"}`,
			consensus: `["value"]`,
		},
		{
			name:      `bracket_notation_after_recursive_descent`,
			selector:  `$..[0]`,
			document:  `["first", {"key": ["first nested", {"more": [{"nested": ["deepest", "second"]}, ["more", "values"]]}]}]`,
			consensus: `["first", "first nested", {"nested": ["deepest", "second"]}, "more", "deepest"]`,
			// consensus: `["deepest", "first nested", "first", "more", {"nested": ["deepest", "second"]}]`,
		},
		{
			name:      `bracket_notation_with_dot`,
			selector:  `$['two.some']`,
			document:  `{"one": {"key": "value"}, "two": {"some": "more", "key": "other value"}, "two.some": "42"}`,
			consensus: `["42"]`,
		},
		{
			name:      `bracket_notation_with_double_quotes`,
			selector:  `$["key"]`,
			document:  `{"key": "value"}`,
			consensus: `["value"]`,
		},
		{
			name:      `bracket_notation_with_empty_string`,
			selector:  `$['']`,
			document:  `{"": 42, "''": 123, "\"\"": 222}`,
			consensus: `[42]`,
		},
		{
			name:      `bracket_notation_with_number`,
			selector:  `$[2]`,
			document:  `["first", "second", "third", "forth", "fifth"]`,
			consensus: `["third"]`,
		},
		{
			name:      `bracket_notation_with_number_-1`,
			selector:  `$[-1]`,
			document:  `["first", "second", "third"]`,
			consensus: `["third"]`,
		},
		{
			name:      `bracket_notation_with_number_0`,
			selector:  `$[0]`,
			document:  `["first", "second", "third", "forth", "fifth"]`,
			consensus: `["first"]`,
		},
		{
			name:      `bracket_notation_with_number_after_dot_notation_with_wildcard_on_nested_arrays_with_different_length`,
			selector:  `$.*[1]`,
			document:  `[[1], [2, 3]]`,
			consensus: `[3]`,
		},
		{
			name:      `bracket_notation_with_quoted_array_slice_literal`,
			selector:  `$[':']`,
			document:  `{":": "value", "another": "entry"}`,
			consensus: `["value"]`,
		},
		{
			name:      `bracket_notation_with_quoted_closing_bracket_literal`,
			selector:  `$[']']`,
			document:  `{"]": 42}`,
			consensus: `[42]`,
		},
		{
			name:      `bracket_notation_with_quoted_current_object_literal`,
			selector:  `$['@']`,
			document:  `{"@": "value", "another": "entry"}`,
			consensus: `["value"]`,
		},
		{
			name:      `bracket_notation_with_quoted_dot_literal`,
			selector:  `$['.']`,
			document:  `{".": "value", "another": "entry"}`,
			consensus: `["value"]`,
		},
		{
			name:      `bracket_notation_with_quoted_dot_wildcard`,
			selector:  `$['.*']`,
			document:  `{"key": 42, ".*": 1, "": 10}`,
			consensus: `[1]`,
		},
		{
			name:      `bracket_notation_with_quoted_double_quote_literal`,
			selector:  `$['"']`,
			document:  `{"\"": "value", "another": "entry"}`,
			consensus: `["value"]`,
		},
		{
			name:      `bracket_notation_with_quoted_number_on_object`,
			selector:  `$['0']`,
			document:  `{"0": "value"}`,
			consensus: `["value"]`,
		},
		{
			name:      `bracket_notation_with_quoted_root_literal`,
			selector:  `$['$']`,
			document:  `{"$": "value", "another": "entry"}`,
			consensus: `["value"]`,
		},
		{
			name:      `bracket_notation_with_quoted_union_literal`,
			selector:  `$[',']`,
			document:  `{",": "value", "another": "entry"}`,
			consensus: `["value"]`,
		},
		{
			name:      `bracket_notation_with_quoted_wildcard_literal`,
			selector:  `$['*']`,
			document:  `{"*": "value", "another": "entry"}`,
			consensus: `["value"]`,
		},
		{
			name:      `bracket_notation_with_string_including_dot_wildcard`,
			selector:  `$['ni.*']`,
			document:  `{"nice": 42, "ni.*": 1, "mice": 100}`,
			consensus: `[1]`,
		},
		{
			name:      `bracket_notation_with_wildcard_after_array_slice`,
			selector:  `$[0:2][*]`,
			document:  `[[1, 2], ["a", "b"], [0, 0]]`,
			consensus: `[1, 2, "a", "b"]`,
		},
		{
			name:      `bracket_notation_with_wildcard_after_dot_notation_after_bracket_notation_with_wildcard`,
			selector:  `$[*].bar[*]`,
			document:  `[{"bar": [42]}]`,
			consensus: `[42]`,
		},
		{
			name:      `bracket_notation_with_wildcard_after_recursive_descent`,
			selector:  `$..[*]`,
			document:  `{"key": "value", "another key": {"complex": "string", "primitives": [0, 1]}}`,
			consensus: `[{"complex": "string", "primitives": [0, 1]}, "value", "string", [0, 1], 0, 1]`,
			// consensus: `["string", "value", 0, 1, [0, 1], {"complex": "string", "primitives": [0, 1]}]`,
		},
		{
			name:      `bracket_notation_with_wildcard_on_array`,
			selector:  `$[*]`,
			document:  `["string", 42, {"key": "value"}, [0, 1]]`,
			consensus: `["string", 42, {"key": "value"}, [0, 1]]`,
		},
		{
			name:      `bracket_notation_with_wildcard_on_empty_array`,
			selector:  `$[*]`,
			document:  `[]`,
			consensus: `[]`,
		},
		{
			name:      `bracket_notation_with_wildcard_on_empty_object`,
			selector:  `$[*]`,
			document:  `{}`,
			consensus: `[]`,
		},
		{
			name:      `bracket_notation_with_wildcard_on_null_value_array`,
			selector:  `$[*]`,
			document:  `[40, null, 42]`,
			consensus: `[40, null, 42]`,
		},
		{
			name:      `bracket_notation_with_wildcard_on_object`,
			selector:  `$[*]`,
			document:  `{"some": "string", "int": 42, "object": {"key": "value"}, "array": [0, 1]}`,
			consensus: `[[0, 1], 42, {"key": "value"}, "string"]`,
			// consensus: `["string", 42, [0, 1], {"key": "value"}]`,
		},
		{
			name:      `dot_notation`,
			selector:  `$.key`,
			document:  `{"key": "value"}`,
			consensus: `["value"]`,
		},
		{
			name:      `dot_notation_after_bracket_notation_with_wildcard`,
			selector:  `$[*].a`,
			document:  `[{"a": 1}, {"a": 1}]`,
			consensus: `[1, 1]`,
		},
		{
			name:      `dot_notation_after_bracket_notation_with_wildcard_on_one_matching`,
			selector:  `$[*].a`,
			document:  `[{"a": 1}]`,
			consensus: `[1]`,
		},
		{
			name:      `dot_notation_after_bracket_notation_with_wildcard_on_some_matching`,
			selector:  `$[*].a`,
			document:  `[{"a": 1}, {"b": 1}]`,
			consensus: `[1]`,
		},
		{
			name:      `dot_notation_after_filter_expression`,
			selector:  `$[?(@.id==42)].name`,
			document:  `[{"id": 42, "name": "forty-two"}, {"id": 1, "name": "one"}]`,
			consensus: `["forty-two"]`,
		},
		{
			name:      `dot_notation_after_recursive_descent`,
			selector:  `$..key`,
			document:  `{"object": {"key": "value", "array": [{"key": "something"}, {"key": {"key": "russian dolls"}}]}, "key": "top"}`,
			consensus: `["top", "value", "something", {"key": "russian dolls"}, "russian dolls"]`,
			// consensus: `["russian dolls", "something", "top", "value", {"key": "russian dolls"}]`,
		},
		{
			name:      `dot_notation_after_recursive_descent_after_dot_notation`,
			selector:  `$.store..price`,
			document:  `{"store": {"book": [{"category": "reference", "author": "Nigel Rees", "title": "Sayings of the Century", "price": 8.95}, {"category": "fiction", "author": "Evelyn Waugh", "title": "Sword of Honour", "price": 12.99}, {"category": "fiction", "author": "Herman Melville", "title": "Moby Dick", "isbn": "0-553-21311-3", "price": 8.99}, {"category": "fiction", "author": "J. R. R. Tolkien", "title": "The Lord of the Rings", "isbn": "0-395-19395-8", "price": 22.99}], "bicycle": {"color": "red", "price": 19.95}}}`,
			consensus: `[19.95, 8.95, 12.99, 8.99, 22.99]`,
			// consensus: `[12.99, 19.95, 22.99, 8.95, 8.99]`,
		},
		{
			name:      `dot_notation_after_union`,
			selector:  `$[0,2].key`,
			document:  `[{"key": "ey"}, {"key": "bee"}, {"key": "see"}]`,
			consensus: `["ey", "see"]`,
		},
		{
			name:      `dot_notation_after_union_with_keys`,
			selector:  `$['one','three'].key`,
			document:  `{"one": {"key": "value"}, "two": {"k": "v"}, "three": {"some": "more", "key": "other value"}}`,
			consensus: `["value", "other value"]`,
		},
		{
			name:      `dot_notation_on_array_value`,
			selector:  `$.key`,
			document:  `{"key": ["first", "second"]}`,
			consensus: `[["first", "second"]]`,
		},
		{
			name:      `dot_notation_on_empty_object_value`,
			selector:  `$.key`,
			document:  `{"key": {}}`,
			consensus: `[{}]`,
		},
		{
			name:      `dot_notation_on_null_value`,
			selector:  `$.key`,
			document:  `{"key": null}`,
			consensus: `[null]`,
		},
		{
			name:      `dot_notation_with_dash`,
			selector:  `$.key-dash`,
			document:  `{"key-dash": "value"}`,
			consensus: `["value"]`,
		},
		{
			name:      `dot_notation_with_key_named_in`,
			selector:  `$.in`,
			document:  `{"in": "value"}`,
			consensus: `["value"]`,
		},
		{
			name:      `dot_notation_with_key_named_null`,
			selector:  `$.null`,
			document:  `{"null": "value"}`,
			consensus: `["value"]`,
		},
		{
			name:      `dot_notation_with_key_named_true`,
			selector:  `$.true`,
			document:  `{"true": "value"}`,
			consensus: `["value"]`,
		},
		{
			name:      `dot_notation_with_non_ASCII_key`,
			selector:  `$.屬性`,
			document:  `{"\u5c6c\u6027": "value"}`,
			consensus: `["value"]`,
		},
		{
			name:      `dot_notation_with_number_on_object`,
			selector:  `$.2`,
			document:  `{"a": "first", "2": "second", "b": "third"}`,
			consensus: `["second"]`,
		},
		{
			name:      `dot_notation_with_wildcard_after_dot_notation_after_dot_notation_with_wildcard`,
			selector:  `$.*.bar.*`,
			document:  `[{"bar": [42]}]`,
			consensus: `[42]`,
		},
		{
			name:      `dot_notation_with_wildcard_after_dot_notation_with_wildcard_on_nested_arrays`,
			selector:  `$.*.*`,
			document:  `[[1, 2, 3], [4, 5, 6]]`,
			consensus: `[1, 2, 3, 4, 5, 6]`,
		},
		{
			name:      `dot_notation_with_wildcard_after_recursive_descent`,
			selector:  `$..*`,
			document:  `{"key": "value", "another key": {"complex": "string", "primitives": [0, 1]}}`,
			consensus: `[{"complex": "string", "primitives": [0, 1]}, "value", "string", [0, 1], 0, 1]`,
			// consensus: `["string", "value", 0, 1, [0, 1], {"complex": "string", "primitives": [0, 1]}]`,
		},
		{
			name:      `dot_notation_with_wildcard_after_recursive_descent_on_null_value_array`,
			selector:  `$..*`,
			document:  `[40, null, 42]`,
			consensus: `[40, null, 42]`,
			// consensus: `[40, 42, null]`,
		},
		{
			name:      `dot_notation_with_wildcard_after_recursive_descent_on_scalar`,
			selector:  `$..*`,
			document:  `42`,
			consensus: `[]`,
		},
		{
			name:      `dot_notation_with_wildcard_on_array`,
			selector:  `$.*`,
			document:  `["string", 42, {"key": "value"}, [0, 1]]`,
			consensus: `["string", 42, {"key": "value"}, [0, 1]]`,
		},
		{
			name:      `dot_notation_with_wildcard_on_empty_array`,
			selector:  `$.*`,
			document:  `[]`,
			consensus: `[]`,
		},
		{
			name:      `dot_notation_with_wildcard_on_empty_object`,
			selector:  `$.*`,
			document:  `{}`,
			consensus: `[]`,
		},
		{
			name:      `dot_notation_with_wildcard_on_object`,
			selector:  `$.*`,
			document:  `{"some": "string", "int": 42, "object": {"key": "value"}, "array": [0, 1]}`,
			consensus: `[[0, 1], 42, {"key": "value"}, "string"]`,
			// consensus: `["string", 42, [0, 1], {"key": "value"}]`,
		},
		{
			name:      `filter_expression_with_bracket_notation`,
			selector:  `$[?(@['key']==42)]`,
			document:  `[{"key": 0}, {"key": 42}, {"key": -1}, {"key": 41}, {"key": 43}, {"key": 42.0001}, {"key": 41.9999}, {"key": 100}, {"some": "value"}]`,
			consensus: `[{"key": 42}]`,
		},
		{
			name:      `filter_expression_with_bracket_notation_and_current_object_literal`,
			selector:  `$[?(@['@key']==42)]`,
			document:  `[{"@key": 0}, {"@key": 42}, {"key": 42}, {"@key": 43}, {"some": "value"}]`,
			consensus: `[{"@key": 42}]`,
		},
		{
			name:      `filter_expression_with_bracket_notation_with_number`,
			selector:  `$[?(@[1]=='b')]`,
			document:  `[["a", "b"], ["x", "y"]]`,
			consensus: `[["a", "b"]]`,
		},
		{
			name:      `filter_expression_with_equals_on_array_without_match`,
			selector:  `$[?(@.key==43)]`,
			document:  `[{"key": 42}]`,
			consensus: `[]`,
		},
		{
			name:      `filter_expression_with_equals_string_with_current_object_literal`,
			selector:  `$[?(@.key=="hi@example.com")]`,
			document:  `[{"key": "some"}, {"key": "value"}, {"key": "hi@example.com"}]`,
			consensus: `[{"key": "hi@example.com"}]`,
		},
		{
			name:      `filter_expression_with_equals_string_with_dot_literal`,
			selector:  `$[?(@.key=="some.value")]`,
			document:  `[{"key": "some"}, {"key": "value"}, {"key": "some.value"}]`,
			consensus: `[{"key": "some.value"}]`,
		},
		{
			name:      `filter_expression_with_equals_string_with_single_quotes`,
			selector:  `$[?(@.key=='value')]`,
			document:  `[{"key": "some"}, {"key": "value"}]`,
			consensus: `[{"key": "value"}]`,
		},
		{
			name:      `root`,
			selector:  `$`,
			document:  `{"key": "value", "another key": {"complex": ["a", 1]}}`,
			consensus: `[{"another key": {"complex": ["a", 1]}, "key": "value"}]`,
		},
		{
			name:      `root_on_scalar`,
			selector:  `$`,
			document:  `42`,
			consensus: `[42]`,
		},
		{
			name:      `root_on_scalar_false`,
			selector:  `$`,
			document:  `false`,
			consensus: `[false]`,
		},
		{
			name:      `root_on_scalar_true`,
			selector:  `$`,
			document:  `true`,
			consensus: `[true]`,
		},
		{
			name:      `union`,
			selector:  `$[0,1]`,
			document:  `["first", "second", "third"]`,
			consensus: `["first", "second"]`,
		},
		{
			name:      `union_with_keys`,
			selector:  `$['key','another']`,
			document:  `{"key": "value", "another": "entry"}`,
			consensus: `["value", "entry"]`,
		},
		{
			name:      `union_with_keys_after_bracket_notation`,
			selector:  `$[0]['c','d']`,
			document:  `[{"c": "cc1", "d": "dd1", "e": "ee1"}, {"c": "cc2", "d": "dd2", "e": "ee2"}]`,
			consensus: `["cc1", "dd1"]`,
		},
		{
			name:      `union_with_keys_on_object_without_key`,
			selector:  `$['missing','key']`,
			document:  `{"key": "value", "another": "entry"}`,
			consensus: `["value"]`,
		},
		{
			name:      `union_with_numbers_in_decreasing_order`,
			selector:  `$[4,1]`,
			document:  `[1, 2, 3, 4, 5]`,
			consensus: `[5, 2]`,
		},
		{
			name:      `union_with_spaces`,
			selector:  `$[ 0 , 1 ]`,
			document:  `["first", "second", "third"]`,
			consensus: `["first", "second"]`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nodes, err := JSONPath([]byte(test.document), test.selector)
			if err != nil {
				t.Errorf("JSONPath() error = %v. got = %v", err, nodes)
				return
			}

			results := make([]interface{}, 0)
			for _, node := range nodes {
				value, err := node.Unpack()
				if err != nil {
					t.Errorf("Unpack(): unexpected error: %v", err)
					return
				}
				results = append(results, value)
			}

			expected, err := Must(Unmarshal([]byte(test.consensus))).Unpack()
			if err != nil {
				t.Errorf("Unpack(): unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(expected, results) {
				t.Errorf("JSONPath(): wrong result:\nSelector: %#+v\nDocument: %s\nExpected: %#+v\nActual:   %#+v", test.selector, test.document, expected, results)
			}
		})
	}
}

func TestJSONPath_special_requests(t *testing.T) {
	tests := []struct {
		selector  string
		document  string
		consensus string
	}{
		{
			selector:  `$.[?(@.name=='special\'')]`,
			document:  `[{"name":"special'"}, {"name":"special"}]`,
			consensus: `[{"name":"special'"}]`,
		},
		{
			selector:  `$.[?(@.name=='special\n')]`,
			document:  `[{"name":"special\n"}, {"name":"special"}]`,
			consensus: `[{"name":"special\n"}]`,
		},
		{
			selector:  `$.[?(@.name==')special(')]`,
			document:  `[{"name":")special("}, {"name":"special"}]`,
			consensus: `[{"name":")special("}]`,
		},
		{
			selector:  `$.[?(@.name==']special[')]`,
			document:  `[{"name":"]special["}, {"name":"special"}]`,
			consensus: `[{"name":"]special["}]`,
		},
		{
			selector:  `$.[?(@.name=='special?')]`,
			document:  `[{"name":"special?"}, {"name":"special"}]`,
			consensus: `[{"name":"special?"}]`,
		},
		{
			selector:  `$.[?(@.name=='special\u3210')]`,
			document:  `[{"name":"special\u3210"}, {"name":"special"}]`,
			consensus: `[{"name":"special\u3210"}]`,
		},
		{
			selector:  `$.[?(@.['special\u3210']=='name')]`,
			document:  `[{"special\u3210":"name"}, {"special":"another"}]`,
			consensus: `[{"special\u3210":"name"}]`,
		},
		{
			selector:  `$..name.title`,
			document:  jpStubs["random_user"],
			consensus: `["Miss"]`,
		},
	}
	for _, test := range tests {
		t.Run(test.selector, func(t *testing.T) {
			nodes, err := JSONPath([]byte(test.document), test.selector)
			if err != nil {
				t.Errorf("JSONPath() error = %v. got = %v", err, nodes)
				return
			}

			results := make([]interface{}, 0)
			for _, node := range nodes {
				value, err := node.Unpack()
				if err != nil {
					t.Errorf("Unpack(): unexpected error: %v", err)
					return
				}
				results = append(results, value)
			}

			expected, err := Must(Unmarshal([]byte(test.consensus))).Unpack()
			if err != nil {
				t.Errorf("Unpack(): unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(expected, results) {
				t.Errorf("JSONPath(): wrong result:\nSelector: %#+v\nDocument: %s\nExpected: %#+v\nActual:   %#+v", test.selector, test.document, expected, results)
			}
		})
	}
}

func TestApplyJSONPath(t *testing.T) {
	node1 := NumericNode("", 1.)
	node2 := NumericNode("", 2.)
	cpy := func(n Node) *Node {
		return &n
	}
	array := ArrayNode("", []*Node{cpy(*node1), cpy(*node2)})

	type args struct {
		node     *Node
		commands []string
	}
	tests := []struct {
		name       string
		args       args
		wantResult []*Node
		wantErr    bool
	}{
		{
			name: "nil",
			args: args{
				node:     nil,
				commands: nil,
			},
			wantResult: make([]*Node, 0),
			wantErr:    false,
		},
		{
			name: "root",
			args: args{
				node:     node1,
				commands: []string{"$"},
			},
			wantResult: []*Node{node1},
			wantErr:    false,
		},
		{
			name: "second",
			args: args{
				node:     array,
				commands: []string{"$", "1"},
			},
			wantResult: []*Node{array.children["1"]},
			wantErr:    false,
		},
		{
			name: "both",
			args: args{
				node:     array,
				commands: []string{"$", "1,0"},
			},
			wantResult: []*Node{array.children["1"], array.children["0"]},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := ApplyJSONPath(tt.args.node, tt.args.commands)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyJSONPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("ApplyJSONPath() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func ExampleApplyJSONPath() {
	json := `[
		[0, 1, 2, 3, 4, 5, 6, 7, 8, 9],
		[0, 1, 2, 3, 4, 5, 6, 7, 8, 9],
		[0, 1, 2, 3, 4, 5, 6, 7, 8, 9],
		[0, 1, 2, 3, 4, 5, 6, 7, 8, 9],
		[0, 1, 2, 3, 4, 5, 6, 7, 8, 9],
		[0, 1, 2, 3, 4, 5, 6, 7, 8, 9],
		[0, 1, 2, 3, 4, 5, 6, 7, 8, 9],
		[0, 1, 2, 3, 4, 5, 6, 7, 8, 9],
		[0, 1, 2, 3, 4, 5, 6, 7, 8, 9],
		[0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
	]`
	node := Must(Unmarshal([]byte(json)))
	for i := 0; i < 10; i++ {
		key1 := strconv.Itoa(i)
		key2 := strconv.Itoa(4 - i)
		nodes, _ := ApplyJSONPath(node, []string{"$", key1, key2})
		fmt.Printf("%s", nodes)
	}

	// Output:
	// [4][3][2][1][0][9][8][7][6][5]
}
