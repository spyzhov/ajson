package ajson

import "testing"

func TestMarshal_Primitive(t *testing.T) {
	tests := []struct {
		name   string
		node   *Node
		result string
	}{
		{
			name: "null",
			node: NullNode(""),
		},
		{
			name: "true",
			node: BoolNode("", true),
		},
		{
			name: "false",
			node: BoolNode("", false),
		},
		{
			name: `"string"`,
			node: StringNode("", "string"),
		},
		{
			name: "100500",
			node: NumericNode("", 100500),
		},
		{
			name: "100.5",
			node: NumericNode("", 100.5),
		},
		{
			name: "[1,2,3]",
			node: ArrayNode("", []*Node{
				NumericNode("0", 1),
				NumericNode("2", 2),
				NumericNode("3", 3),
			}),
		},
		{
			name: `{"foo":"bar"}`,
			node: ObjectNode("", map[string]*Node{
				"foo": StringNode("foo", "bar"),
			}),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, err := Marshal(test.node)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			} else if string(value) != test.name {
				t.Errorf("wrong result: '%s', expected '%s'", value, test.name)
			}
		})
	}
}
