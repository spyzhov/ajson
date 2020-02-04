package ajson

import (
	"reflect"
	"testing"
)

// Add test: [{"foo":"bar"}], get "bar", parent=nil => get value/path/etc.

func TestNode_SetNull(t *testing.T) {
	tests := []struct {
		name string
		node *Node
	}{
		{
			name: "Null",
			node: NullNode(""),
		},
		{
			name: "parsed Null",
			node: Must(Unmarshal(_null)),
		},
		{
			name: "Bool",
			node: BoolNode("", false),
		},
		{
			name: "parsed Bool",
			node: Must(Unmarshal(_true)),
		},
		{
			name: "String",
			node: StringNode("", "String value"),
		},
		{
			name: "parsed String",
			node: Must(Unmarshal([]byte(`"some String value"`))),
		},
		{
			name: "Numeric",
			node: NumericNode("", 123.456),
		},
		{
			name: "parsed Numeric",
			node: Must(Unmarshal([]byte(`123.456`))),
		},
		{
			name: "Array",
			node: ArrayNode("", []*Node{
				NumericNode("0", 123.456),
				BoolNode("1", false),
				NullNode("2"),
			}),
		},
		{
			name: "parsed Array",
			node: Must(Unmarshal([]byte(`[0,1,null,true,"example"]`))),
		},
		{
			name: "Object",
			node: ObjectNode("", map[string]*Node{
				"foo": NumericNode("foo", 123.456),
				"bar": BoolNode("bar", false),
				"baz": NullNode("baz"),
			}),
		},
		{
			name: "parsed Object",
			node: Must(Unmarshal([]byte(`{"foo": true, "bar": null, "baz": 123}`))),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.node.SetNull()
			if value, err := test.node.GetNull(); err != nil {
				t.Errorf("GetNull returns error: %s", err)
			} else if value != nil {
				t.Errorf("GetNull returns wrong value: %v != nil", value)
			}
			if test.node.ready() {
				t.Errorf("modified Node is ready")
			}
			if !test.node.IsDirty() {
				t.Errorf("modified Node is not dirty")
			}
			if test.node.children != nil {
				t.Errorf("modified Node has children")
			}
		})
	}
}

func TestNode_SetNumeric(t *testing.T) {
	expected := 123.456
	tests := []struct {
		name string
		node *Node
	}{
		{
			name: "Null",
			node: NullNode(""),
		},
		{
			name: "parsed Null",
			node: Must(Unmarshal(_null)),
		},
		{
			name: "Bool",
			node: BoolNode("", false),
		},
		{
			name: "parsed Bool",
			node: Must(Unmarshal(_true)),
		},
		{
			name: "String",
			node: StringNode("", "String value"),
		},
		{
			name: "parsed String",
			node: Must(Unmarshal([]byte(`"some String value"`))),
		},
		{
			name: "Numeric",
			node: NumericNode("", 123.456),
		},
		{
			name: "parsed Numeric",
			node: Must(Unmarshal([]byte(`123.456`))),
		},
		{
			name: "Array",
			node: ArrayNode("", []*Node{
				NumericNode("0", 123.456),
				BoolNode("1", false),
				NullNode("2"),
			}),
		},
		{
			name: "parsed Array",
			node: Must(Unmarshal([]byte(`[0,1,null,true,"example"]`))),
		},
		{
			name: "Object",
			node: ObjectNode("", map[string]*Node{
				"foo": NumericNode("foo", 123.456),
				"bar": BoolNode("bar", false),
				"baz": NullNode("baz"),
			}),
		},
		{
			name: "parsed Object",
			node: Must(Unmarshal([]byte(`{"foo": true, "bar": null, "baz": 123}`))),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.node.SetNumeric(expected)
			if value, err := test.node.GetNumeric(); err != nil {
				t.Errorf("GetNumeric returns error: %s", err)
			} else if value != expected {
				t.Errorf("GetNumeric returns wrong value: %v != %v", value, expected)
			}
			if test.node.ready() {
				t.Errorf("modified Node is ready")
			}
			if !test.node.IsDirty() {
				t.Errorf("modified Node is not dirty")
			}
			if test.node.children != nil {
				t.Errorf("modified Node has children")
			}
		})
	}
}

func TestNode_SetString(t *testing.T) {
	expected := "expected value"
	tests := []struct {
		name string
		node *Node
	}{
		{
			name: "Null",
			node: NullNode(""),
		},
		{
			name: "parsed Null",
			node: Must(Unmarshal(_null)),
		},
		{
			name: "Bool",
			node: BoolNode("", false),
		},
		{
			name: "parsed Bool",
			node: Must(Unmarshal(_true)),
		},
		{
			name: "String",
			node: StringNode("", "String value"),
		},
		{
			name: "parsed String",
			node: Must(Unmarshal([]byte(`"some String value"`))),
		},
		{
			name: "Numeric",
			node: NumericNode("", 123.456),
		},
		{
			name: "parsed Numeric",
			node: Must(Unmarshal([]byte(`123.456`))),
		},
		{
			name: "Array",
			node: ArrayNode("", []*Node{
				NumericNode("0", 123.456),
				BoolNode("1", false),
				NullNode("2"),
			}),
		},
		{
			name: "parsed Array",
			node: Must(Unmarshal([]byte(`[0,1,null,true,"example"]`))),
		},
		{
			name: "Object",
			node: ObjectNode("", map[string]*Node{
				"foo": NumericNode("foo", 123.456),
				"bar": BoolNode("bar", false),
				"baz": NullNode("baz"),
			}),
		},
		{
			name: "parsed Object",
			node: Must(Unmarshal([]byte(`{"foo": true, "bar": null, "baz": 123}`))),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.node.SetString(expected)
			if value, err := test.node.GetString(); err != nil {
				t.Errorf("GetString returns error: %s", err)
			} else if value != expected {
				t.Errorf("GetString returns wrong value: %v != %v", value, expected)
			}
			if test.node.ready() {
				t.Errorf("modified Node is ready")
			}
			if !test.node.IsDirty() {
				t.Errorf("modified Node is not dirty")
			}
			if test.node.children != nil {
				t.Errorf("modified Node has children")
			}
		})
	}
}

func TestNode_SetBool(t *testing.T) {
	expected := true
	tests := []struct {
		name string
		node *Node
	}{
		{
			name: "Null",
			node: NullNode(""),
		},
		{
			name: "parsed Null",
			node: Must(Unmarshal(_null)),
		},
		{
			name: "Bool",
			node: BoolNode("", false),
		},
		{
			name: "parsed Bool",
			node: Must(Unmarshal(_true)),
		},
		{
			name: "String",
			node: StringNode("", "String value"),
		},
		{
			name: "parsed String",
			node: Must(Unmarshal([]byte(`"some String value"`))),
		},
		{
			name: "Numeric",
			node: NumericNode("", 123.456),
		},
		{
			name: "parsed Numeric",
			node: Must(Unmarshal([]byte(`123.456`))),
		},
		{
			name: "Array",
			node: ArrayNode("", []*Node{
				NumericNode("0", 123.456),
				BoolNode("1", false),
				NullNode("2"),
			}),
		},
		{
			name: "parsed Array",
			node: Must(Unmarshal([]byte(`[0,1,null,true,"example"]`))),
		},
		{
			name: "Object",
			node: ObjectNode("", map[string]*Node{
				"foo": NumericNode("foo", 123.456),
				"bar": BoolNode("bar", false),
				"baz": NullNode("baz"),
			}),
		},
		{
			name: "parsed Object",
			node: Must(Unmarshal([]byte(`{"foo": true, "bar": null, "baz": 123}`))),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.node.SetBool(expected)
			if value, err := test.node.GetBool(); err != nil {
				t.Errorf("GetBool returns error: %s", err)
			} else if value != expected {
				t.Errorf("GetBool returns wrong value: %v != %v", value, expected)
			}
			if test.node.ready() {
				t.Errorf("modified Node is ready")
			}
			if !test.node.IsDirty() {
				t.Errorf("modified Node is not dirty")
			}
			if test.node.children != nil {
				t.Errorf("modified Node has children")
			}
		})
	}
}

func TestNode_SetArray(t *testing.T) {
	expected := []*Node{
		NullNode("0"),
		BoolNode("1", false),
		StringNode("2", "Foo"),
		NumericNode("3", 1),
	}
	tests := []struct {
		name string
		node *Node
	}{
		{
			name: "Null",
			node: NullNode(""),
		},
		{
			name: "parsed Null",
			node: Must(Unmarshal(_null)),
		},
		{
			name: "Bool",
			node: BoolNode("", false),
		},
		{
			name: "parsed Bool",
			node: Must(Unmarshal(_true)),
		},
		{
			name: "String",
			node: StringNode("", "String value"),
		},
		{
			name: "parsed String",
			node: Must(Unmarshal([]byte(`"some String value"`))),
		},
		{
			name: "Numeric",
			node: NumericNode("", 123.456),
		},
		{
			name: "parsed Numeric",
			node: Must(Unmarshal([]byte(`123.456`))),
		},
		{
			name: "Array",
			node: ArrayNode("", []*Node{
				NumericNode("0", 123.456),
				BoolNode("1", false),
				NullNode("2"),
			}),
		},
		{
			name: "parsed Array",
			node: Must(Unmarshal([]byte(`[0,1,null,true,"example"]`))),
		},
		{
			name: "Object",
			node: ObjectNode("", map[string]*Node{
				"foo": NumericNode("foo", 123.456),
				"bar": BoolNode("bar", false),
				"baz": NullNode("baz"),
			}),
		},
		{
			name: "parsed Object",
			node: Must(Unmarshal([]byte(`{"foo": true, "bar": null, "baz": 123}`))),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.node.SetArray(expected)
			if value, err := test.node.GetArray(); err != nil {
				t.Errorf("GetArray returns error: %s", err)
			} else if !reflect.DeepEqual(value, expected) {
				t.Errorf("GetArray returns wrong value: %v != %v", value, expected)
			}
			if test.node.ready() {
				t.Errorf("modified Node is ready")
			}
			if !test.node.IsDirty() {
				t.Errorf("modified Node is not dirty")
			}
			if test.node.children != nil {
				t.Errorf("modified Node has children")
			}
		})
	}
}

func TestNode_SetObject(t *testing.T) {
	expected := map[string]*Node{
		"foo": NullNode("foo"),
		"bar": BoolNode("bar", false),
	}
	tests := []struct {
		name string
		node *Node
	}{
		{
			name: "Null",
			node: NullNode(""),
		},
		{
			name: "parsed Null",
			node: Must(Unmarshal(_null)),
		},
		{
			name: "Bool",
			node: BoolNode("", false),
		},
		{
			name: "parsed Bool",
			node: Must(Unmarshal(_true)),
		},
		{
			name: "String",
			node: StringNode("", "String value"),
		},
		{
			name: "parsed String",
			node: Must(Unmarshal([]byte(`"some String value"`))),
		},
		{
			name: "Numeric",
			node: NumericNode("", 123.456),
		},
		{
			name: "parsed Numeric",
			node: Must(Unmarshal([]byte(`123.456`))),
		},
		{
			name: "Array",
			node: ArrayNode("", []*Node{
				NumericNode("0", 123.456),
				BoolNode("1", false),
				NullNode("2"),
			}),
		},
		{
			name: "parsed Array",
			node: Must(Unmarshal([]byte(`[0,1,null,true,"example"]`))),
		},
		{
			name: "Object",
			node: ObjectNode("", map[string]*Node{
				"foo": NumericNode("foo", 123.456),
				"bar": BoolNode("bar", false),
				"baz": NullNode("baz"),
			}),
		},
		{
			name: "parsed Object",
			node: Must(Unmarshal([]byte(`{"foo": true, "bar": null, "baz": 123}`))),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.node.SetObject(expected)
			if value, err := test.node.GetObject(); err != nil {
				t.Errorf("GetObject returns error: %s", err)
			} else if !reflect.DeepEqual(value, expected) {
				t.Errorf("GetObject returns wrong value: %v != %v", value, expected)
			}
			if test.node.ready() {
				t.Errorf("modified Node is ready")
			}
			if !test.node.IsDirty() {
				t.Errorf("modified Node is not dirty")
			}
			if test.node.children != nil {
				t.Errorf("modified Node has children")
			}
		})
	}
}

func TestNode_mutations(t *testing.T) {
	root := Must(Unmarshal([]byte(`[{"foo":"bar"}]`)))
	nodes, err := root.JSONPath(`$..foo`)
	if err != nil {
		t.Errorf("JSONPath returns error: %s", err)
	} else if len(nodes) != 1 {
		t.Errorf("JSONPath returns wrong result size: %d", len(nodes))
	}

	node := nodes[0]
	if value, err := Marshal(node); err != nil {
		t.Errorf("Marshal returns error: %s", err)
	} else if string(value) != `"bar"` {
		t.Errorf("Marshal returns wrong value: %v", string(value))
	}

	root.SetNull()
	if value, err := Marshal(node); err != nil {
		t.Errorf("Marshal returns error: %s", err)
	} else if string(value) != `"bar"` {
		t.Errorf("Marshal returns wrong value: %v", string(value))
	}

	newRoot := node.root()
	if value, err := Marshal(newRoot); err != nil {
		t.Errorf("Marshal returns error: %s", err)
	} else if string(value) != `{"foo":"bar"}` {
		t.Errorf("Marshal returns wrong value: %v", string(value))
	}

	newRoot.SetNull()
	if value, err := Marshal(node); err != nil {
		t.Errorf("Marshal returns error: %s", err)
	} else if string(value) != `"bar"` {
		t.Errorf("Marshal returns wrong value: %v", string(value))
	}

	lastRoot := node.root()
	if value, err := Marshal(lastRoot); err != nil {
		t.Errorf("Marshal returns error: %s", err)
	} else if string(value) != `"bar"` {
		t.Errorf("Marshal returns wrong value: %v", string(value))
	}
}

func TestNode_AppendArray(t *testing.T) {
	if err := Must(Unmarshal([]byte(`[{"foo":"bar"}]`))).AppendArray(NullNode("")); err != nil {
		t.Errorf("AppendArray should return error")
	}

	root := Must(Unmarshal([]byte(`[{"foo":"bar"}]`)))

	if err := root.AppendArray(NullNode("")); err != nil {
		t.Errorf("AppendArray returns error: %v", err)
	}
	if value, err := Marshal(root); err != nil {
		t.Errorf("Marshal returns error: %v", err)
	} else if string(value) != `[{"foo":"bar"},null]` {
		t.Errorf("Marshal returns wrong value: %s", string(value))
	}

	if err := root.AppendArray(
		NumericNode("", 1),
		StringNode("", "foo"),
		Must(Unmarshal([]byte(`[0,1,null,true,"example"]`))),
		Must(Unmarshal([]byte(`{"foo": true, "bar": null, "baz": 123}`))),
	); err != nil {
		t.Errorf("AppendArray returns error: %v", err)
	}
	if value, err := Marshal(root); err != nil {
		t.Errorf("Marshal returns error: %v", err)
	} else if string(value) != `[{"foo":"bar"},null,1,"foo",[0,1,null,true,"example"],{"foo": true, "bar": null, "baz": 123}]` {
		t.Errorf("Marshal returns wrong value: %s", string(value))
	}
}

func TestNode_AppendArray_self(t *testing.T) {
	root := Must(Unmarshal([]byte(`[{"foo":"bar"},null]`)))

	if err := root.AppendArray(root); err == nil {
		t.Errorf("AppendArray must returns error")
	}

	nodes, err := root.JSONPath(`$..foo`)
	if err != nil {
		t.Errorf("JSONPath returns error: %s", err)
	}
	if err := root.AppendArray(nodes...); err != nil {
		t.Errorf("AppendArray returns error: %s", err)
	}

	if value, err := Marshal(root); err != nil {
		t.Errorf("Marshal returns error: %v", err)
	} else if string(value) != `[{},null,"bar"]` {
		t.Errorf("Marshal returns wrong value: %s", string(value))
	}

	object := root.MustIndex(0)
	if value, err := Marshal(object); err != nil {
		t.Errorf("Marshal returns error: %v", err)
	} else if string(value) != `{}` {
		t.Errorf("Marshal returns wrong value: %s", string(value))
	}

	err = root.AppendArray(object)
	if err != nil {
		t.Errorf("AppendArray returns error: %s", err)
	}
	if value, err := Marshal(root); err != nil {
		t.Errorf("Marshal returns error: %v", err)
	} else if string(value) != `[null,"bar",{}]` {
		t.Errorf("Marshal returns wrong value: %s", string(value))
	}

	err = root.AppendArray(ArrayNode("", nil))
	if err != nil {
		t.Errorf("AppendArray returns error: %s", err)
	}
	if value, err := Marshal(root); err != nil {
		t.Errorf("Marshal returns error: %v", err)
	} else if string(value) != `[null,"bar",{},[]]` {
		t.Errorf("Marshal returns wrong value: %s", string(value))
	}

	err = root.MustIndex(3).AppendArray(
		root.MustIndex(0),
		root.MustIndex(1),
		root.MustIndex(2),
	)
	if err != nil {
		t.Errorf("AppendArray returns error: %s", err)
	}
	if value, err := Marshal(root); err != nil {
		t.Errorf("Marshal returns error: %v", err)
	} else if string(value) != `[[null,"bar",{}]]` {
		t.Errorf("Marshal returns wrong value: %s", string(value))
	}

	err = root.AppendArray(
		root.MustIndex(0).MustIndex(0),
		root.MustIndex(0).MustIndex(1),
		root.MustIndex(0).MustIndex(2),
	)
	if err != nil {
		t.Errorf("AppendArray returns error: %s", err)
	}
	if value, err := Marshal(root); err != nil {
		t.Errorf("Marshal returns error: %v", err)
	} else if string(value) != `[[],null,"bar",{}]` {
		t.Errorf("Marshal returns wrong value: %s", string(value))
	}
}
