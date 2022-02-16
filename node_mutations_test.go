package ajson

import (
	"fmt"
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
			node: NewNull(),
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
				NewNull(),
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
				"baz": NewNull(),
			}),
		},
		{
			name: "parsed Object",
			node: Must(Unmarshal([]byte(`{"foo": true, "bar": null, "baz": 123}`))),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.node.SetNull()
			if err != nil {
				t.Errorf("SetNull returns error: %s", err)
			} else if value, err := test.node.GetNull(); err != nil {
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
			node: NewNull(),
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
				NewNull(),
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
				"baz": NewNull(),
			}),
		},
		{
			name: "parsed Object",
			node: Must(Unmarshal([]byte(`{"foo": true, "bar": null, "baz": 123}`))),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.node.SetNumeric(expected)
			if err != nil {
				t.Errorf("SetNumeric returns error: %s", err)
			} else if value, err := test.node.GetNumeric(); err != nil {
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
			node: NewNull(),
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
				NewNull(),
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
				"baz": NewNull(),
			}),
		},
		{
			name: "parsed Object",
			node: Must(Unmarshal([]byte(`{"foo": true, "bar": null, "baz": 123}`))),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.node.SetString(expected)
			if err != nil {
				t.Errorf("SetString returns error: %s", err)
			} else if value, err := test.node.GetString(); err != nil {
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
			node: NewNull(),
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
				NewNull(),
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
				"baz": NewNull(),
			}),
		},
		{
			name: "parsed Object",
			node: Must(Unmarshal([]byte(`{"foo": true, "bar": null, "baz": 123}`))),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.node.SetBool(expected)
			if err != nil {
				t.Errorf("SetBool returns error: %s", err)
			} else if value, err := test.node.GetBool(); err != nil {
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
		withKey(NewNull(), "0"),
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
			node: NewNull(),
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
				NewNull(),
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
				"baz": NewNull(),
			}),
		},
		{
			name: "parsed Object",
			node: Must(Unmarshal([]byte(`{"foo": true, "bar": null, "baz": 123}`))),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.node.SetArray(expected)
			if err != nil {
				t.Errorf("SetArray returns error: %s", err)
			} else if value, err := test.node.GetArray(); err != nil {
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
		})
	}
}

func TestNode_SetObject(t *testing.T) {
	expected := map[string]*Node{
		"foo": NewNull(),
		"bar": BoolNode("bar", false),
	}
	tests := []struct {
		name string
		node *Node
	}{
		{
			name: "Null",
			node: NewNull(),
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
				NewNull(),
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
				"baz": NewNull(),
			}),
		},
		{
			name: "parsed Object",
			node: Must(Unmarshal([]byte(`{"foo": true, "bar": null, "baz": 123}`))),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.node.SetObject(expected)
			if err != nil {
				t.Errorf("SetArray returns error: %s", err)
			} else if value, err := test.node.GetObject(); err != nil {
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

	err = root.SetNull()
	if err != nil {
		t.Errorf("SetNull returns error: %s", err)
	}
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

	err = newRoot.SetNull()
	if err != nil {
		t.Errorf("SetNull returns error: %s", err)
	}
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
	if err := Must(Unmarshal([]byte(`[{"foo":"bar"}]`))).AppendArray(NewNull()); err != nil {
		t.Errorf("AppendArray should return error")
	}

	root := Must(Unmarshal([]byte(`[{"foo":"bar"}]`)))

	if err := root.AppendArray(NewNull()); err != nil {
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
	if err := root.MustIndex(0).AppendArray(NewNull()); err == nil {
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
	if err := Must(Unmarshal([]byte(`{"foo":"bar","baz":null}`))).AppendObject("biz", NewNull()); err != nil {
		t.Errorf("AppendArray should return error")
	}

	root := Must(Unmarshal([]byte(`{"foo":"bar"}`)))

	if err := root.AppendObject("biz", NewNull()); err != nil {
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
	if err := root.MustKey("fiz").AppendObject("fiz", NewNull()); err == nil {
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

	if err := root.DeleteNode(NewNull()); err == nil {
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

	node := NewNull()
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

func TestIssue22_SetArray_not_working(t *testing.T) {
	data := []byte(`{"key": [1, 2, 3]}`)
	node := NumericNode("", 4)
	expected := `{"key":[1,2,4]}`

	root := Must(Unmarshal(data))
	parent := root.MustKey("key")

	vals := parent.MustArray()
	vals[2] = node
	err := parent.SetArray(vals)
	if err != nil {
		t.Errorf("SetArray returns error: %s", err)
	}

	actual, err := Marshal(root)
	if err != nil {
		t.Errorf("error on Marshal(): %s", err)
	} else if string(actual) != expected {
		t.Errorf("actual != expected: \n'%s'\n'%s'", string(actual), expected)
	}
}

func TestNode_SetArray1(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		path     string
		value    []*Node
		wantErr  bool
		expected string
	}{
		{
			name:     "null -> []",
			json:     `null`,
			path:     `$`,
			value:    []*Node{},
			wantErr:  false,
			expected: `[]`,
		},
		{
			name:     "null -> [1,2,3]",
			json:     `null`,
			path:     `$`,
			value:    []*Node{NumericNode("", 1), NumericNode("", 2), NumericNode("", 3)},
			wantErr:  false,
			expected: `[1,2,3]`,
		},
		{
			name:     `{"key": null} -> {"key": [1,2,3]}`,
			json:     `{"key": null}`,
			path:     `$.key`,
			value:    []*Node{NumericNode("", 1), NumericNode("", 2), NumericNode("", 3)},
			wantErr:  false,
			expected: `{"key":[1,2,3]}`,
		},
		{
			name:     `{"key": [1,2,3]} -> {"key": [1,4,3]}`,
			json:     `{"key": [1,2,3]}`,
			path:     `$.key`,
			value:    []*Node{NumericNode("", 1), NumericNode("", 4), NumericNode("", 3)},
			wantErr:  false,
			expected: `{"key":[1,4,3]}`,
		},
		{
			name:     `{"key": [[1,2,3],2,3]} -> {"key": [[4,5,6],2,3]}`,
			json:     `{"key": [[1,2,3],2,3]}`,
			path:     `$.key[0]`,
			value:    []*Node{NumericNode("", 4), NumericNode("", 5), NumericNode("", 6)},
			wantErr:  false,
			expected: `{"key":[[4,5,6],2,3]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := Must(Unmarshal([]byte(tt.json)))
			nodes, err := root.JSONPath(tt.path)
			if err != nil {
				t.Errorf("JSONPath() error = %v", err)
			}
			if len(nodes) != 1 {
				t.Errorf("JSONPath() wrong response")
			}
			if err := nodes[0].SetArray(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("SetArray() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			actual, err := Marshal(root)
			if err != nil {
				t.Errorf("Marshal() error = %v", err)
			}
			if string(actual) != tt.expected {
				t.Errorf("actual != expected: \n'%s'\n'%s'", string(actual), tt.expected)
			}
		})
	}
}

func TestNode_SetObject1(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		path     string
		value    map[string]*Node
		wantErr  bool
		expected string
	}{
		{
			name:     "null -> {}",
			json:     `null`,
			path:     `$`,
			value:    map[string]*Node{},
			wantErr:  false,
			expected: `{}`,
		},
		{
			name:     `null -> {"foo": "bar"}`,
			json:     `null`,
			path:     `$`,
			value:    map[string]*Node{"foo": StringNode("foo", "bar")},
			wantErr:  false,
			expected: `{"foo":"bar"}`,
		},
		{
			name:     `{"key": null} -> {"key": {"foo": "bar"}}`,
			json:     `{"key": null}`,
			path:     `$.key`,
			value:    map[string]*Node{"foo": StringNode("foo", "bar")},
			wantErr:  false,
			expected: `{"key":{"foo":"bar"}}`,
		},
		{
			name:     `{"key": [1,2,3]} -> {"key": {"foo":"bar"}}`,
			json:     `{"key": [1,2,3]}`,
			path:     `$.key`,
			value:    map[string]*Node{"foo": StringNode("foo", "bar")},
			wantErr:  false,
			expected: `{"key":{"foo":"bar"}}`,
		},
		{
			name:     `{"key": [[1,2,3],2,3]} -> {"key": [{"foo":"bar"},2,3]}`,
			json:     `{"key": [[1,2,3],2,3]}`,
			path:     `$.key[0]`,
			value:    map[string]*Node{"foo": StringNode("foo", "bar")},
			wantErr:  false,
			expected: `{"key":[{"foo":"bar"},2,3]}`,
		},
		{
			name:     `{"key": {"baz": [null]}} -> {"key": {"baz": [{"foo":"bar"}]}}`,
			json:     `{"key": {"baz": [null]}}`,
			path:     `$.key.baz[0]`,
			value:    map[string]*Node{"foo": StringNode("foo", "bar")},
			wantErr:  false,
			expected: `{"key":{"baz":[{"foo":"bar"}]}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := Must(Unmarshal([]byte(tt.json)))
			nodes, err := root.JSONPath(tt.path)
			if err != nil {
				t.Errorf("JSONPath() error = %v", err)
			}
			if len(nodes) != 1 {
				t.Errorf("JSONPath() wrong response")
			}
			if err := nodes[0].SetObject(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("SetObject() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			actual, err := Marshal(root)
			if err != nil {
				t.Errorf("Marshal() error = %v", err)
			}
			if string(actual) != tt.expected {
				t.Errorf("actual != expected: \n'%s'\n'%s'", string(actual), tt.expected)
			}
		})
	}
}

func TestNode_update(t *testing.T) {
	type args struct {
		_type NodeType
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Null: success",
			args: args{
				_type: Null,
				value: nil,
			},
			wantErr: false,
		},
		{
			name: "Null: fail",
			args: args{
				_type: Null,
				value: "string",
			},
			wantErr: true,
		},
		{
			name: "String: success",
			args: args{
				_type: String,
				value: "",
			},
			wantErr: false,
		},
		{
			name: "String: fail",
			args: args{
				_type: String,
				value: nil,
			},
			wantErr: true,
		},
		{
			name: "Numeric: success",
			args: args{
				_type: Numeric,
				value: 1.1,
			},
			wantErr: false,
		},
		{
			name: "String: fail",
			args: args{
				_type: Numeric,
				value: nil,
			},
			wantErr: true,
		},
		{
			name: "Bool: success",
			args: args{
				_type: Bool,
				value: false,
			},
			wantErr: false,
		},
		{
			name: "Bool: fail",
			args: args{
				_type: Bool,
				value: nil,
			},
			wantErr: true,
		},
		{
			name: "Array: success",
			args: args{
				_type: Array,
				value: []*Node{},
			},
			wantErr: false,
		},
		{
			name: "Array: success nil",
			args: args{
				_type: Array,
				value: nil,
			},
			wantErr: false,
		},
		{
			name: "Array: fail",
			args: args{
				_type: Array,
				value: 1.1,
			},
			wantErr: true,
		},
		{
			name: "Object: success",
			args: args{
				_type: Object,
				value: map[string]*Node{},
			},
			wantErr: false,
		},
		{
			name: "Object: success nil",
			args: args{
				_type: Object,
				value: nil,
			},
			wantErr: false,
		},
		{
			name: "Object: fail",
			args: args{
				_type: Object,
				value: 1.1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewNull()
			if err := node.update(tt.args._type, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNode_update_fail(t *testing.T) {
	i0 := 0
	node := NewNull()
	parent := NewNull()
	broken := NewNull()
	broken.index = &i0
	broken.parent = parent

	if err := node.SetArray([]*Node{broken}); err == nil {
		t.Errorf("SetArray() error is nil")
	}
	if err := node.SetObject(map[string]*Node{"foo": broken}); err == nil {
		t.Errorf("SetObject() error is nil")
	}
}

func TestNode_Clone(t *testing.T) {
	node := NumericNode("", 1.1)
	null := NewNull()
	array := ArrayNode("", []*Node{node, null})
	object := ObjectNode("", map[string]*Node{"array": array})

	tests := []struct {
		name string
		node *Node
		json string
	}{
		{
			name: "null",
			node: null,
			json: "null",
		},
		{
			name: "node",
			node: node,
			json: "1.1",
		},
		{
			name: "array",
			node: array,
			json: "[1.1,null]",
		},
		{
			name: "object",
			node: object,
			json: `{"array":[1.1,null]}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			clone := test.node.Clone()
			if clone.parent != nil {
				t.Error("Clone().parent != nil")
			} else if clone.index != nil {
				t.Error("Clone().index != nil")
			} else if clone.key != nil {
				t.Error("Clone().key != nil")
			}

			if result, err := Marshal(clone); err != nil {
				t.Errorf("Marshal() error: %s", err)
			} else if string(result) != test.json {
				t.Errorf("Marshal() clone not match: \nExpected: %s\nActual: %s", test.json, result)
			} else if base, err := Marshal(test.node); err != nil {
				t.Errorf("Marshal() error: %s", err)
			} else if string(base) != test.json {
				t.Errorf("Marshal() base not match: \nExpected: %s\nActual: %s", test.json, base)
			}
		})
	}
}

func ExampleNode_Clone() {
	root := Must(Unmarshal(jsonPathTestData))
	nodes, _ := root.JSONPath("$..price")
	for i, node := range nodes {
		nodes[i] = node.Clone()
	}

	result, _ := Marshal(ArrayNode("", nodes))
	fmt.Printf("Array: %s\n", result)

	result, _ = Marshal(root)
	fmt.Printf("Basic: %s\n", result)

	// Output:
	// Array: [19.95,8.95,12.99,8.99,22.99]
	// Basic: { "store": {
	//     "book": [
	//       { "category": "reference",
	//         "author": "Nigel Rees",
	//         "title": "Sayings of the Century",
	//         "price": 8.95
	//       },
	//       { "category": "fiction",
	//         "author": "Evelyn Waugh",
	//         "title": "Sword of Honour",
	//         "price": 12.99
	//       },
	//       { "category": "fiction",
	//         "author": "Herman Melville",
	//         "title": "Moby Dick",
	//         "isbn": "0-553-21311-3",
	//         "price": 8.99
	//       },
	//       { "category": "fiction",
	//         "author": "J. R. R. Tolkien",
	//         "title": "The Lord of the Rings",
	//         "isbn": "0-395-19395-8",
	//         "price": 22.99
	//       }
	//     ],
	//     "bicycle": {
	//       "color": "red",
	//       "price": 19.95
	//     }
	//   }
	// }
	//
}

func TestNode_SetNode(t *testing.T) {
	iValue := `{"foo": [{"bar":"baz"}]}`
	idempotent := Must(Unmarshal([]byte(iValue)))
	child := StringNode("", "example")
	parent := ArrayNode("", []*Node{child})
	array := ArrayNode("", []*Node{})
	proxy := func(root *Node) *Node {
		return root
	}

	tests := []struct {
		name    string
		root    *Node
		getter  func(root *Node) *Node
		value   *Node
		result  string
		after   func(t *testing.T)
		wantErr bool
	}{
		{
			name:    "Null->Numeric(1)",
			root:    NewNull(),
			getter:  proxy,
			value:   NumericNode("", 1),
			result:  `1`,
			wantErr: false,
		},
		{
			name:    "Null->Object",
			root:    NewNull(),
			getter:  proxy,
			value:   Must(Unmarshal([]byte(`{"bar":"baz"}`))),
			result:  `{"bar":"baz"}`,
			wantErr: false,
		},
		{
			name:    "Null->Object(ref)",
			root:    NewNull(),
			getter:  proxy,
			value:   idempotent.MustObject()["foo"].MustArray()[0],
			result:  `{"bar":"baz"}`,
			wantErr: false,
		},
		{
			name: "[Numeric(1)]->[Object]",
			root: Must(Unmarshal([]byte(`[1]`))),
			getter: func(root *Node) *Node {
				return root.MustArray()[0]
			},
			value:   Must(Unmarshal([]byte(`{"bar":"baz"}`))),
			result:  `[{"bar":"baz"}]`,
			wantErr: false,
		},
		{
			name: "[Numeric(1)]->[Object(ref)]",
			root: Must(Unmarshal([]byte(`[1]`))),
			getter: func(root *Node) *Node {
				return root.MustIndex(0)
			},
			value:   idempotent.MustKey("foo").MustIndex(0),
			result:  `[{"bar":"baz"}]`,
			wantErr: false,
		},
		{
			name: "{foo:Null}->{foo:Object(ref)}",
			root: Must(Unmarshal([]byte(`{"foo":null}`))),
			getter: func(root *Node) *Node {
				return root.MustKey("foo")
			},
			value:   idempotent.MustKey("foo").MustIndex(0),
			result:  `{"foo":{"bar":"baz"}}`,
			wantErr: false,
		},
		{
			name: "parent[child]->parent[parent]",
			root: parent,
			getter: func(_ *Node) *Node {
				return child
			},
			value:   parent,
			result:  ``,
			wantErr: true,
		},
		{
			name:    "array[]->array[]",
			root:    array,
			getter:  proxy,
			value:   array,
			result:  `[]`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := tt.getter(tt.root)
			if err := node.SetNode(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("SetNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				return
			}
			if tt.root.String() != tt.result {
				t.Errorf("SetNode() value not match: \nExpected: %s\nActual: %s", tt.result, tt.root.String())
				return
			}
			if idempotent.String() != iValue {
				t.Errorf("SetNode() unexpected update value: \nUpdated: %s", idempotent.String())
				return
			}
			if tt.after != nil {
				tt.after(t)
			}
		})
	}
}

func TestNode_Set(t *testing.T) {
	node := func(data string) *Node {
		return Must(Unmarshal([]byte(data)))
	}
	tests := []struct {
		name    string
		node    *Node
		getter  func(root *Node) *Node
		value   interface{}
		result  string
		wantErr bool
	}{
		{
			name:    "Null->float64(123.456)",
			node:    node("null"),
			value:   float64(123.456),
			result:  "123.456",
			wantErr: false,
		},
		{
			name:    "Null->float32(123)",
			node:    node("null"),
			value:   float32(123),
			result:  "123",
			wantErr: false,
		},
		{
			name:    "Null->int(123)",
			node:    node("null"),
			value:   123,
			result:  "123",
			wantErr: false,
		},
		{
			name:    "Null->int8(123)",
			node:    node("null"),
			value:   int8(123),
			result:  "123",
			wantErr: false,
		},
		{
			name:    "Null->int16(123)",
			node:    node("null"),
			value:   int16(123),
			result:  "123",
			wantErr: false,
		},
		{
			name:    "Null->int32(123)",
			node:    node("null"),
			value:   int32(123),
			result:  "123",
			wantErr: false,
		},
		{
			name:    "Null->int64(123)",
			node:    node("null"),
			value:   int64(123),
			result:  "123",
			wantErr: false,
		},
		{
			name:    "Null->uint8(123)",
			node:    node("null"),
			value:   uint8(123),
			result:  "123",
			wantErr: false,
		},
		{
			name:    "Null->uint16(123)",
			node:    node("null"),
			value:   uint16(123),
			result:  "123",
			wantErr: false,
		},
		{
			name:    "Null->uint32(123)",
			node:    node("null"),
			value:   uint32(123),
			result:  "123",
			wantErr: false,
		},
		{
			name:    "Null->uint64(123)",
			node:    node("null"),
			value:   uint64(123),
			result:  "123",
			wantErr: false,
		},
		{
			name:    "Null->uint(123)",
			node:    node("null"),
			value:   uint(123),
			result:  "123",
			wantErr: false,
		},
		{
			name:    "Array[]->string",
			node:    node("[123]"),
			value:   "example value",
			result:  `"example value"`,
			wantErr: false,
		},
		{
			name:    "Object[]->bool",
			node:    node(`{"foo":["bar"]}`),
			value:   true,
			result:  `true`,
			wantErr: false,
		},
		{
			name: "Object[V]->bool",
			node: node(`{"foo":["bar"]}`),
			getter: func(root *Node) *Node {
				return root.MustKey("foo").MustIndex(0)
			},
			value:   true,
			result:  `{"foo":[true]}`,
			wantErr: false,
		},
		{
			name: "Object[V]->nil",
			node: node(`{"foo":["bar"]}`),
			getter: func(root *Node) *Node {
				return root.MustKey("foo").MustIndex(0)
			},
			value:   nil,
			result:  `{"foo":[null]}`,
			wantErr: false,
		},
		{
			name: "Array[V]->Array[Array[]]",
			node: node(`[null]`),
			getter: func(root *Node) *Node {
				return root.MustIndex(0)
			},
			value:   []*Node{node(`1`)},
			result:  `[[1]]`,
			wantErr: false,
		},
		{
			name: "Array[V]->Array[Object[]]",
			node: node(`[null]`),
			getter: func(root *Node) *Node {
				return root.MustIndex(0)
			},
			value:   map[string]*Node{},
			result:  `[{}]`,
			wantErr: false,
		},
		{
			name: "Array[V]->Array[Node]",
			node: node(`[null]`),
			getter: func(root *Node) *Node {
				return root.MustIndex(0)
			},
			value:   node(`{}`),
			result:  `[{}]`,
			wantErr: false,
		},
		{
			name:    "wrong_type",
			node:    node(`[null]`),
			value:   new(string),
			wantErr: true,
		},
		{
			name:    "nil",
			node:    nil,
			value:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := tt.node
			if tt.getter != nil {
				value = tt.getter(tt.node)
			}
			if err := value.Set(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if tt.node.String() != tt.result {
				t.Errorf("Set() value not match: \nExpected: %s\nActual: %s", tt.result, tt.node.String())
				return
			}
		})
	}
}
