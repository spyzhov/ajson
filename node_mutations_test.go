package ajson

import (
	"reflect"
	"testing"
)

func testEqObject(str []byte, variants []string) bool {
	for _, val := range variants {
		if string(str) == val {
			return true
		}
	}
	return false
}

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
	if err := root.MustIndex(0).AppendArray(NullNode("")); err == nil {
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

func TestNode_AppendObject(t *testing.T) {
	if err := Must(Unmarshal([]byte(`{"foo":"bar","baz":null}`))).AppendObject("biz", NullNode("")); err != nil {
		t.Errorf("AppendArray should return error")
	}

	root := Must(Unmarshal([]byte(`{"foo":"bar"}`)))

	if err := root.AppendObject("biz", NullNode("")); err != nil {
		t.Errorf("AppendArray returns error: %v", err)
	}
	if value, err := Marshal(root); err != nil {
		t.Errorf("Marshal returns error: %v", err)
	} else if !testEqObject(value, []string{`{"foo":"bar","biz":null}`, `{"biz":null,"foo":"bar"}`}) {
		t.Errorf("Marshal returns wrong value: %s", string(value))
	}

	if err := root.AppendObject("foo", NumericNode("", 1)); err != nil {
		t.Errorf("AppendArray returns error: %v", err)
	}
	if value, err := Marshal(root); err != nil {
		t.Errorf("Marshal returns error: %v", err)
	} else if !testEqObject(value, []string{`{"foo":1,"biz":null}`, `{"biz":null,"foo":1}`}) {
		t.Errorf("Marshal returns wrong value: %s", string(value))
	}
}

func TestNode_AppendObject_self(t *testing.T) {
	root := Must(Unmarshal([]byte(`{"foo":{"bar":"baz"},"fiz":null}`)))

	if err := root.AppendObject("foo", root); err == nil {
		t.Errorf("AppendArray must returns error")
	}
	if err := root.MustKey("fiz").AppendObject("fiz", NullNode("")); err == nil {
		t.Errorf("AppendArray must returns error: not object")
	}

	if err := root.MustKey("foo").AppendObject("bar", root); err == nil {
		t.Errorf("AppendArray must returns error: self")
	}

	nodes, err := root.JSONPath(`$..bar`)
	if err != nil {
		t.Errorf("JSONPath returns error: %s", err)
	}
	if err := root.AppendObject("bar", nodes[0]); err != nil {
		t.Errorf("AppendArray returns error: %s", err)
	}

	if value, err := Marshal(root); err != nil {
		t.Errorf("Marshal returns error: %v", err)
	} else if !testEqObject(value, []string{
		`{"bar":"baz","foo":{},"fiz":null}`,
		`{"bar":"baz","fiz":null,"foo":{}}`,
		`{"fiz":null,"bar":"baz","foo":{}}`,
		`{"fiz":null,"foo":{},"bar":"baz"}`,
		`{"foo":{},"bar":"baz","fiz":null}`,
		`{"foo":{},"fiz":null,"bar":"baz"}`,
	}) {
		t.Errorf("Marshal returns wrong value: %s", string(value))
	}

	object := root.MustKey("foo")
	if value, err := Marshal(object); err != nil {
		t.Errorf("Marshal returns error: %v", err)
	} else if string(value) != `{}` {
		t.Errorf("Marshal returns wrong value: %s", string(value))
	}

	err = root.AppendObject("bar", object)
	if err != nil {
		t.Errorf("AppendArray returns error: %s", err)
	}
	if value, err := Marshal(root); err != nil {
		t.Errorf("Marshal returns error: %v", err)
	} else if !testEqObject(value, []string{
		`{"bar":{},"fiz":null}`,
		`{"fiz":null,"bar":{}}`,
	}) {
		t.Errorf("Marshal returns wrong value: %s", string(value))
	}
}

func TestNode_Delete(t *testing.T) {
	root := Must(Unmarshal([]byte(`{"foo":"bar"}`)))
	if err := root.Delete(); err != nil {
		t.Errorf("root.Delete returns error: %v", err)
	}
	if value, err := Marshal(root); err != nil {
		t.Errorf("Marshal returns error: %v", err)
	} else if string(value) != `{"foo":"bar"}` {
		t.Errorf("Marshal returns wrong value: %s", string(value))
	}

	foo := root.MustKey("foo")
	if err := foo.Delete(); err != nil {
		t.Errorf("foo.Delete returns error: %v", err)
	}
	if value, err := Marshal(root); err != nil {
		t.Errorf("Marshal returns error: %v", err)
	} else if string(value) != `{}` {
		t.Errorf("Marshal returns wrong value: %s", string(value))
	}
	if value, err := Marshal(foo); err != nil {
		t.Errorf("Marshal returns error: %v", err)
	} else if string(value) != `"bar"` {
		t.Errorf("Marshal returns wrong value: %s", string(value))
	}
	if foo.Parent() != nil {
		t.Errorf("Delete didn't remove parent")
	}
}

func TestNode_DeleteIndex(t *testing.T) {
	tests := []struct {
		json     string
		expected string
		index    int
		fail     bool
	}{
		{`null`, ``, 0, true},
		{`1`, ``, 0, true},
		{`{}`, ``, 0, true},
		{`{"foo":"bar"}`, ``, 0, true},
		{`true`, ``, 0, true},
		{`[]`, ``, 0, true},
		{`[]`, ``, -1, true},
		{`[1]`, `[]`, 0, false},
		{`[{}]`, `[]`, 0, false},
		{`[{}]`, `[]`, -1, false},
		{`[{},[],1]`, `[{},[]]`, -1, false},
		{`[{},[],1]`, `[{},1]`, 1, false},
		{`[{},[],1]`, ``, 10, true},
		{`[{},[],1]`, ``, -10, true},
	}
	for _, test := range tests {
		t.Run(test.json, func(t *testing.T) {
			root := Must(Unmarshal([]byte(test.json)))
			err := root.DeleteIndex(test.index)
			if test.fail {
				if err == nil {
					t.Errorf("Expected error")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				result, err := Marshal(root)
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				} else if string(result) != test.expected {
					t.Errorf("Unexpected result: %s", result)
				}
			}
		})
	}
}

func TestNode_PopIndex(t *testing.T) {
	tests := []struct {
		json     string
		expected string
		index    int
		fail     bool
	}{
		{`null`, ``, 0, true},
		{`1`, ``, 0, true},
		{`{}`, ``, 0, true},
		{`{"foo":"bar"}`, ``, 0, true},
		{`true`, ``, 0, true},
		{`[]`, ``, 0, true},
		{`[]`, ``, -1, true},
		{`[1]`, `[]`, 0, false},
		{`[{}]`, `[]`, 0, false},
		{`[{}]`, `[]`, -1, false},
		{`[{},[],1]`, `[{},[]]`, -1, false},
		{`[{},[],1]`, `[{},1]`, 1, false},
		{`[{},[],1]`, ``, 10, true},
		{`[{},[],1]`, ``, -10, true},
	}
	for _, test := range tests {
		t.Run(test.json, func(t *testing.T) {
			if test.fail {
				root := Must(Unmarshal([]byte(test.json)))
				_, err := root.PopIndex(test.index)
				if err == nil {
					t.Errorf("Expected error")
				}
			} else {
				root := Must(Unmarshal([]byte(test.json)))
				expected := root.MustIndex(test.index)
				node, err := root.PopIndex(test.index)
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if node == nil {
					t.Errorf("Unexpected node")
				}
				if node.Parent() != nil {
					t.Errorf("node.Parent is not nil")
				}
				if node != expected {
					t.Errorf("node is not expected")
				}
				result, err := Marshal(root)
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				} else if string(result) != test.expected {
					t.Errorf("Unexpected result: %s", result)
				}
			}
		})
	}
}

func TestNode_DeleteKey(t *testing.T) {
	tests := []struct {
		json     string
		expected string
		key      string
		fail     bool
	}{
		{`null`, ``, "", true},
		{`1`, ``, "", true},
		{`[]`, ``, "", true},
		{`[1,2,3]`, ``, "", true},
		{`true`, ``, "", true},
		{`{}`, ``, "", true},
		{`{}`, ``, "foo", true},
		{`{"foo":"bar"}`, ``, "bar", true},
		{`{"foo":"bar"}`, `{}`, "foo", false},
		{`{"foo":"bar","baz":1}`, `{"baz":1}`, "foo", false},
		{`{"foo":"bar","baz":1}`, `{"foo":"bar"}`, "baz", false},
		{`{"foo":"bar","baz":1}`, ``, "fiz", true},
	}
	for _, test := range tests {
		t.Run(test.json, func(t *testing.T) {
			root := Must(Unmarshal([]byte(test.json)))
			err := root.DeleteKey(test.key)
			if test.fail {
				if err == nil {
					t.Errorf("Expected error")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				result, err := Marshal(root)
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				} else if string(result) != test.expected {
					t.Errorf("Unexpected result: %s", result)
				}
			}
		})
	}
}

func TestNode_PopKey(t *testing.T) {
	tests := []struct {
		json     string
		expected string
		key      string
		fail     bool
	}{
		{`null`, ``, "", true},
		{`1`, ``, "", true},
		{`[]`, ``, "", true},
		{`[1,2,3]`, ``, "", true},
		{`true`, ``, "", true},
		{`{}`, ``, "", true},
		{`{}`, ``, "foo", true},
		{`{"foo":"bar"}`, ``, "bar", true},
		{`{"foo":"bar"}`, `{}`, "foo", false},
		{`{"foo":"bar","baz":1}`, `{"baz":1}`, "foo", false},
		{`{"foo":"bar","baz":1}`, `{"foo":"bar"}`, "baz", false},
		{`{"foo":"bar","baz":1}`, ``, "fiz", true},
	}
	for _, test := range tests {
		t.Run(test.json, func(t *testing.T) {
			root := Must(Unmarshal([]byte(test.json)))
			if test.fail {
				_, err := root.PopKey(test.key)
				if err == nil {
					t.Errorf("Expected error")
				}
			} else {
				expected := root.MustKey(test.key)
				node, err := root.PopKey(test.key)
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if node == nil {
					t.Errorf("Unexpected node")
				}
				if node.Parent() != nil {
					t.Errorf("node.Parent is not nil")
				}
				if node != expected {
					t.Errorf("node is not expected")
				}
				result, err := Marshal(root)
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				} else if string(result) != test.expected {
					t.Errorf("Unexpected result: %s", result)
				}
			}
		})
	}
}

func TestNode_DeleteNode(t *testing.T) {
	initial := `{"foo":{"bar":["baz",1,null]},"biz":"zip"}`
	root := Must(Unmarshal([]byte(initial)))

	if err := root.DeleteNode(NullNode("")); err == nil {
		t.Errorf("Expected error")
	}
	if err := root.DeleteNode(StringNode("biz", "zip")); err == nil {
		t.Errorf("Expected error")
	}
	if err := root.MustKey("biz").DeleteNode(root.MustKey("biz")); err == nil {
		t.Errorf("Expected error")
	}
	if err := root.MustKey("foo").DeleteNode(root.MustKey("foo")); err == nil {
		t.Errorf("Expected error")
	}

	node := NullNode("")
	if err := root.AppendObject("key", node); err != nil {
		t.Errorf("UnExpected error: %v", err)
	}
	if err := root.DeleteNode(node); err != nil {
		t.Errorf("UnExpected error: %v", err)
	}
	if value, err := Marshal(root); err != nil {
		t.Errorf("Marshal returns error: %v", err)
	} else if !testEqObject(value, []string{
		`{"foo":{"bar":["baz",1,null]},"biz":"zip"}`,
		`{"biz":"zip","foo":{"bar":["baz",1,null]}}`,
	}) {
		t.Errorf("Marshal returns wrong value: %s", string(value))
	}
}
