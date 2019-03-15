package ajson

import (
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
			current := &node{
				_type:   test._type,
				borders: [2]int{0, len(test.bytes)},
				data:    &test.bytes,
			}
			value, err := current.Value()
			if err != nil {
				t.Errorf("Error on get value: %s", err.Error())
			} else if value != test.expected {
				t.Errorf("Error on get value: '%v' != '%v'", value, test.expected)
			}
		})
	}
}
