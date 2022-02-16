package ajson

import (
	"bytes"
	"encoding/json"
	"math"
	"reflect"
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
		{name: "string error", bytes: []byte("\"foo\nbar\""), _type: String, error: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			current := &Node{
				_type:   test._type,
				borders: [2]int{0, len(test.bytes)},
				data:    &test.bytes,
			}
			value, err := current.getValue()
			if err != nil {
				if !test.error {
					t.Errorf("Error on get value: %s", err.Error())
				}
			} else if value != test.expected {
				t.Errorf("Error on get value: '%v' != '%v'", value, test.expected)
			} else if value2, err := current.getValue(); err != nil {
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

func TestNode_Unpack_nil(t *testing.T) {
	_, err := (*Node)(nil).Unpack()
	if err == nil {
		t.Errorf("(nil).Unpack() should be an error")
	}
}

func TestNode_getValue(t *testing.T) {
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
	iValue, err := root.getValue()
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
	iValue, err := root.getValue()
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
	if (*Node)(nil).Empty() {
		t.Errorf("(nil).Empty() is empty")
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

	root = NewNull()
	_, err = root.GetArray()
	if err == nil {
		t.Errorf("Error on root.GetArray(): NewNull")
	}
	if _, err := (*Node)(nil).GetArray(); err == nil {
		t.Errorf("(nil).GetArray() should be an error")
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

	root = NewNull()
	_, err = root.GetBool()
	if err == nil {
		t.Errorf("Error on root.GetBool(): NewNull")
	}
	if _, err := (*Node)(nil).GetBool(); err == nil {
		t.Errorf("(nil).GetBool() should be an error")
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
	if _, err := (*Node)(nil).GetIndex(0); err == nil {
		t.Errorf("(nil).GetIndex() should be an error")
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
	if _, err := (*Node)(nil).GetKey(""); err == nil {
		t.Errorf("(nil).GetKey() should be an error")
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

	root = NewNumeric(1)
	_, err = root.GetNull()
	if err == nil {
		t.Errorf("Error expected on root.GetNull() using NewNumeric")
	}
	if _, err := (*Node)(nil).GetNull(); err == nil {
		t.Errorf("(nil).GetNull() should be an error")
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

	root = NewString("")
	_, err = root.GetNumeric()
	if err == nil {
		t.Errorf("Error on root.GetNumeric() using NewString")
	}

	root = valueNode(nil, "", Numeric, "foo")
	_, err = root.GetNumeric()
	if err == nil {
		t.Errorf("Error on root.GetNumeric() wrong data")
	}
	if _, err := (*Node)(nil).GetNumeric(); err == nil {
		t.Errorf("(nil).GetNumeric() should be an error")
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

	root = NewNull()
	_, err = root.GetObject()
	if err == nil {
		t.Errorf("Error on root.GetArray(): NewNull")
	}
	if _, err := (*Node)(nil).GetObject(); err == nil {
		t.Errorf("(nil).GetObject() should be an error")
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

	root = NewNumeric(1)
	_, err = root.GetString()
	if err == nil {
		t.Errorf("Error on root.GetString(): NewNumeric")
	}
	if _, err := (*Node)(nil).GetString(); err == nil {
		t.Errorf("(nil).GetString() should be an error")
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
	if (*Node)(nil).Index() != -1 {
		t.Errorf("Wrong value for (*Node)(nil).Index()")
	}
	if NewNull().Index() != -1 {
		t.Errorf("Wrong value for Null.Index()")
	}
	if ObjectNode("", nil).Index() != -1 {
		t.Errorf("Wrong value for Null.Index()")
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
	if (*Node)(nil).Key() != "" {
		t.Errorf("Wrong value for (*Node)(nil).Key()")
	}
	if root.MustKey("foo").Clone().Key() != "" {
		t.Errorf("Wrong value for Cloned.Key()")
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
	root, err := Unmarshal([]byte(`-1.23e-2`))
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
	if (*Node)(nil).IsArray() {
		t.Errorf("Wrong (*Node)(nil).IsArray()")
	}
	if root.IsObject() {
		t.Errorf("Wrong root.IsObject()")
	}
	if (*Node)(nil).IsObject() {
		t.Errorf("Wrong (*Node)(nil).IsObject()")
	}
	if root.IsString() {
		t.Errorf("Wrong root.IsString()")
	}
	if (*Node)(nil).IsString() {
		t.Errorf("Wrong (*Node)(nil).IsString()")
	}
	if root.IsNumeric() {
		t.Errorf("Wrong root.IsNumeric()")
	}
	if (*Node)(nil).IsNumeric() {
		t.Errorf("Wrong (*Node)(nil).IsNumeric()")
	}
	if root.IsBool() {
		t.Errorf("Wrong root.IsBool()")
	}
	if (*Node)(nil).IsBool() {
		t.Errorf("Wrong (*Node)(nil).IsBool()")
	}
	if !root.IsNull() {
		t.Errorf("Wrong root.IsNull()")
	}
	if (*Node)(nil).IsNull() {
		t.Errorf("Wrong (*Node)(nil).IsNull()")
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
	if (*Node)(nil).Keys() != nil {
		t.Errorf("Wrong value for (*Node)(nil).Keys()")
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
	if (*Node)(nil).Size() != 0 {
		t.Errorf("Wrong (*Node)(nil).Size()")
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
	if (*Node)(nil).Parent() != nil {
		t.Errorf("Wrong value for (*Node)(nil).Parent()")
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
	if (*Node)(nil).Source() != nil {
		t.Errorf("Wrong value for (*Node)(nil).Source()")
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

	root = NewString("foo")
	value = root.String()
	if value != `"foo"` {
		t.Errorf("Wrong (NewString) root.String()")
	}

	root = NewNull()
	value = root.String()
	if value != "null" {
		t.Errorf("Wrong (NewNull) root.String()")
	}
	if (*Node)(nil).String() != "" {
		t.Errorf("Wrong value for (*Node)(nil).String()")
	}

	node := Must(Unmarshal([]byte(`{"foo":"bar"}`)))
	node.borders[1] = 0 // broken borders
	broken := node.String()
	if broken != "Error: not parsed yet" {
		t.Errorf("Wrong broken.String() value, actual value: %s", broken)
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

func TestNode_Type_null(t *testing.T) {
	if (*Node)(nil).Type() != Null {
		t.Errorf("Wrong value for (*Node)(nil).Type()")
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
	if (*Node)(nil).HasKey("baz") {
		t.Errorf("Wrong (*Node)(nil).HasKey('bar')")
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
	if (*Node)(nil).Path() != "" {
		t.Errorf("Wrong (nil).Path()")
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
			name:     "filled maps",
			left:     valueNode(nil, "{}", Object, map[string]*Node{"foo": NewString("bar")}),
			right:    valueNode(nil, "{}", Object, map[string]*Node{"foo": NewString("bar")}),
			expected: true,
		},
		{
			name:     "filled arrays",
			left:     valueNode(nil, "[]", Array, []*Node{NewNumeric(1)}),
			right:    valueNode(nil, "[]", Array, []*Node{NewNumeric(1)}),
			expected: true,
		},
		{
			name:     "filled maps: different",
			left:     valueNode(nil, "{}", Object, map[string]*Node{"foo": NewString("bar")}),
			right:    valueNode(nil, "{}", Object, map[string]*Node{"foo": NewString("baz")}),
			expected: false,
		},
		{
			name:     "filled maps: different keys",
			left:     valueNode(nil, "{}", Object, map[string]*Node{"foo": NewString("bar")}),
			right:    valueNode(nil, "{}", Object, map[string]*Node{"baz": NewString("bar")}),
			expected: false,
		},
		{
			name:     "filled arrays: different",
			left:     valueNode(nil, "[]", Array, []*Node{NewNumeric(1)}),
			right:    valueNode(nil, "[]", Array, []*Node{NewNumeric(2)}),
			expected: false,
		},
		{
			name:  "filled maps: errors",
			left:  valueNode(nil, "{}", Object, map[string]*Node{"foo": NewString("bar")}),
			right: valueNode(nil, "{}", Object, map[string]*Node{"foo": valueNode(nil, "", String, 123)}),
			error: true,
		},
		{
			name:  "filled arrays: errors",
			left:  valueNode(nil, "[]", Array, []*Node{NewNumeric(1)}),
			right: valueNode(nil, "[]", Array, []*Node{valueNode(nil, "", Numeric, "foo")}),
			error: true,
		},
		{
			name:     "floats 1",
			left:     NewNumeric(1.1),
			right:    NewNumeric(1.2),
			expected: false,
		},
		{
			name:     "floats 2",
			left:     NewNumeric(-1),
			right:    NewNumeric(1),
			expected: false,
		},
		{
			name:     "floats 3",
			left:     NewNumeric(1.0001),
			right:    NewNumeric(1.00011),
			expected: false,
		},
		{
			name:  "error 1",
			left:  valueNode(nil, "", Numeric, "foo"),
			right: NewNumeric(1.00011),
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
		{
			name:  "nil/value",
			left:  nil,
			right: NewString("foo"),
			error: true,
		},
		{
			name:  "value/nil",
			left:  NewString("foo"),
			right: nil,
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
		error       bool
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
			left:     NewNumeric(1.1),
			right:    NewNumeric(1.2),
			expected: true,
		},
		{
			name:     "floats 2",
			left:     NewNumeric(-1),
			right:    NewNumeric(1),
			expected: true,
		},
		{
			name:     "floats 3",
			left:     NewNumeric(1.0001),
			right:    NewNumeric(1.00011),
			expected: true,
		},
		{
			name:  "nil/value",
			left:  nil,
			right: NewString("foo"),
			error: true,
		},
		{
			name:  "value/nil",
			left:  NewString("foo"),
			right: nil,
			error: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := test.left.Neq(test.right)
			if test.error {
				if err == nil {
					t.Errorf("Error expected: nil given")
				}
			} else if err != nil {
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
			left:  NewNull(),
			right: NewNull(),
			error: true,
		},
		{
			name:  "array",
			left:  NewArray(nil),
			right: NewArray(nil),
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
			left:     NewNumeric(3.1),
			right:    NewNumeric(3),
			expected: true,
		},
		{
			name:     "float 2",
			left:     NewNumeric(0),
			right:    NewNumeric(-3),
			expected: true,
		},
		{
			name:     "float 3",
			left:     NewNumeric(0),
			right:    NewNumeric(0),
			expected: false,
		},
		{
			name:     "float 4",
			left:     NewNumeric(math.MaxFloat64),
			right:    NewNumeric(math.SmallestNonzeroFloat64),
			expected: true,
		},
		{
			name:     "float 5",
			left:     NewNumeric(math.SmallestNonzeroFloat64),
			right:    NewNumeric(math.MaxFloat64),
			expected: false,
		},
		{
			name:     "string 1",
			left:     NewString("z"),
			right:    NewString("a"),
			expected: true,
		},
		{
			name:     "string 2",
			left:     NewString("a"),
			right:    NewString("a"),
			expected: false,
		},
		{
			name:     "wrong type 1",
			left:     NewString("z"),
			right:    NewNumeric(math.MaxFloat64),
			expected: false,
		},
		{
			name:     "wrong type 2",
			left:     NewNumeric(math.MaxFloat64),
			right:    NewString("z"),
			expected: false,
		},
		{
			name:  "error 1",
			left:  valueNode(nil, "e1", Numeric, string("e1")),
			right: NewNumeric(1),
			error: true,
		},
		{
			name:  "error 2",
			left:  valueNode(nil, "e1", String, float64(1)),
			right: NewString("foo"),
			error: true,
		},
		{
			name:  "nil/value",
			left:  nil,
			right: NewString("foo"),
			error: true,
		},
		{
			name:  "value/nil",
			left:  NewString("foo"),
			right: nil,
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
			left:  NewNull(),
			right: NewNull(),
			error: true,
		},
		{
			name:  "array",
			left:  NewArray(nil),
			right: NewArray(nil),
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
			left:     NewNumeric(3.1),
			right:    NewNumeric(3),
			expected: true,
		},
		{
			name:     "float 2",
			left:     NewNumeric(0),
			right:    NewNumeric(-3),
			expected: true,
		},
		{
			name:     "float 3",
			left:     NewNumeric(0),
			right:    NewNumeric(0),
			expected: true,
		},
		{
			name:     "float 4",
			left:     NewNumeric(math.MaxFloat64),
			right:    NewNumeric(math.SmallestNonzeroFloat64),
			expected: true,
		},
		{
			name:     "float 5",
			left:     NewNumeric(math.SmallestNonzeroFloat64),
			right:    NewNumeric(math.MaxFloat64),
			expected: false,
		},
		{
			name:     "string 1",
			left:     NewString("z"),
			right:    NewString("a"),
			expected: true,
		},
		{
			name:     "string 2",
			left:     NewString("a"),
			right:    NewString("a"),
			expected: true,
		},
		{
			name:     "wrong type 1",
			left:     NewString("z"),
			right:    NewNumeric(math.MaxFloat64),
			expected: false,
		},
		{
			name:     "wrong type 2",
			left:     NewNumeric(math.MaxFloat64),
			right:    NewString("z"),
			expected: false,
		},
		{
			name:  "error 1",
			left:  valueNode(nil, "e1", Numeric, string("e1")),
			right: NewNumeric(1),
			error: true,
		},
		{
			name:  "error 2",
			left:  valueNode(nil, "e1", String, float64(1)),
			right: NewString("foo"),
			error: true,
		},
		{
			name:  "nil/value",
			left:  nil,
			right: NewString("foo"),
			error: true,
		},
		{
			name:  "value/nil",
			left:  NewString("foo"),
			right: nil,
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
			left:  NewNull(),
			right: NewNull(),
			error: true,
		},
		{
			name:  "array",
			left:  NewArray(nil),
			right: NewArray(nil),
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
			left:     NewNumeric(3.1),
			right:    NewNumeric(3),
			expected: false,
		},
		{
			name:     "float 2",
			left:     NewNumeric(0),
			right:    NewNumeric(-3),
			expected: false,
		},
		{
			name:     "float 3",
			left:     NewNumeric(0),
			right:    NewNumeric(0),
			expected: false,
		},
		{
			name:     "float 4",
			left:     NewNumeric(math.MaxFloat64),
			right:    NewNumeric(math.SmallestNonzeroFloat64),
			expected: false,
		},
		{
			name:     "float 5",
			left:     NewNumeric(math.SmallestNonzeroFloat64),
			right:    NewNumeric(math.MaxFloat64),
			expected: true,
		},
		{
			name:     "string 1",
			left:     NewString("z"),
			right:    NewString("a"),
			expected: false,
		},
		{
			name:     "string 2",
			left:     NewString("a"),
			right:    NewString("a"),
			expected: false,
		},
		{
			name:     "wrong type 1",
			left:     NewString("z"),
			right:    NewNumeric(math.MaxFloat64),
			expected: false,
		},
		{
			name:     "wrong type 2",
			left:     NewNumeric(math.MaxFloat64),
			right:    NewString("z"),
			expected: false,
		},
		{
			name:  "error 1",
			left:  valueNode(nil, "e1", Numeric, string("e1")),
			right: NewNumeric(1),
			error: true,
		},
		{
			name:  "error 2",
			left:  valueNode(nil, "e1", String, float64(1)),
			right: NewString("foo"),
			error: true,
		},
		{
			name:  "nil/value",
			left:  nil,
			right: NewString("foo"),
			error: true,
		},
		{
			name:  "value/nil",
			left:  NewString("foo"),
			right: nil,
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
			left:  NewNull(),
			right: NewNull(),
			error: true,
		},
		{
			name:  "array",
			left:  NewArray(nil),
			right: NewArray(nil),
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
			left:     NewNumeric(3.1),
			right:    NewNumeric(3),
			expected: false,
		},
		{
			name:     "float 2",
			left:     NewNumeric(0),
			right:    NewNumeric(-3),
			expected: false,
		},
		{
			name:     "float 3",
			left:     NewNumeric(0),
			right:    NewNumeric(0),
			expected: true,
		},
		{
			name:     "float 4",
			left:     NewNumeric(math.MaxFloat64),
			right:    NewNumeric(math.SmallestNonzeroFloat64),
			expected: false,
		},
		{
			name:     "float 5",
			left:     NewNumeric(math.SmallestNonzeroFloat64),
			right:    NewNumeric(math.MaxFloat64),
			expected: true,
		},
		{
			name:     "string 1",
			left:     NewString("z"),
			right:    NewString("a"),
			expected: false,
		},
		{
			name:     "string 2",
			left:     NewString("a"),
			right:    NewString("a"),
			expected: true,
		},
		{
			name:     "wrong type 1",
			left:     NewString("z"),
			right:    NewNumeric(math.MaxFloat64),
			expected: false,
		},
		{
			name:     "wrong type 2",
			left:     NewNumeric(math.MaxFloat64),
			right:    NewString("z"),
			expected: false,
		},
		{
			name:  "error 1",
			left:  valueNode(nil, "e1", Numeric, string("e1")),
			right: NewNumeric(1),
			error: true,
		},
		{
			name:  "error 2",
			left:  valueNode(nil, "e1", String, float64(1)),
			right: NewString("foo"),
			error: true,
		},
		{
			name:  "nil/value",
			left:  nil,
			right: NewString("foo"),
			error: true,
		},
		{
			name:  "value/nil",
			left:  NewString("foo"),
			right: nil,
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
	node := NewNull()
	if node.MustNull() != nil {
		t.Errorf("Failed")
	}
}

func TestNumericNode(t *testing.T) {
	node := NewNumeric(1.5)
	if node.MustNumeric() != 1.5 {
		t.Errorf("Failed")
	}
}

func TestStringNode(t *testing.T) {
	node := NewString("check")
	if node.MustString() != "check" {
		t.Errorf("Failed")
	}
}

func TestBoolNode(t *testing.T) {
	node := NewBool(true)
	if !node.MustBool() {
		t.Errorf("Failed")
	}
}

func TestArrayNode(t *testing.T) {
	array := []*Node{
		NewNull(),
		NewNumeric(1),
		NewString("foo"),
	}
	node := NewArray(array)
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
		"zero": NewNull(),
		"foo":  NewNumeric(1),
		"bar":  NewString("foo"),
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
				"zero": NewNull(),
				"foo":  NewNumeric(1),
				"bar":  NewString("foo"),
			}),
			expected: []*Node{
				NewString("foo"),
				NewNumeric(1),
				withKey(NewNull(), "0"),
			},
		},
		{
			name: "array",
			node: NewArray([]*Node{
				NewNull(),
				NewNumeric(1),
				NewString("foo"),
			}),
			expected: []*Node{
				withKey(NewNull(), "0"),
				NewNumeric(1),
				NewString("foo"),
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
	root, err := Unmarshal(jsonPathTestData)
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

func TestNode_JSONPath_error(t *testing.T) {
	root, err := Unmarshal(jsonPathTestData)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
		return
	}
	_, err = root.MustKey("store").MustKey("book").JSONPath("XXX")
	if err == nil {
		t.Errorf("JSONPath() Expected error")
	}
}

func TestNode_IsDirty(t *testing.T) {
	root, err := Unmarshal(jsonPathTestData)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
		return
	}
	tests := []struct {
		name     string
		node     *Node
		expected bool
	}{
		{
			name:     "simple",
			node:     valueNode(nil, "bool", Bool, true),
			expected: true,
		},
		{
			name:     "null",
			node:     valueNode(nil, "null", Null, nil),
			expected: true,
		},
		{
			name:     "float",
			node:     valueNode(nil, "123.5", Numeric, float64(123.5)),
			expected: true,
		},
		{
			name:     "blank array",
			node:     valueNode(nil, "[]", Array, []*Node{}),
			expected: true,
		},
		{
			name:     "blank map",
			node:     valueNode(nil, "{}", Object, map[string]*Node{}),
			expected: true,
		},
		{
			name:     "NewNumeric",
			node:     NewNumeric(1.1),
			expected: true,
		},
		{
			name:     "NewNull",
			node:     NewNull(),
			expected: true,
		},
		{
			name:     "NewArray",
			node:     NewArray(nil),
			expected: true,
		},
		{
			name:     "NewString",
			node:     NewString(""),
			expected: true,
		},
		{
			name:     "NewBool",
			node:     NewBool(false),
			expected: true,
		},
		{
			name:     "Unmarshal",
			node:     root,
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.node.dirty != test.expected {
				t.Errorf("Node dirty is not correct")
			}
		})
	}
}

func Test_newNode(t *testing.T) {
	var nilKey *string
	fillKey := "key"
	relFillKey := &fillKey
	type args struct {
		parent *Node
		buf    *buffer
		_type  NodeType
		key    **string
	}
	tests := []struct {
		name        string
		args        args
		wantCurrent *Node
		wantErr     bool
	}{
		{
			name: "blank key for Object",
			args: args{
				parent: ObjectNode("", make(map[string]*Node)),
				buf:    newBuffer(make([]byte, 10)),
				_type:  Bool,
				key:    &nilKey,
			},
			wantCurrent: nil,
			wantErr:     true,
		},
		{
			name: "child for non Object/Array",
			args: args{
				parent: NewBool(true),
				buf:    newBuffer(make([]byte, 10)),
				_type:  Bool,
				key:    &relFillKey,
			},
			wantCurrent: nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCurrent, err := newNode(tt.args.parent, tt.args.buf, tt.args._type, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("newNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(gotCurrent, tt.wantCurrent) {
				t.Errorf("newNode() gotCurrent = %v, want %v", gotCurrent, tt.wantCurrent)
			}
		})
	}
}

func TestNode_Value(t *testing.T) {
	array := NewArray([]*Node{
		NewNumeric(0),
		NewString("bar"),
	})
	object := ObjectNode("", map[string]*Node{
		"foo": NewNumeric(0),
		"bar": NewString("bar"),
	})
	tests := []struct {
		name      string
		node      *Node
		wantValue interface{}
		wantErr   bool
	}{
		{
			name:      "null",
			node:      NewNull(),
			wantValue: nil,
			wantErr:   false,
		},
		{
			name:      "string",
			node:      NewString("foo"),
			wantValue: "foo",
			wantErr:   false,
		},
		{
			name:      "string error",
			node:      valueNode(nil, "", String, false),
			wantValue: nil,
			wantErr:   true,
		},
		{
			name:      "numeric",
			node:      NewNumeric(1e3),
			wantValue: float64(1000),
			wantErr:   false,
		},
		{
			name:      "numeric error",
			node:      valueNode(nil, "", Numeric, false),
			wantValue: nil,
			wantErr:   true,
		},
		{
			name:      "bool",
			node:      NewBool(true),
			wantValue: true,
			wantErr:   false,
		},
		{
			name:      "bool error",
			node:      valueNode(nil, "", Bool, nil),
			wantValue: nil,
			wantErr:   true,
		},
		{
			name: "array",
			node: array,
			wantValue: []*Node{
				array.children["0"],
				array.children["1"],
			},
			wantErr: false,
		},
		{
			name:      "array error",
			node:      valueNode(nil, "", Array, false),
			wantValue: nil,
			wantErr:   true,
		},
		{
			name: "object",
			node: object,
			wantValue: map[string]*Node{
				"foo": object.children["foo"],
				"bar": object.children["bar"],
			},
			wantErr: false,
		},
		{
			name:      "object error",
			node:      valueNode(nil, "", Array, false),
			wantValue: nil,
			wantErr:   true,
		},
		{
			name:      "type error",
			node:      valueNode(nil, "", 10000, false),
			wantValue: nil,
			wantErr:   true,
		},
		{
			name:      "nil",
			node:      nil,
			wantValue: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, err := tt.node.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(gotValue, tt.wantValue) {
					t.Errorf("Value() gotValue = %v, want %v", gotValue, tt.wantValue)
				}
			}
		})
	}
}

func withKey(node *Node, key string) *Node {
	node.key = &key
	return node
}
