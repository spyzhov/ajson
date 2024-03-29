package ajson

import (
	"io"
	"strings"
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
		{name: "string", value: `@['foo)(]]"[[[.[?(@.bar \' == 1 && < @.length)'].*`, index: 50, fail: false},

		{name: "part 1", value: "@.foo+@.bar", index: 5, fail: false},
		{name: "part 2", value: "@.foo && @.bar", index: 5, fail: false},
		{name: "part 3", value: "@.foo,3", index: 5, fail: false},
		{name: "part 4", value: "@.length-1", index: 8, fail: false},

		{name: "number 1", value: "1", index: 1, fail: false},
		{name: "number 2", value: "1.3e2", index: 5, fail: false},
		{name: "number 3", value: "-1.3e2", index: 6, fail: false},
		{name: "number 4", value: "-1.3e-2", index: 7, fail: false},

		{name: "string 1", value: "'1'", index: 3, fail: false},
		{name: "string 2", value: "'foo \\'bar '", index: 12, fail: false},
		{name: "string 3", value: `"foo \"bar "`, index: 12, fail: false},

		{name: "fail 1", value: "@.foo[", fail: true},
		{name: "fail 2", value: "@.foo[(]", fail: true},
		{name: "fail 3", value: "'", fail: true},
		{name: "fail 4", value: "'x", fail: true},

		{name: "parentheses 0", value: "()", index: 2, fail: false},
		{name: "parentheses 1", value: "(@)", index: 3, fail: false},
		{name: "parentheses 2", value: "(", fail: true},
		{name: "parentheses 3", value: ")", fail: true},
		{name: "parentheses 4", value: "(x", fail: true},
		{name: "parentheses 5", value: "((())", fail: true},
		{name: "parentheses 6", value: "@)", index: 1, fail: false},
		{name: "parentheses 7", value: "[)", fail: true},
		{name: "parentheses 8", value: "[())", fail: true},

		{name: "bracket 0", value: "[]", index: 2, fail: false},
		{name: "bracket 1", value: "[@]", index: 3, fail: false},
		{name: "bracket 2", value: "[", fail: true},
		{name: "bracket 3", value: "]", fail: true},
		{name: "bracket 4", value: "[x", fail: true},
		{name: "bracket 5", value: "[[[]]", fail: true},
		{name: "bracket 6", value: "@]", index: 1, fail: false},
		{name: "bracket 7", value: "(]", fail: true},
		{name: "bracket 8", value: "([]]", fail: true},

		{name: "sign 1", value: "+X", index: 1},
		{name: "sign 2", value: "-X", index: 1},
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

func TestBuffer_RPN_long_operations_name(t *testing.T) {
	jsonpath := `@.key !@#$%^&* 1`
	_, err := newBuffer([]byte(jsonpath)).rpn()
	if err == nil {
		t.Errorf("Expected error, got nothing")
		return
	}

	expected := []string{"@.key", "1", "!@#$%^&*"}
	AddOperation(`!@#$%^&*`, 3, false, func(left *Node, right *Node) (result *Node, err error) {
		return NullNode(""), nil
	})
	result, err := newBuffer([]byte(jsonpath)).rpn()
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	} else if !sliceEqual(expected, result) {
		t.Errorf("Error on RPN(%s): result doesn't match\nExpected: %s\nActual:   %s", jsonpath, sliceString(expected), sliceString(result))
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
		{name: "example_8", value: "@.length-1", expected: []string{"@.length", "1", "-"}},
		{name: "example_9", value: "@.length+-1", expected: []string{"@.length", "-1", "+"}},
		{name: "example_10", value: "@.length/e", expected: []string{"@.length", "e", "/"}},
		{name: "example_12", value: "123.456", expected: []string{"123.456"}},
		{name: "example_13", value: " 123.456 ", expected: []string{"123.456"}},

		{name: "1 /", value: "1 /", expected: []string{"1", "/"}},
		{name: "1 + ", value: "1 + ", expected: []string{"1", "+"}},
		{name: "1 -", value: "1 -", expected: []string{"1", "-"}},
		{name: "1 * ", value: "1 * ", expected: []string{"1", "*"}},
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

func TestBuffer_RPNError(t *testing.T) {
	tests := []struct {
		value string
	}{
		{value: "1 + / 1"},
		{value: "1 * / 1"},
		{value: "1 - / 1"},
		{value: "1 / / 1"},

		{value: "1 + * 1"},
		{value: "1 * * 1"},
		{value: "1 - * 1"},
		{value: "1 / * 1"},

		{value: "1e1.1 + 1"},

		{value: "len('string)"},
		{value: "'Hello ' + 'World"},

		{value: "@.length + $['length')"},
		{value: "2 + 2)"},
		{value: "(2 + 2"},

		{value: "e + q"},
		{value: "foo(e)"},
		{value: "++2"},
		{value: ""},
	}
	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			buf := newBuffer([]byte(test.value))
			result, err := buf.rpn()
			if err == nil {
				t.Errorf("Expected error, nil given, with result: %v", strings.Join(result, ", "))
			}
		})
	}
}

func TestTokenize(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected []string
		fail     bool
	}{
		{name: "example_1", value: "@.length", expected: []string{"@.length"}},
		{name: "example_2", value: "1 + 2", expected: []string{"1", "+", "2"}},
		{name: "example_3", value: "1+2", expected: []string{"1", "+", "2"}},
		{name: "example_4", value: "1:", expected: []string{"1", ":"}},
		{name: "example_5", value: ":2 :1", expected: []string{":", "2", ":", "1"}},
		{name: "example_6", value: "1 ,2,'foo'", expected: []string{"1", ",", "2", ",", "'foo'"}},
		{name: "example_7", value: "(@.length-1)", expected: []string{"(", "@.length", "-", "1", ")"}},
		{name: "example_8", value: "?(@.length-1)", expected: []string{"?", "(", "@.length", "-", "1", ")"}},
		{name: "example_9", value: "'foo'", expected: []string{"'foo'"}},
		{name: "example_10", value: "$.foo[(@.length - 3):3:]", expected: []string{"$.foo[(@.length - 3):3:]"}},
		{name: "example_11", value: "$..", expected: []string{"$.."}},
		{name: "blank", value: "", expected: []string{}},
		{name: "number", value: "1e", fail: true},
		{name: "string", value: "'foo", fail: true},
		{name: "fail", value: "@.[", fail: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := tokenize(test.value)
			if test.fail {
				if err == nil {
					t.Error("Expected error: nil given")
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
			} else if !sliceEqual(test.expected, result) {
				t.Errorf("Error on RPN(%s): result doesn't match\nExpected: %s\nActual:   %s", test.value, sliceString(test.expected), sliceString(result))
			}
		})
	}
}

func TestBuffer_Current(t *testing.T) {
	buf := newBuffer([]byte{})
	_, err := buf.current()
	if err != io.EOF {
		t.Error("Unexpected result: io.EOF expected")
	}
}

func TestBuffer_Numeric(t *testing.T) {
	tests := []struct {
		value string
		index int
		fail  bool
	}{
		{value: "1", index: 1, fail: false},
		{value: "0", index: 1, fail: false},
		{value: "1.3e2", index: 5, fail: false},
		{value: "-1.3e2", index: 6, fail: false},
		{value: "-1.3e-2", index: 7, fail: false},
		{value: "..3", index: 0, fail: true},
		{value: "e.", index: 0, fail: true},
		{value: ".e.", index: 0, fail: true},
		{value: "1.e1", index: 0, fail: true},
		{value: "0.e0", index: 0, fail: true},
		{value: "0+0", index: 1, fail: false},
		{value: "0-1", index: 1, fail: false},
		{value: "++1", index: 0, fail: true},
		{value: "--1", index: 0, fail: true},
		{value: "-+1", index: 0, fail: true},
		{value: "+-1", index: 0, fail: true},
		{value: "+", index: 0, fail: true},
		{value: "-", index: 0, fail: true},
		{value: ".", index: 0, fail: true},
		{value: "e", index: 0, fail: true},
		{value: "+a", index: 0, fail: true},
	}
	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			buf := newBuffer([]byte(test.value))
			err := buf.numeric(true)
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
