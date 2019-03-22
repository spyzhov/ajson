package ajson

import (
	"io"
	"testing"
)

func TestBuffer_Token(t *testing.T) {
	tests := []struct {
		name  string
		value string
		index int
		fail  bool
	}{
		{name: "simple", value: "@.length", index: 8, fail: false},
		{name: "combined", value: "@['foo'].0.bar", index: 14, fail: false},
		{name: "formula", value: "@['foo'].[(@.length - 1)].*", index: 27, fail: false},
		{name: "filter", value: "@['foo'].[?(@.bar == 1 && @.baz < @.length)].*", index: 46, fail: false},
		{name: "string", value: "@['foo)(]][[[.[?(@.bar \\' == 1 && < @.length)'].*", index: 49, fail: false},

		{name: "part 1", value: "@.foo+@.bar", index: 5, fail: false},
		{name: "part 2", value: "@.foo && @.bar", index: 5, fail: false},

		{name: "fail 1", value: "@.foo[", fail: true},
		{name: "fail 2", value: "@.foo[(]", fail: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := newBuffer([]byte(test.value))
			err := buf.token()
			if !test.fail && err != nil && err != io.EOF {
				t.Errorf("Unexpected error: %s", err.Error())
			} else if test.fail && (err == nil || err == io.EOF) {
				t.Errorf("Expected error, got nothing")
			} else if !test.fail && test.index != buf.index {
				t.Errorf("Wrong index: expected %d, got %d", test.index, buf.index)
			}
		})
	}
}

func TestBuffer_RPN(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected []string
	}{
		{name: "example_1", value: "@.length", expected: []string{"@.length"}},
		{name: "example_2", value: "1 + 2", expected: []string{"1", "2", "+"}},
		{name: "example_3", value: "3 + 4 * 2 / (1 - 5)**2", expected: []string{"3", "4", "2", "*", "1", "5", "-", "2", "**", "/", "+"}},
		{name: "example_4", value: "'foo' == pi", expected: []string{"'foo'", "pi", "=="}},
		{name: "example_5", value: "pi != 'bar'", expected: []string{"pi", "'bar'", "!="}},
		{name: "example_6", value: "3 + 4 * -2 / (-1 - 5)**-2", expected: []string{"3", "4", "-2", "*", "-1", "5", "-", "-2", "**", "/", "+"}},
		{name: "example_7", value: "1.3e2 + sin(2*pi/3)", expected: []string{"1.3e2", "2", "pi", "*", "3", "/", "sin", "+"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := newBuffer([]byte(test.value))
			result, err := buf.rpn()
			if err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
			} else if !sliceEqual(test.expected, result) {
				t.Errorf("Error on RPN(%s): result doesn't match\nExpected: %s\nActual:   %s", test.value, sliceString(test.expected), sliceString(result))
			}
		})
	}
}
