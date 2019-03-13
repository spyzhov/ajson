package ajson

import (
	"bytes"
	"testing"
)

type testCase struct {
	name  string
	input []byte
	_type NodeType
	value []byte
}

func simpleCorrupted(name string) testCase {
	return testCase{name: name, input: []byte(name)}
}

func simpleValid(test *testCase, t *testing.T) {
	root, err := Unmarshal(test.input, false)
	if err != nil {
		t.Errorf("Error on Unmarshal(%s): %s", test.name, err.Error())
	} else if root == nil {
		t.Errorf("Error on Unmarshal(%s): root is nil", test.name)
	} else if root.Type() != test._type {
		t.Errorf("Error on Unmarshal(%s): wrong type", test.name)
	} else if !bytes.Equal(root.Source(), test.value) {
		t.Errorf("Error on Unmarshal(%s): %s != %s", test.name, root.Source(), test.value)
	}
}

func simpleInvalid(test *testCase, t *testing.T) {
	root, err := Unmarshal(test.input, false)
	if err == nil {
		t.Errorf("Error on Unmarshal(%s): error expected, got '%s'", test.name, root.Source())
	} else if root != nil {
		t.Errorf("Error on Unmarshal(%s): root is not nil", test.name)
	}
}

func TestUnmarshal_NumericSimpleSuccess(t *testing.T) {
	tests := []testCase{
		{name: "1", input: []byte("1"), _type: Numeric, value: []byte("1")},
		{name: "+1", input: []byte("+1"), _type: Numeric, value: []byte("+1")},
		{name: "-1", input: []byte("-1"), _type: Numeric, value: []byte("-1")},

		{name: "1234567890", input: []byte("1234567890"), _type: Numeric, value: []byte("1234567890")},
		{name: "+123", input: []byte("+123"), _type: Numeric, value: []byte("+123")},
		{name: "-123", input: []byte("-123"), _type: Numeric, value: []byte("-123")},

		{name: "123.456", input: []byte("123.456"), _type: Numeric, value: []byte("123.456")},
		{name: "+123.456", input: []byte("+123.456"), _type: Numeric, value: []byte("+123.456")},
		{name: "-123.456", input: []byte("-123.456"), _type: Numeric, value: []byte("-123.456")},

		{name: ".456", input: []byte(".456"), _type: Numeric, value: []byte(".456")},
		{name: "+.456", input: []byte("+.456"), _type: Numeric, value: []byte("+.456")},
		{name: "-.456", input: []byte("-.456"), _type: Numeric, value: []byte("-.456")},

		{name: "1e3", input: []byte("1e3"), _type: Numeric, value: []byte("1e3")},
		{name: "1e+3", input: []byte("1e+3"), _type: Numeric, value: []byte("1e+3")},
		{name: "1e-3", input: []byte("1e-3"), _type: Numeric, value: []byte("1e-3")},
		{name: "+1e3", input: []byte("+1e3"), _type: Numeric, value: []byte("+1e3")},
		{name: "+1e+3", input: []byte("+1e+3"), _type: Numeric, value: []byte("+1e+3")},
		{name: "+1e-3", input: []byte("+1e-3"), _type: Numeric, value: []byte("+1e-3")},
		{name: "-1e3", input: []byte("-1e3"), _type: Numeric, value: []byte("-1e3")},
		{name: "-1e+3", input: []byte("-1e+3"), _type: Numeric, value: []byte("-1e+3")},
		{name: "-1e-3", input: []byte("-1e-3"), _type: Numeric, value: []byte("-1e-3")},

		{name: "1.123e3.456", input: []byte("1.123e3.456"), _type: Numeric, value: []byte("1.123e3.456")},
		{name: "1.123e+3.456", input: []byte("1.123e+3.456"), _type: Numeric, value: []byte("1.123e+3.456")},
		{name: "1.123e-3.456", input: []byte("1.123e-3.456"), _type: Numeric, value: []byte("1.123e-3.456")},
		{name: "+1.123e3.456", input: []byte("+1.123e3.456"), _type: Numeric, value: []byte("+1.123e3.456")},
		{name: "+1.123e+3.456", input: []byte("+1.123e+3.456"), _type: Numeric, value: []byte("+1.123e+3.456")},
		{name: "+1.123e-3.456", input: []byte("+1.123e-3.456"), _type: Numeric, value: []byte("+1.123e-3.456")},
		{name: "-1.123e3.456", input: []byte("-1.123e3.456"), _type: Numeric, value: []byte("-1.123e3.456")},
		{name: "-1.123e+3.456", input: []byte("-1.123e+3.456"), _type: Numeric, value: []byte("-1.123e+3.456")},
		{name: "-1.123e-3.456", input: []byte("-1.123e-3.456"), _type: Numeric, value: []byte("-1.123e-3.456")},

		{name: "1E3", input: []byte("1E3"), _type: Numeric, value: []byte("1E3")},
		{name: "1E+3", input: []byte("1E+3"), _type: Numeric, value: []byte("1E+3")},
		{name: "1E-3", input: []byte("1E-3"), _type: Numeric, value: []byte("1E-3")},
		{name: "+1E3", input: []byte("+1E3"), _type: Numeric, value: []byte("+1E3")},
		{name: "+1E+3", input: []byte("+1E+3"), _type: Numeric, value: []byte("+1E+3")},
		{name: "+1E-3", input: []byte("+1E-3"), _type: Numeric, value: []byte("+1E-3")},
		{name: "-1E3", input: []byte("-1E3"), _type: Numeric, value: []byte("-1E3")},
		{name: "-1E+3", input: []byte("-1E+3"), _type: Numeric, value: []byte("-1E+3")},
		{name: "-1E-3", input: []byte("-1E-3"), _type: Numeric, value: []byte("-1E-3")},

		{name: "1.123E3.456", input: []byte("1.123E3.456"), _type: Numeric, value: []byte("1.123E3.456")},
		{name: "1.123E+3.456", input: []byte("1.123E+3.456"), _type: Numeric, value: []byte("1.123E+3.456")},
		{name: "1.123E-3.456", input: []byte("1.123E-3.456"), _type: Numeric, value: []byte("1.123E-3.456")},
		{name: "+1.123E3.456", input: []byte("+1.123E3.456"), _type: Numeric, value: []byte("+1.123E3.456")},
		{name: "+1.123E+3.456", input: []byte("+1.123E+3.456"), _type: Numeric, value: []byte("+1.123E+3.456")},
		{name: "+1.123E-3.456", input: []byte("+1.123E-3.456"), _type: Numeric, value: []byte("+1.123E-3.456")},
		{name: "-1.123E3.456", input: []byte("-1.123E3.456"), _type: Numeric, value: []byte("-1.123E3.456")},
		{name: "-1.123E+3.456", input: []byte("-1.123E+3.456"), _type: Numeric, value: []byte("-1.123E+3.456")},
		{name: "-1.123E-3.456", input: []byte("-1.123E-3.456"), _type: Numeric, value: []byte("-1.123E-3.456")},

		{name: "-1.123E-3.456 with spaces", input: []byte(" \r -1.123E-3.456 \t\n"), _type: Numeric, value: []byte("-1.123E-3.456")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			simpleValid(&test, t)
		})
	}
}

func TestUnmarshal_NumericSimpleCorrupted(t *testing.T) {
	tests := []testCase{
		simpleCorrupted("x1"),
		simpleCorrupted("1+1"),
		simpleCorrupted("-1+"),
		simpleCorrupted("."),
		simpleCorrupted("-"),
		simpleCorrupted("+"),
		simpleCorrupted("-."),
		simpleCorrupted("+."),
		simpleCorrupted("e"),
		simpleCorrupted("e+"),
		simpleCorrupted("e+1-"),
		simpleCorrupted("1null"),
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			simpleInvalid(&test, t)
		})
	}
}

func TestUnmarshal_StringSimpleSuccess(t *testing.T) {
	tests := []testCase{
		{name: "blank", input: []byte("\"\""), _type: String, value: []byte("\"\"")},
		{name: "char", input: []byte("\"c\""), _type: String, value: []byte("\"c\"")},
		{name: "word", input: []byte("\"cat\""), _type: String, value: []byte("\"cat\"")},
		{name: "spaces", input: []byte("  \"good cat\n\tor dog\"\r\n "), _type: String, value: []byte("\"good cat\n\tor dog\"")},
		{name: "backslash", input: []byte("\"good \\\"cat\\\"\""), _type: String, value: []byte("\"good \\\"cat\\\"\"")},
		{name: "backslash 2", input: []byte("\"good \\\\\\\"cat\\\"\""), _type: String, value: []byte("\"good \\\\\\\"cat\\\"\"")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			simpleValid(&test, t)
		})
	}
}

func TestUnmarshal_StringSimpleCorrupted(t *testing.T) {
	tests := []testCase{
		{name: "one quote", input: []byte("\"")},
		{name: "one quote char", input: []byte("\"c")},
		{name: "wrong quotes", input: []byte("'cat'")},
		{name: "quotes in quotes", input: []byte("\"good \"cat\"\"")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			simpleInvalid(&test, t)
		})
	}
}

func TestUnmarshal_NullSimpleSuccess(t *testing.T) {
	tests := []testCase{
		{name: "lower", input: []byte("null"), _type: Null, value: []byte("null")},
		{name: "upper", input: []byte("NULL"), _type: Null, value: []byte("NULL")},
		{name: "CamelCase", input: []byte("NuLl"), _type: Null, value: []byte("NuLl")},
		{name: "spaces", input: []byte("  Null\r\n "), _type: Null, value: []byte("Null")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			simpleValid(&test, t)
		})
	}
}

func TestUnmarshal_NullSimpleCorrupted(t *testing.T) {
	tests := []testCase{
		{name: "nul", input: []byte("nul")},
		{name: "NILL", input: []byte("NILL")},
		{name: "spaces", input: []byte("Nu ll")},
		{name: "null1", input: []byte("null1")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			simpleInvalid(&test, t)
		})
	}
}

func TestUnmarshal_BoolSimpleSuccess(t *testing.T) {
	tests := []testCase{
		{name: "lower true", input: []byte("true"), _type: Bool, value: []byte("true")},
		{name: "lower false", input: []byte("false"), _type: Bool, value: []byte("false")},
		{name: "upper true", input: []byte("TRUE"), _type: Bool, value: []byte("TRUE")},
		{name: "upper false", input: []byte("FALSE"), _type: Bool, value: []byte("FALSE")},
		{name: "CamelCase true", input: []byte("TrUe"), _type: Bool, value: []byte("TrUe")},
		{name: "CamelCase false", input: []byte("FaLsE"), _type: Bool, value: []byte("FaLsE")},
		{name: "spaces true", input: []byte("  True\r\n "), _type: Bool, value: []byte("True")},
		{name: "spaces false", input: []byte("  False\r\n "), _type: Bool, value: []byte("False")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			simpleValid(&test, t)
		})
	}
}

func TestUnmarshal_BoolSimpleCorrupted(t *testing.T) {
	tests := []testCase{
		simpleCorrupted("tru"),
		simpleCorrupted("fals"),
		simpleCorrupted("tre"),
		simpleCorrupted("fal se"),
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			simpleInvalid(&test, t)
		})
	}
}

func TestUnmarshal_ArraySimpleSuccess(t *testing.T) {
	tests := []testCase{
		{name: "[]", input: []byte("[]"), _type: Array, value: []byte("[]")},
		{name: "[1]", input: []byte("[1]"), _type: Array, value: []byte("[1]")},
		{name: "[1,2,3]", input: []byte("[1,2,3]"), _type: Array, value: []byte("[1,2,3]")},
		{name: "[1, 2, 3]", input: []byte("[1, 2, 3]"), _type: Array, value: []byte("[1, 2, 3]")},
		{name: "[1,[2],3]", input: []byte("[1,[2],3]"), _type: Array, value: []byte("[1,[2],3]")},
		{name: "[[],[],[]]", input: []byte("[[],[],[]]"), _type: Array, value: []byte("[[],[],[]]")},
		{name: "[[[[[]]]]]", input: []byte("[[[[[]]]]]"), _type: Array, value: []byte("[[[[[]]]]]")},
		{name: "[true,null,1,\"foo\",[]]", input: []byte("[true,null,1,\"foo\",[]]"), _type: Array, value: []byte("[true,null,1,\"foo\",[]]")},
		{name: "spaces", input: []byte("\n\r [\n1\n ]\r\n"), _type: Array, value: []byte("[\n1\n ]")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			simpleValid(&test, t)
		})
	}
}

func TestUnmarshal_ArraySimpleCorrupted(t *testing.T) {
	tests := []testCase{
		simpleCorrupted("[,]"),
		simpleCorrupted("[]\\"),
		simpleCorrupted("[1,]"),
		simpleCorrupted("[[]"),
		simpleCorrupted("[]]"),
		simpleCorrupted("1[]"),
		simpleCorrupted("[]1"),
		simpleCorrupted("[[]1]"),
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			simpleInvalid(&test, t)
		})
	}
}

func TestUnmarshal_ObjectSimpleSuccess(t *testing.T) {
	tests := []testCase{
		{name: "{}", input: []byte("{}"), _type: Object, value: []byte("{}")},
		{name: `{ \r\n }`, input: []byte("{ \r\n }"), _type: Object, value: []byte("{ \r\n }")},
		{name: `{"key":1}`, input: []byte(`{"key":1}`), _type: Object, value: []byte(`{"key":1}`)},
		{name: `{"key":true}`, input: []byte(`{"key":true}`), _type: Object, value: []byte(`{"key":true}`)},
		{name: `{"key":"value"}`, input: []byte(`{"key":"value"}`), _type: Object, value: []byte(`{"key":"value"}`)},
		{name: `{"foo":"bar","baz":"foo"}`, input: []byte(`{"foo":"bar", "baz":"foo"}`), _type: Object, value: []byte(`{"foo":"bar", "baz":"foo"}`)},
		{name: "spaces", input: []byte(`  {  "foo"  :  "bar"  , "baz"   :   "foo"   }    `), _type: Object, value: []byte(`{  "foo"  :  "bar"  , "baz"   :   "foo"   }`)},
		{name: "nested", input: []byte(`{"foo":{"bar":{"baz":{}}}}`), _type: Object, value: []byte(`{"foo":{"bar":{"baz":{}}}}`)},
		{name: "array", input: []byte(`{"array":[{},{},{"foo":[{"bar":["baz"]}]}]}`), _type: Object, value: []byte(`{"array":[{},{},{"foo":[{"bar":["baz"]}]}]}`)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			simpleValid(&test, t)
		})
	}
}

func TestUnmarshal_ObjectSimpleCorrupted(t *testing.T) {
	tests := []testCase{
		simpleCorrupted("{,}"),
		simpleCorrupted("{:}"),
		simpleCorrupted(`{"foo"}`),
		simpleCorrupted(`{"foo":}`),
		simpleCorrupted(`{:"foo"}`),
		simpleCorrupted(`{"foo":bar}`),
		simpleCorrupted(`{"foo":"bar",}`),
		simpleCorrupted(`{}{}`),
		simpleCorrupted(`{},{}`),
		simpleCorrupted(`{[},{]}`),
		simpleCorrupted(`{[,]}`),
		simpleCorrupted(`{[]}`),
		simpleCorrupted(`{}1`),
		simpleCorrupted(`1{}`),
		simpleCorrupted(`{"x"::1}`),
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			simpleInvalid(&test, t)
		})
	}
}

func TestUnmarshal_Array(t *testing.T) {
	root, err := Unmarshal([]byte(" [1,[\"1\",[1,[1,2,3]]]]\r\n"), false)
	if err != nil {
		t.Errorf("Error on Unmarshal: %s", err.Error())
	} else if root == nil {
		t.Errorf("Error on Unmarshal: root is nil")
	} else if root.Type() != Array {
		t.Errorf("Error on Unmarshal: wrong type")
	} else {
		array, err := root.GetArray()
		if err != nil {
			t.Errorf("Error on root.GetArray(): %s", err.Error())
		} else if len(array) != 2 {
			t.Errorf("Error on root.GetArray(): expected 2 elements")
		} else if val, err := array[0].GetNumeric(); err != nil {
			t.Errorf("Error on array[0].GetNumeric(): %s", err.Error())
		} else if val != 1 {
			t.Errorf("Error on array[0].GetNumeric(): expected to be '1'")
		} else if val, err := array[1].GetArray(); err != nil {
			t.Errorf("Error on array[1].GetArray(): %s", err.Error())
		} else if len(val) != 2 {
			t.Errorf("Error on array[1].GetArray(): expected 2 elements")
		} else if el, err := val[0].GetString(); err != nil {
			t.Errorf("Error on val[0].GetString(): %s", err.Error())
		} else if el != "1" {
			t.Errorf("Error on val[0].GetString(): expected to be '\"1\"'")
		}
	}
}

func TestUnmarshal_Object(t *testing.T) {
	root, err := Unmarshal([]byte(`{"foo":{"bar":[null]}, "baz":true}`), false)
	if err != nil {
		t.Errorf("Error on Unmarshal: %s", err.Error())
	} else if root == nil {
		t.Errorf("Error on Unmarshal: root is nil")
	} else if !root.IsObject() {
		t.Errorf("Error on Unmarshal: wrong type")
	} else {
		object, err := root.GetObject()
		if err != nil {
			t.Errorf("Error on root.GetObject(): %s", err.Error())
		} else if foo, ok := object["foo"]; !ok {
			t.Errorf("Error on getting foo from map")
		} else if !foo.IsObject() {
			t.Errorf("Child element type error [foo]")
		} else if obj, err := foo.GetObject(); err != nil {
			t.Errorf("Error on foo.GetObject(): %s", err.Error())
		} else if bar, ok := obj["bar"]; !ok {
			t.Errorf("Error on getting bar from map")
		} else if !bar.IsArray() {
			t.Errorf("Child element type error [bar]")
		}
	}
}
