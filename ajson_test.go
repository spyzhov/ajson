package ajson

import (
	"bytes"
	"testing"
)

func TestNumericSimpleSuccess(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		_type NodeType
		value []byte
	}{
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
			root, err := Unmarshal(test.input, false)
			if err != nil {
				t.Errorf("Error on Unmarshal(%s): %s", test.name, err.Error())
			} else if root == nil {
				t.Errorf("Error on Unmarshal(%s): root is nil", test.name)
			} else if root.Type() != test._type {
				t.Errorf("Error on Unmarshal(%s): wrong type", test.name)
			} else if !bytes.Equal(root.Value(), test.value) {
				t.Errorf("Error on Unmarshal(%s): %s != %s", test.name, root.Value(), test.value)
			}
		})
	}
}

func TestNumericSimpleCorrupted(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{name: "x1", input: []byte("x1")},
		{name: "1+1", input: []byte("1+1")},
		{name: "-1+", input: []byte("-1+")},
		{name: ".", input: []byte(".")},
		{name: "-", input: []byte("-")},
		{name: "+", input: []byte("+")},
		{name: "-.", input: []byte("-")},
		{name: "+.", input: []byte("+")},
		{name: "e", input: []byte("e")},
		{name: "e+", input: []byte("e+")},
		{name: "e+1-", input: []byte("e+1-")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root, err := Unmarshal(test.input, false)
			if err == nil {
				t.Errorf("Error on Unmarshal(%s): error expected", test.name)
			} else if root != nil {
				t.Errorf("Error on Unmarshal(%s): root is not nil", test.name)
			}
		})
	}
}

func TestStringSimpleSuccess(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		_type NodeType
		value []byte
	}{
		{name: "blank", input: []byte("\"\""), _type: String, value: []byte("\"\"")},
		{name: "char", input: []byte("\"c\""), _type: String, value: []byte("\"c\"")},
		{name: "word", input: []byte("\"cat\""), _type: String, value: []byte("\"cat\"")},
		{name: "spaces", input: []byte("  \"good cat\n\tor dog\"\r\n "), _type: String, value: []byte("\"good cat\n\tor dog\"")},
		{name: "backslash", input: []byte("\"good \\\"cat\\\"\""), _type: String, value: []byte("\"good \\\"cat\\\"\"")},
		{name: "backslash 2", input: []byte("\"good \\\\\\\"cat\\\"\""), _type: String, value: []byte("\"good \\\\\\\"cat\\\"\"")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root, err := Unmarshal(test.input, false)
			if err != nil {
				t.Errorf("Error on Unmarshal(%s): %s", test.name, err.Error())
			} else if root == nil {
				t.Errorf("Error on Unmarshal(%s): root is nil", test.name)
			} else if root.Type() != test._type {
				t.Errorf("Error on Unmarshal(%s): wrong type", test.name)
			} else if !bytes.Equal(root.Value(), test.value) {
				t.Errorf("Error on Unmarshal(%s): %s != %s", test.name, root.Value(), test.value)
			}
		})
	}
}

func TestStringSimpleCorrupted(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{name: "one quote", input: []byte("\"")},
		{name: "one quote char", input: []byte("\"c")},
		{name: "wrong quotes", input: []byte("'cat'")},
		{name: "quotes in quotes", input: []byte("\"good \"cat\"\"")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root, err := Unmarshal(test.input, false)
			if err == nil {
				t.Errorf("Error on Unmarshal(%s): error expected", test.name)
			} else if root != nil {
				t.Errorf("Error on Unmarshal(%s): root is not nil", test.name)
			}
		})
	}
}
