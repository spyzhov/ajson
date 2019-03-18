package ajson

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestNode_ValueSimple(t *testing.T) {
	tests := []struct {
		name     string
		bytes    []byte
		_type    NodeType
		expected interface{}
	}{
		{name: "null", bytes: []byte("null"), _type: Null, expected: nil},
		{name: "1", bytes: []byte("1"), _type: Numeric, expected: float64(1)},
		{name: ".1", bytes: []byte(".1"), _type: Numeric, expected: float64(.1)},
		{name: "-.1e1", bytes: []byte("-.1e1"), _type: Numeric, expected: float64(-1)},
		{name: "string", bytes: []byte("\"foo\""), _type: String, expected: "foo"},
		{name: "space", bytes: []byte("\"foo bar\""), _type: String, expected: "foo bar"},
		{name: "true", bytes: []byte("true"), _type: Bool, expected: true},
		{name: "false", bytes: []byte("false"), _type: Bool, expected: false},
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
				t.Errorf("Error on get value: %s", err.Error())
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
	}
	array, err := root.GetArray()
	if err != nil {
		t.Errorf("Error on root.GetArray(): %s", err.Error())
	}
	if len(array) != 3 {
		t.Errorf("root.GetArray() is corrupted")
	}
}

func TestNode_MustArray(t *testing.T) {
	root, err := Unmarshal([]byte(`[1, 2, 3]`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
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
	}
	value, err := root.GetBool()
	if err != nil {
		t.Errorf("Error on root.GetBool(): %s", err.Error())
	}
	if !value {
		t.Errorf("root.GetBool() is corrupted")
	}
}

func TestNode_MustBool(t *testing.T) {
	root, err := Unmarshal([]byte(`true`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
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
	}
	value, err := root.GetIndex(1)
	if err != nil {
		t.Errorf("Error on root.GetIndex(): %s", err.Error())
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
	}
	value, err := root.GetKey("foo")
	if err != nil {
		t.Errorf("Error on root.GetKey(): %s", err.Error())
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
	}
	value, err := root.GetNull()
	if err != nil {
		t.Errorf("Error on root.GetNull(): %s", err.Error())
	}
	if value != nil {
		t.Errorf("root.GetNull() is corrupted")
	}
}

func TestNode_MustNull(t *testing.T) {
	root, err := Unmarshal([]byte(`null`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
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
	}
	value, err := root.GetNumeric()
	if err != nil {
		t.Errorf("Error on root.GetNumeric(): %s", err.Error())
	}
	if value != float64(123) {
		t.Errorf("root.GetNumeric() is corrupted")
	}
}

func TestNode_MustNumeric(t *testing.T) {
	root, err := Unmarshal([]byte(`123`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
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
}

func TestNode_MustObject(t *testing.T) {
	root, err := Unmarshal([]byte(`{"foo":true,"bar":null}`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
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
	}
	value, err := root.GetString()
	if err != nil {
		t.Errorf("Error on root.GetString(): %s", err.Error())
	}
	if value != "123" {
		t.Errorf("root.GetString() is corrupted")
	}
}

func TestNode_MustString(t *testing.T) {
	root, err := Unmarshal([]byte(`"123"`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
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
	root, err := Unmarshal([]byte(`+1.23e-1.01`))
	if err != nil {
		t.Errorf("Error on Unmarshal(): %s", err.Error())
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
	}
	value := root.String()
	if value != `{"foo":true,"bar":null}` {
		t.Errorf("Wrong root.String()")
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
