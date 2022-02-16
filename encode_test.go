package ajson

import (
	"fmt"
	"testing"
)

func ExampleMarshal() {
	data := []byte(`[{"latitude":1,"longitude":2},{"other":"value"},null,{"internal":{"name": "unknown", "longitude":22, "latitude":11}}]`)
	root := Must(Unmarshal(data))
	locations, _ := root.JSONPath("$..[?(@.latitude && @.longitude)]")
	for _, location := range locations {
		name := fmt.Sprintf("At [%v, %v]", location.MustKey("latitude").MustNumeric(), location.MustKey("longitude").MustNumeric())
		_ = location.AppendObject("name", NewString(name))
	}
	result, _ := Marshal(root)
	fmt.Printf("%s", result)
	// JSON Output:
	// [
	// 	{
	// 		"latitude":1,
	// 		"longitude":2,
	// 		"name":"At [1, 2]"
	// 	},
	// 	{
	// 		"other":"value"
	// 	},
	// 	null,
	// 	{
	// 		"internal":{
	// 			"name":"At [11, 22]",
	// 			"longitude":22,
	// 			"latitude":11
	// 		}
	// 	}
	// ]
}

func TestMarshal_Primitive(t *testing.T) {
	tests := []struct {
		name string
		node *Node
	}{
		{
			name: "null",
			node: NewNull(),
		},
		{
			name: "true",
			node: NewBool(true),
		},
		{
			name: "false",
			node: NewBool(false),
		},
		{
			name: `"string"`,
			node: NewString("string"),
		},
		{
			name: `"one \"encoded\" string"`,
			node: NewString(`one "encoded" string`),
		},
		{
			name: `"spec.symbols: \r\n\t; UTF-8: 😹; \u2028 \u0000"`,
			node: NewString("spec.symbols: \r\n\t; UTF-8: 😹; \u2028 \000"),
		},
		{
			name: "100500",
			node: NewNumeric(100500),
		},
		{
			name: "100.5",
			node: NewNumeric(100.5),
		},
		{
			name: "[1,2,3]",
			node: NewArray([]*Node{
				NewNumeric(1),
				NewNumeric(2),
				NewNumeric(3),
			}),
		},
		{
			name: `{"foo":"bar"}`,
			node: ObjectNode("", map[string]*Node{
				"foo": NewString("bar"),
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

func TestMarshal_Unparsed(t *testing.T) {
	node := Must(Unmarshal([]byte(`{"foo":"bar"}`)))
	node.borders[1] = 0 // broken borders
	_, err := Marshal(node)
	if err == nil {
		t.Errorf("expected error")
	} else if current, ok := err.(Error); !ok {
		t.Errorf("unexpected error type: %T %s", err, err)
	} else if current.Error() != "not parsed yet" {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestMarshal_Encoded(t *testing.T) {
	base := `"one \"encoded\" string"`
	node := Must(Unmarshal([]byte(base)))

	value, err := Marshal(node)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else if string(value) != base {
		t.Errorf("wrong result: '%s', expected '%s'", value, base)
	}
}

func TestMarshal_Errors(t *testing.T) {
	tests := []struct {
		name string
		node func() (node *Node)
	}{
		{
			name: "nil",
			node: func() (node *Node) {
				return
			},
		},
		{
			name: "broken",
			node: func() (node *Node) {
				node = Must(Unmarshal([]byte(`{}`)))
				node.borders[1] = 0
				return
			},
		},
		{
			name: "Numeric",
			node: func() (node *Node) {
				return valueNode(nil, "", Numeric, false)
			},
		},
		{
			name: "String",
			node: func() (node *Node) {
				return valueNode(nil, "", String, false)
			},
		},
		{
			name: "Bool",
			node: func() (node *Node) {
				return valueNode(nil, "", Bool, 1)
			},
		},
		{
			name: "Array_1",
			node: func() (node *Node) {
				node = NewArray(nil)
				node.children["1"] = withKey(NewNull(), "1")
				return
			},
		},
		{
			name: "Array_2",
			node: func() (node *Node) {
				return NewArray([]*Node{valueNode(nil, "", Bool, 1)})
			},
		},
		{
			name: "Object",
			node: func() (node *Node) {
				return ObjectNode("", map[string]*Node{"key": valueNode(nil, "key", Bool, 1)})
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, err := Marshal(test.node())
			if err == nil {
				t.Errorf("expected error")
			} else if len(value) != 0 {
				t.Errorf("wrong result")
			}
		})
	}
}
