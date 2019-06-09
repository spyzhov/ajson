package ajson

import (
	"bytes"
	"encoding/json"
	"math"
	"testing"
)

func TestNode_Value_Simple(t *testing.T) {
	tests := []struct {
		name     string
		bytes    []byte
		_type    NodeType
		expected interface{}
		error    bool
	}{
		{name: "null", bytes: []byte("null"), _type: Null, expected: nil},
		{name: "1", bytes: []byte("1"), _type: Numeric, expected: float64(1)},
		{name: ".1", bytes: []byte(".1"), _type: Numeric, expected: float64(.1)},
		{name: "-.1e1", bytes: []byte("-.1e1"), _type: Numeric, expected: float64(-1)},
		{name: "string", bytes: []byte("\"foo\""), _type: String, expected: "foo"},
		{name: "space", bytes: []byte("\"foo bar\""), _type: String, expected: "foo bar"},
		{name: "true", bytes: []byte("true"), _type: Bool, expected: true},
		{name: "false", bytes: []byte("false"), _type: Bool, expected: false},
		{name: "e1", bytes: []byte("e1"), _type: Numeric, error: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			current := &Node{
				_type:   test._type,
				borders: [2]int{0, len(test.bytes)},
				data:    &test.bytes,
			}
			value, err := current.Value()
			if err != nil {
				if !test.error {
					t.Errorf("Error on get value: %s", err.Error())
				}
			} else if value != test.expected {
				t.Errorf("Error on get value: '%v' != '%v'", value, test.expected)
			} else if value2, err := current.Value(); err != nil {
				t.Errorf("Error on get value 2: %s", err.Error())
			} else if value != value2 {
				t.Errorf("Error on get value 2: '%v' != '%v'", value, value2)
			}
		})
	}
}

func TestNode_Unpack(t *testing.T) {
	tests := []struct {
		value string
	}{
		{value: `1`},
		{value: `true`},
		{value: `null`},
		{value: `{}`},
		{value: `[]`},
		{value: `[1,2,3]`},
		{value: `[1,{},null]`},
		{value: `{"foo":["bar",null]}`},
	}
	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			root, err := Unmarshal([]byte(test.value))
			if err != nil {
				t.Errorf("Error on Unmarshal(): %s", err.Error())
				return
			}
			unpacked, err := root.Unpack()
			if err != nil {
				t.Errorf("Error on root.Unpack(): %s", err.Error())
			}
			marshalled, err := json.Marshal(unpacked)
			if err != nil {
				t.Errorf("Error on json.Marshal(): %s", err.Error())
			}
			if string(marshalled) != test.value {
				t.Errorf("Wrong structure: '%s' != '%s'", string(marshalled), test.value)
			}
		})
	}
}

func TestNode_Value(t *testing.T) {
	root, err := Unmarshal([]byte(`{ "category": null,
        "author": "Evelyn Waugh",
        "title": "Sword of Honour",
        "price": 12.99,
        "ordered": true,
        "tags": [],
        "sub": {}
      }`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	iValue, err := root.Value()
	if err != nil {
		t.Errorf("Error on root.Value(): %s", err.Error())
	}
	value, ok := iValue.(map[string]*Node)
	if !ok {
		t.Errorf("Value is not an Object map")
	}
	keys := []string{"category", "author", "title", "price", "ordered", "tags", "sub"}
	for _, key := range keys {
		if _, ok := value[key]; !ok {
			t.Errorf("Object map has no field: " + key)
		}
	}
}

func TestNode_Empty(t *testing.T) {
	root, err := Unmarshal([]byte(`{
        "tag1": [1, 2, 3],
        "tag2": [],
        "sub1": {},
        "sub2": {"foo":null}
      }`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	iValue, err := root.Value()
	if err != nil {
		t.Errorf("Error on root.Value(): %s", err.Error())
	}
	value, ok := iValue.(map[string]*Node)
	if !ok {
		t.Errorf("Value is not an Object map")
	}
	if value["tag1"].Empty() {
		t.Errorf("Node `tag1` is not empty")
	}
	if !value["tag2"].Empty() {
		t.Errorf("Node `tag2` is empty")
	}
	if value["sub2"].Empty() {
		t.Errorf("Node `sub2` is not empty")
	}
	if !value["sub1"].Empty() {
		t.Errorf("Node `sub1` is empty")
	}
}

func TestNode_GetArray(t *testing.T) {
	root, err := Unmarshal([]byte(`[1, 2, 3]`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	array, err := root.GetArray()
	if err != nil {
		t.Errorf("Error on root.GetArray(): %s", err.Error())
	}
	if len(array) != 3 {
		t.Errorf("root.GetArray() is corrupted")
	}

	root = NullNode("")
	_, err = root.GetArray()
	if err == nil {
		t.Errorf("Error on root.GetArray(): NullNode")
	}
}

func TestNode_MustArray(t *testing.T) {
	root, err := Unmarshal([]byte(`[1, 2, 3]`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	array := root.MustArray()
	if len(array) != 3 {
		t.Errorf("root.GetArray() is corrupted")
	}
}

func TestNode_GetBool(t *testing.T) {
	root, err := Unmarshal([]byte(`true`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	value, err := root.GetBool()
	if err != nil {
		t.Errorf("Error on root.GetBool(): %s", err.Error())
	}
	if !value {
		t.Errorf("root.GetBool() is corrupted")
	}

	root = NullNode("")
	_, err = root.GetBool()
	if err == nil {
		t.Errorf("Error on root.GetBool(): NullNode")
	}
}

func TestNode_MustBool(t *testing.T) {
	root, err := Unmarshal([]byte(`true`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	value := root.MustBool()
	if !value {
		t.Errorf("root.MustBool() is corrupted")
	}
}

func TestNode_GetIndex(t *testing.T) {
	root, err := Unmarshal([]byte(`[1, 2, 3]`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	value, err := root.GetIndex(1)
	if err != nil {
		t.Errorf("Error on root.GetIndex(): %s", err.Error())
		return
	}
	if value.MustNumeric() != 2 {
		t.Errorf("root.GetIndex() is corrupted")
	}
	value, err = root.GetIndex(10)
	if err == nil {
		t.Errorf("Error on root.GetIndex() - out of range")
	}
	if value != nil {
		t.Errorf("Error on root.GetIndex() - wrong value")
	}
}

func TestNode_MustIndex(t *testing.T) {
	root, err := Unmarshal([]byte(`[1, 2, 3]`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	value := root.MustIndex(1)
	if value.MustNumeric() != 2 {
		t.Errorf("root.GetIndex() is corrupted")
	}
}

func TestNode_GetKey(t *testing.T) {
	root, err := Unmarshal([]byte(`{"foo":2,"bar":null}`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	value, err := root.GetKey("foo")
	if err != nil {
		t.Errorf("Error on root.GetKey(): %s", err.Error())
		return
	}
	if value.MustNumeric() != 2 {
		t.Errorf("root.GetKey() is corrupted")
	}
	value, err = root.GetKey("baz")
	if err == nil {
		t.Errorf("Error on root.GetKey() - wrong element")
	}
	if value != nil {
		t.Errorf("Error on root.GetKey() - wrong value")
	}
}

func TestNode_MustKey(t *testing.T) {
	root, err := Unmarshal([]byte(`{"foo":2,"bar":null}`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	value := root.MustKey("foo")
	if value.MustNumeric() != 2 {
		t.Errorf("root.GetKey() is corrupted")
	}
}

func TestNode_GetNull(t *testing.T) {
	root, err := Unmarshal([]byte(`null`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	value, err := root.GetNull()
	if err != nil {
		t.Errorf("Error on root.GetNull(): %s", err.Error())
	}
	if value != nil {
		t.Errorf("root.GetNull() is corrupted")
	}

	root = NumericNode("", 1)
	_, err = root.GetNull()
	if err == nil {
		t.Errorf("Error expected on root.GetNull() using NumericNode")
	}
}

func TestNode_MustNull(t *testing.T) {
	root, err := Unmarshal([]byte(`null`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	value := root.MustNull()
	if value != nil {
		t.Errorf("root.MustNull() is corrupted")
	}
}

func TestNode_GetNumeric(t *testing.T) {
	root, err := Unmarshal([]byte(`123`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	value, err := root.GetNumeric()
	if err != nil {
		t.Errorf("Error on root.GetNumeric(): %s", err.Error())
	}
	if value != float64(123) {
		t.Errorf("root.GetNumeric() is corrupted")
	}

	root = StringNode("", "")
	_, err = root.GetNumeric()
	if err == nil {
		t.Errorf("Error on root.GetNumeric() using StringNode")
	}

	root = valueNode(nil, "", Numeric, "foo")
	_, err = root.GetNumeric()
	if err == nil {
		t.Errorf("Error on root.GetNumeric() wrong data")
	}
}

func TestNode_MustNumeric(t *testing.T) {
	root, err := Unmarshal([]byte(`123`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	value := root.MustNumeric()
	if value != float64(123) {
		t.Errorf("root.GetNumeric() is corrupted")
	}
}

func TestNode_GetObject(t *testing.T) {
	root, err := Unmarshal([]byte(`{"foo":true,"bar":null}`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	value, err := root.GetObject()
	if err != nil {
		t.Errorf("Error on root.GetObject(): %s", err.Error())
	}
	if _, ok := value["foo"]; !ok {
		t.Errorf("root.GetObject() is corrupted: foo")
	}
	if _, ok := value["bar"]; !ok {
		t.Errorf("root.GetObject() is corrupted: bar")
	}

	root = NullNode("")
	_, err = root.GetObject()
	if err == nil {
		t.Errorf("Error on root.GetArray(): NullNode")
	}
}

func TestNode_MustObject(t *testing.T) {
	root, err := Unmarshal([]byte(`{"foo":true,"bar":null}`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	value := root.MustObject()
	if _, ok := value["foo"]; !ok {
		t.Errorf("root.GetObject() is corrupted: foo")
	}
	if _, ok := value["bar"]; !ok {
		t.Errorf("root.GetObject() is corrupted: bar")
	}
}

func TestNode_GetString(t *testing.T) {
	root, err := Unmarshal([]byte(`"123"`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	value, err := root.GetString()
	if err != nil {
		t.Errorf("Error on root.GetString(): %s", err.Error())
	}
	if value != "123" {
		t.Errorf("root.GetString() is corrupted")
	}

	root = NumericNode("", 1)
	_, err = root.GetString()
	if err == nil {
		t.Errorf("Error on root.GetString(): NumericNode")
	}
}

func TestNode_MustString(t *testing.T) {
	root, err := Unmarshal([]byte(`"123"`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	value := root.MustString()
	if value != "123" {
		t.Errorf("root.GetString() is corrupted")
	}
}

func TestNode_Index(t *testing.T) {
	root, err := Unmarshal([]byte(`[1, 2, 3]`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	array := root.MustArray()
	for i, node := range array {
		if i != node.Index() {
			t.Errorf("Wrong node.Index(): %d != %d", i, node.Index())
		}
	}
}

func TestNode_Key(t *testing.T) {
	root, err := Unmarshal([]byte(`{"foo":"bar", "baz":null}`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	object := root.MustObject()
	for key, node := range object {
		if key != node.Key() {
			t.Errorf("Wrong node.Index(): '%s' != '%s'", key, node.Key())
		}
	}
}

func TestNode_IsArray(t *testing.T) {
	root, err := Unmarshal([]byte(`[1, 2, 3]`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	if !root.IsArray() {
		t.Errorf("Wrong root.IsArray()")
	}
	if root.IsObject() {
		t.Errorf("Wrong root.IsObject()")
	}
	if root.IsString() {
		t.Errorf("Wrong root.IsString()")
	}
	if root.IsNumeric() {
		t.Errorf("Wrong root.IsNumeric()")
	}
	if root.IsBool() {
		t.Errorf("Wrong root.IsBool()")
	}
	if root.IsNull() {
		t.Errorf("Wrong root.IsNull()")
	}
}

func TestNode_IsObject(t *testing.T) {
	root, err := Unmarshal([]byte(`{"foo":null}`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	if root.IsArray() {
		t.Errorf("Wrong root.IsArray()")
	}
	if !root.IsObject() {
		t.Errorf("Wrong root.IsObject()")
	}
	if root.IsString() {
		t.Errorf("Wrong root.IsString()")
	}
	if root.IsNumeric() {
		t.Errorf("Wrong root.IsNumeric()")
	}
	if root.IsBool() {
		t.Errorf("Wrong root.IsBool()")
	}
	if root.IsNull() {
		t.Errorf("Wrong root.IsNull()")
	}
}

func TestNode_IsString(t *testing.T) {
	root, err := Unmarshal([]byte(`"123"`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	if root.IsArray() {
		t.Errorf("Wrong root.IsArray()")
	}
	if root.IsObject() {
		t.Errorf("Wrong root.IsObject()")
	}
	if !root.IsString() {
		t.Errorf("Wrong root.IsString()")
	}
	if root.IsNumeric() {
		t.Errorf("Wrong root.IsNumeric()")
	}
	if root.IsBool() {
		t.Errorf("Wrong root.IsBool()")
	}
	if root.IsNull() {
		t.Errorf("Wrong root.IsNull()")
	}
}

func TestNode_IsNumeric(t *testing.T) {
	root, err := Unmarshal([]byte(`+1.23e-2`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	if root.IsArray() {
		t.Errorf("Wrong root.IsArray()")
	}
	if root.IsObject() {
		t.Errorf("Wrong root.IsObject()")
	}
	if root.IsString() {
		t.Errorf("Wrong root.IsString()")
	}
	if !root.IsNumeric() {
		t.Errorf("Wrong root.IsNumeric()")
	}
	if root.IsBool() {
		t.Errorf("Wrong root.IsBool()")
	}
	if root.IsNull() {
		t.Errorf("Wrong root.IsNull()")
	}
}

func TestNode_IsBool(t *testing.T) {
	root, err := Unmarshal([]byte(`true`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	if root.IsArray() {
		t.Errorf("Wrong root.IsArray()")
	}
	if root.IsObject() {
		t.Errorf("Wrong root.IsObject()")
	}
	if root.IsString() {
		t.Errorf("Wrong root.IsString()")
	}
	if root.IsNumeric() {
		t.Errorf("Wrong root.IsNumeric()")
	}
	if !root.IsBool() {
		t.Errorf("Wrong root.IsBool()")
	}
	if root.IsNull() {
		t.Errorf("Wrong root.IsNull()")
	}
}

func TestNode_IsNull(t *testing.T) {
	root, err := Unmarshal([]byte(`null`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	if root.IsArray() {
		t.Errorf("Wrong root.IsArray()")
	}
	if root.IsObject() {
		t.Errorf("Wrong root.IsObject()")
	}
	if root.IsString() {
		t.Errorf("Wrong root.IsString()")
	}
	if root.IsNumeric() {
		t.Errorf("Wrong root.IsNumeric()")
	}
	if root.IsBool() {
		t.Errorf("Wrong root.IsBool()")
	}
	if !root.IsNull() {
		t.Errorf("Wrong root.IsNull()")
	}
}

func TestNode_Keys(t *testing.T) {
	root, err := Unmarshal([]byte(`{"foo":true,"bar":null}`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	value := root.Keys()
	if len(value) != 2 {
		t.Errorf("Wrong root.Keys()")
	}
	if value[0] != "foo" && value[0] != "bar" {
		t.Errorf("Wrong value in 0")
	}
	if value[1] != "foo" && value[1] != "bar" {
		t.Errorf("Wrong value in 1")
	}
}

func TestNode_Size(t *testing.T) {
	root, err := Unmarshal([]byte(`[1,2,3,4]`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	value := root.Size()
	if value != 4 {
		t.Errorf("Wrong root.Size()")
	}
}

func TestNode_Parent(t *testing.T) {
	root, err := Unmarshal([]byte(`{"foo":true,"bar":null}`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	value := root.Parent()
	if value != nil {
		t.Errorf("Wrong root.Parent()")
	}
	value = root.MustKey("foo")
	if value.Parent().String() != root.String() {
		t.Errorf("Wrong value.Parent()")
	}
}

func TestNode_Source(t *testing.T) {
	root, err := Unmarshal([]byte(`{"foo":true,"bar":null}`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	value := root.Source()
	if !bytes.Equal(value, []byte(`{"foo":true,"bar":null}`)) {
		t.Errorf("Wrong root.Source()")
	}
}

func TestNode_String(t *testing.T) {
	root, err := Unmarshal([]byte(`{"foo":true,"bar":null}`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	value := root.String()
	if value != `{"foo":true,"bar":null}` {
		t.Errorf("Wrong (Unmarshal) root.String()")
	}

	root = StringNode("", "foo")
	value = root.String()
	if value != "foo" {
		t.Errorf("Wrong (StringNode) root.String()")
	}

	root = NullNode("")
	value = root.String()
	if value != "null" {
		t.Errorf("Wrong (NullNode) root.String()")
	}
}

func TestNode_Type(t *testing.T) {
	tests := []struct {
		_type NodeType
		value string
	}{
		{value: "null", _type: Null},
		{value: "123", _type: Numeric},
		{value: "1.23e+3", _type: Numeric},
		{value: `"1.23e+3"`, _type: String},
		{value: `["1.23e+3"]`, _type: Array},
		{value: `[]`, _type: Array},
		{value: `{}`, _type: Object},
		{value: `{"foo":1.23e+3}`, _type: Object},
		{value: `true`, _type: Bool},
		{value: `false`, _type: Bool},
	}
	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			root, err := Unmarshal([]byte(test.value))
			if err != nil {
				t.Errorf("Error on Unmarshal('%s'): %s", test.value, err.Error())
			} else if root.Type() != test._type {
				t.Errorf("Wrong type on Unmarshal('%s')", test.value)
			}
		})
	}
}

func TestNode_HasKey(t *testing.T) {
	root, err := Unmarshal([]byte(`{"foo":true,"bar":null}`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	if !root.HasKey("foo") {
		t.Errorf("Wrong root.HasKey('foo')")
	}
	if !root.HasKey("bar") {
		t.Errorf("Wrong root.HasKey('bar')")
	}
	if root.HasKey("baz") {
		t.Errorf("Wrong root.HasKey('bar')")
	}
}

func TestNode_Path(t *testing.T) {
	data := []byte(`{
        "Image": {
            "Width":  800,
            "Height": 600,
            "Title":  "View from 15th Floor",
            "Thumbnail": {
                "Url":    "http://www.example.com/image/481989943",
                "Height": 125,
                "Width":  100
            },
            "Animated" : false,
            "IDs": [116, 943, 234, 38793]
          }
      }`)
	root, err := Unmarshal(data)
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
		return
	}
	if root.Path() != "$" {
		t.Errorf("Wrong root.Path()")
	}
	element := root.MustKey("Image").MustKey("Thumbnail").MustKey("Url")
	if element.Path() != "$['Image']['Thumbnail']['Url']" {
		t.Errorf("Wrong element.Path()")
	}
}

func TestNode_Eq(t *testing.T) {
	tests := []struct {
		name        string
		left, right *Node
		expected    bool
		error       bool
	}{
		{
			name:     "simple",
			left:     valueNode(nil, "bool", Bool, true),
			right:    valueNode(nil, "bool", Bool, true),
			expected: true,
		},
		{
			name:     "null",
			left:     valueNode(nil, "null", Null, nil),
			right:    valueNode(nil, "null", Null, nil),
			expected: true,
		},
		{
			name:     "float",
			left:     valueNode(nil, "123.5", Numeric, float64(123.5)),
			right:    valueNode(nil, "123.5", Numeric, float64(123.5)),
			expected: true,
		},
		{
			name:     "blank array",
			left:     valueNode(nil, "[]", Array, []*Node{}),
			right:    valueNode(nil, "[]", Array, []*Node{}),
			expected: true,
		},
		{
			name:     "blank map",
			left:     valueNode(nil, "{}", Object, map[string]*Node{}),
			right:    valueNode(nil, "{}", Object, map[string]*Node{}),
			expected: true,
		},
		{
			name:     "blank map and array",
			left:     valueNode(nil, "{}", Object, map[string]*Node{}),
			right:    valueNode(nil, "[]", Array, []*Node{}),
			expected: false,
		},
		{
			name:     "floats 1",
			left:     NumericNode("", 1.1),
			right:    NumericNode("", 1.2),
			expected: false,
		},
		{
			name:     "floats 2",
			left:     NumericNode("", -1),
			right:    NumericNode("", 1),
			expected: false,
		},
		{
			name:     "floats 3",
			left:     NumericNode("", 1.0001),
			right:    NumericNode("", 1.00011),
			expected: false,
		},
		{
			name:  "error 1",
			left:  valueNode(nil, "", Numeric, "foo"),
			right: NumericNode("", 1.00011),
			error: true,
		},
		{
			name:  "error 1",
			left:  valueNode(nil, "", Numeric, "foo"),
			right: valueNode(nil, "", Numeric, float64(1)),
			error: true,
		},
		{
			name:  "error 2",
			left:  valueNode(nil, "", String, "foo"),
			right: valueNode(nil, "", String, float64(1)),
			error: true,
		},
		{
			name:  "error 3",
			left:  valueNode(nil, "", Bool, "foo"),
			right: valueNode(nil, "", Bool, float64(1)),
			error: true,
		},
		{
			name:  "error 4",
			left:  valueNode(nil, "", Array, "foo"),
			right: valueNode(nil, "", Array, float64(1)),
			error: true,
		},
		{
			name:  "error 5",
			left:  valueNode(nil, "", Object, "foo"),
			right: valueNode(nil, "", Object, float64(1)),
			error: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := test.left.Eq(test.right)
			if test.error {
				if err == nil {
					t.Errorf("Error expected: nil given")
				}
			} else if err != nil {
				t.Errorf("Error on node.Eq(): %s", err.Error())
			} else if actual != test.expected {
				t.Errorf("Failed node.Eq()")
			}
		})
	}
}

func TestNode_Neq(t *testing.T) {
	tests := []struct {
		name        string
		left, right *Node
		expected    bool
	}{
		{
			name:     "simple",
			left:     valueNode(nil, "bool", Bool, true),
			right:    valueNode(nil, "bool", Bool, true),
			expected: false,
		},
		{
			name:     "null",
			left:     valueNode(nil, "null", Null, nil),
			right:    valueNode(nil, "null", Null, nil),
			expected: false,
		},
		{
			name:     "float",
			left:     valueNode(nil, "123.5", Numeric, float64(123.5)),
			right:    valueNode(nil, "123.5", Numeric, float64(123.5)),
			expected: false,
		},
		{
			name:     "blank array",
			left:     valueNode(nil, "[]", Array, []*Node{}),
			right:    valueNode(nil, "[]", Array, []*Node{}),
			expected: false,
		},
		{
			name:     "blank map",
			left:     valueNode(nil, "{}", Object, map[string]*Node{}),
			right:    valueNode(nil, "{}", Object, map[string]*Node{}),
			expected: false,
		},
		{
			name:     "blank map and array",
			left:     valueNode(nil, "{}", Object, map[string]*Node{}),
			right:    valueNode(nil, "[]", Array, []*Node{}),
			expected: true,
		},
		{
			name:     "floats 1",
			left:     NumericNode("", 1.1),
			right:    NumericNode("", 1.2),
			expected: true,
		},
		{
			name:     "floats 2",
			left:     NumericNode("", -1),
			right:    NumericNode("", 1),
			expected: true,
		},
		{
			name:     "floats 3",
			left:     NumericNode("", 1.0001),
			right:    NumericNode("", 1.00011),
			expected: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := test.left.Neq(test.right)
			if err != nil {
				t.Errorf("Error on node.Neq(): %s", err.Error())
			} else if actual != test.expected {
				t.Errorf("Failed node.Neq()")
			}
		})
	}
}

func TestNode_Ge(t *testing.T) {
	tests := []struct {
		name        string
		left, right *Node
		expected    bool
		error       bool
	}{
		{
			name:  "null",
			left:  NullNode(""),
			right: NullNode(""),
			error: true,
		},
		{
			name:  "array",
			left:  ArrayNode("", nil),
			right: ArrayNode("", nil),
			error: true,
		},
		{
			name:  "object",
			left:  ObjectNode("", nil),
			right: ObjectNode("", nil),
			error: true,
		},
		{
			name:     "float 1",
			left:     NumericNode("", 3.1),
			right:    NumericNode("", 3),
			expected: true,
		},
		{
			name:     "float 2",
			left:     NumericNode("", 0),
			right:    NumericNode("", -3),
			expected: true,
		},
		{
			name:     "float 3",
			left:     NumericNode("", 0),
			right:    NumericNode("", 0),
			expected: false,
		},
		{
			name:     "float 4",
			left:     NumericNode("", math.MaxFloat64),
			right:    NumericNode("", math.SmallestNonzeroFloat64),
			expected: true,
		},
		{
			name:     "float 5",
			left:     NumericNode("", math.SmallestNonzeroFloat64),
			right:    NumericNode("", math.MaxFloat64),
			expected: false,
		},
		{
			name:     "string 1",
			left:     StringNode("", "z"),
			right:    StringNode("", "a"),
			expected: true,
		},
		{
			name:     "string 2",
			left:     StringNode("", "a"),
			right:    StringNode("", "a"),
			expected: false,
		},
		{
			name:     "wrong type 1",
			left:     StringNode("", "z"),
			right:    NumericNode("", math.MaxFloat64),
			expected: false,
		},
		{
			name:     "wrong type 2",
			left:     NumericNode("", math.MaxFloat64),
			right:    StringNode("", "z"),
			expected: false,
		},
		{
			name:  "error 1",
			left:  valueNode(nil, "e1", Numeric, string("e1")),
			right: NumericNode("", 1),
			error: true,
		},
		{
			name:  "error 2",
			left:  valueNode(nil, "e1", String, float64(1)),
			right: StringNode("", "foo"),
			error: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := test.left.Ge(test.right)
			if err != nil {
				if !test.error {
					t.Errorf("Error on node.Ge(): %s", err.Error())
				}
			} else if actual != test.expected {
				t.Errorf("Failed node.Ge()")
			}

		})
	}
}

func TestNode_Geq(t *testing.T) {
	tests := []struct {
		name        string
		left, right *Node
		expected    bool
		error       bool
	}{
		{
			name:  "null",
			left:  NullNode(""),
			right: NullNode(""),
			error: true,
		},
		{
			name:  "array",
			left:  ArrayNode("", nil),
			right: ArrayNode("", nil),
			error: true,
		},
		{
			name:  "object",
			left:  ObjectNode("", nil),
			right: ObjectNode("", nil),
			error: true,
		},
		{
			name:     "float 1",
			left:     NumericNode("", 3.1),
			right:    NumericNode("", 3),
			expected: true,
		},
		{
			name:     "float 2",
			left:     NumericNode("", 0),
			right:    NumericNode("", -3),
			expected: true,
		},
		{
			name:     "float 3",
			left:     NumericNode("", 0),
			right:    NumericNode("", 0),
			expected: true,
		},
		{
			name:     "float 4",
			left:     NumericNode("", math.MaxFloat64),
			right:    NumericNode("", math.SmallestNonzeroFloat64),
			expected: true,
		},
		{
			name:     "float 5",
			left:     NumericNode("", math.SmallestNonzeroFloat64),
			right:    NumericNode("", math.MaxFloat64),
			expected: false,
		},
		{
			name:     "string 1",
			left:     StringNode("", "z"),
			right:    StringNode("", "a"),
			expected: true,
		},
		{
			name:     "string 2",
			left:     StringNode("", "a"),
			right:    StringNode("", "a"),
			expected: true,
		},
		{
			name:     "wrong type 1",
			left:     StringNode("", "z"),
			right:    NumericNode("", math.MaxFloat64),
			expected: false,
		},
		{
			name:     "wrong type 2",
			left:     NumericNode("", math.MaxFloat64),
			right:    StringNode("", "z"),
			expected: false,
		},
		{
			name:  "error 1",
			left:  valueNode(nil, "e1", Numeric, string("e1")),
			right: NumericNode("", 1),
			error: true,
		},
		{
			name:  "error 2",
			left:  valueNode(nil, "e1", String, float64(1)),
			right: StringNode("", "foo"),
			error: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := test.left.Geq(test.right)
			if err != nil {
				if !test.error {
					t.Errorf("Error on node.Geq(): %s", err.Error())
				}
			} else if actual != test.expected {
				t.Errorf("Failed node.Geq()")
			}

		})
	}
}

func TestNode_Le(t *testing.T) {
	tests := []struct {
		name        string
		left, right *Node
		expected    bool
		error       bool
	}{
		{
			name:  "null",
			left:  NullNode(""),
			right: NullNode(""),
			error: true,
		},
		{
			name:  "array",
			left:  ArrayNode("", nil),
			right: ArrayNode("", nil),
			error: true,
		},
		{
			name:  "object",
			left:  ObjectNode("", nil),
			right: ObjectNode("", nil),
			error: true,
		},
		{
			name:     "float 1",
			left:     NumericNode("", 3.1),
			right:    NumericNode("", 3),
			expected: false,
		},
		{
			name:     "float 2",
			left:     NumericNode("", 0),
			right:    NumericNode("", -3),
			expected: false,
		},
		{
			name:     "float 3",
			left:     NumericNode("", 0),
			right:    NumericNode("", 0),
			expected: false,
		},
		{
			name:     "float 4",
			left:     NumericNode("", math.MaxFloat64),
			right:    NumericNode("", math.SmallestNonzeroFloat64),
			expected: false,
		},
		{
			name:     "float 5",
			left:     NumericNode("", math.SmallestNonzeroFloat64),
			right:    NumericNode("", math.MaxFloat64),
			expected: true,
		},
		{
			name:     "string 1",
			left:     StringNode("", "z"),
			right:    StringNode("", "a"),
			expected: false,
		},
		{
			name:     "string 2",
			left:     StringNode("", "a"),
			right:    StringNode("", "a"),
			expected: false,
		},
		{
			name:     "wrong type 1",
			left:     StringNode("", "z"),
			right:    NumericNode("", math.MaxFloat64),
			expected: false,
		},
		{
			name:     "wrong type 2",
			left:     NumericNode("", math.MaxFloat64),
			right:    StringNode("", "z"),
			expected: false,
		},
		{
			name:  "error 1",
			left:  valueNode(nil, "e1", Numeric, string("e1")),
			right: NumericNode("", 1),
			error: true,
		},
		{
			name:  "error 2",
			left:  valueNode(nil, "e1", String, float64(1)),
			right: StringNode("", "foo"),
			error: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := test.left.Le(test.right)
			if err != nil {
				if !test.error {
					t.Errorf("Error on node.Le(): %s", err.Error())
				}
			} else if actual != test.expected {
				t.Errorf("Failed node.Le()")
			}

		})
	}
}

func TestNode_Leq(t *testing.T) {
	tests := []struct {
		name        string
		left, right *Node
		expected    bool
		error       bool
	}{
		{
			name:  "null",
			left:  NullNode(""),
			right: NullNode(""),
			error: true,
		},
		{
			name:  "array",
			left:  ArrayNode("", nil),
			right: ArrayNode("", nil),
			error: true,
		},
		{
			name:  "object",
			left:  ObjectNode("", nil),
			right: ObjectNode("", nil),
			error: true,
		},
		{
			name:     "float 1",
			left:     NumericNode("", 3.1),
			right:    NumericNode("", 3),
			expected: false,
		},
		{
			name:     "float 2",
			left:     NumericNode("", 0),
			right:    NumericNode("", -3),
			expected: false,
		},
		{
			name:     "float 3",
			left:     NumericNode("", 0),
			right:    NumericNode("", 0),
			expected: true,
		},
		{
			name:     "float 4",
			left:     NumericNode("", math.MaxFloat64),
			right:    NumericNode("", math.SmallestNonzeroFloat64),
			expected: false,
		},
		{
			name:     "float 5",
			left:     NumericNode("", math.SmallestNonzeroFloat64),
			right:    NumericNode("", math.MaxFloat64),
			expected: true,
		},
		{
			name:     "string 1",
			left:     StringNode("", "z"),
			right:    StringNode("", "a"),
			expected: false,
		},
		{
			name:     "string 2",
			left:     StringNode("", "a"),
			right:    StringNode("", "a"),
			expected: true,
		},
		{
			name:     "wrong type 1",
			left:     StringNode("", "z"),
			right:    NumericNode("", math.MaxFloat64),
			expected: false,
		},
		{
			name:     "wrong type 2",
			left:     NumericNode("", math.MaxFloat64),
			right:    StringNode("", "z"),
			expected: false,
		},
		{
			name:  "error 1",
			left:  valueNode(nil, "e1", Numeric, string("e1")),
			right: NumericNode("", 1),
			error: true,
		},
		{
			name:  "error 2",
			left:  valueNode(nil, "e1", String, float64(1)),
			right: StringNode("", "foo"),
			error: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := test.left.Leq(test.right)
			if err != nil {
				if !test.error {
					t.Errorf("Error on node.Leq(): %s", err.Error())
				}
			} else if actual != test.expected {
				t.Errorf("Failed node.Leq()")
			}

		})
	}
}

func TestNullNode(t *testing.T) {
	node := NullNode("test")
	if node.MustNull() != nil {
		t.Errorf("Failed")
	}
}

func TestNumericNode(t *testing.T) {
	node := NumericNode("test", 1.5)
	if node.MustNumeric() != 1.5 {
		t.Errorf("Failed")
	}
}

func TestStringNode(t *testing.T) {
	node := StringNode("test", "check")
	if node.MustString() != "check" {
		t.Errorf("Failed")
	}
}

func TestBoolNode(t *testing.T) {
	node := BoolNode("test", true)
	if !node.MustBool() {
		t.Errorf("Failed")
	}
}

func TestArrayNode(t *testing.T) {
	array := []*Node{
		NullNode("0"),
		NumericNode("1", 1),
		StringNode("str", "foo"),
	}
	node := ArrayNode("test", array)
	result := node.MustArray()
	if len(result) != len(array) {
		t.Errorf("Failed: length")
	}
	for i, val := range result {
		ok, err := val.Eq(array[i])
		if err != nil {
			t.Errorf("Failed: %s", err.Error())
		} else if !ok {
			t.Errorf("Failed: compare '%s' & '%s'", val, array[i])
		}
	}
}

func TestObjectNode(t *testing.T) {
	objects := map[string]*Node{
		"zero": NullNode("0"),
		"foo":  NumericNode("1", 1),
		"bar":  StringNode("str", "foo"),
	}
	node := ObjectNode("test", objects)
	result := node.MustObject()
	if len(result) != len(objects) {
		t.Errorf("Failed: length")
	}
	for i, val := range result {
		ok, err := val.Eq(objects[i])
		if err != nil {
			t.Errorf("Failed: %s", err.Error())
		} else if !ok {
			t.Errorf("Failed: compare '%s' & '%s'", val, objects[i])
		}
	}
}

func TestNode_Inheritors(t *testing.T) {
	tests := []struct {
		name     string
		node     *Node
		expected []*Node
	}{
		{
			name: "object",
			node: ObjectNode("", map[string]*Node{
				"zero": NullNode("0"),
				"foo":  NumericNode("1", 1),
				"bar":  StringNode("str", "foo"),
			}),
			expected: []*Node{
				StringNode("str", "foo"),
				NumericNode("1", 1),
				NullNode("0"),
			},
		},
		{
			name: "array",
			node: ArrayNode("", []*Node{
				NullNode("0"),
				NumericNode("1", 1),
				StringNode("str", "foo"),
			}),
			expected: []*Node{
				NullNode("0"),
				NumericNode("1", 1),
				StringNode("str", "foo"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.node.Inheritors()
			if len(result) != len(test.expected) {
				t.Errorf("Failed: wrong size")
			} else {
				for i, node := range test.expected {
					if ok, err := node.Eq(result[i]); err != nil {
						t.Errorf("Failed: %s", err.Error())
					} else if !ok {
						t.Errorf("Failed: '%s' != '%s'", node, result[i])
					}
				}
			}
		})
	}
}

func TestNode_JSONPath(t *testing.T) {
	root, err := Unmarshal(jsonpathTestData)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
		return
	}
	result, err := root.MustKey("store").MustKey("book").JSONPath("@.*")
	if err != nil {
		t.Errorf("Error: %s", err.Error())
		return
	}
	if len(result) != 4 {
		t.Errorf("Error: JSONPath")
	}
}
