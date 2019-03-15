package ajson

import (
	"strconv"
	"sync/atomic"
)

type Node interface {
	Source() []byte
	Type() NodeType
	Key() string
	Index() int
	Size() int
	Keys() []string

	String() string

	IsArray() bool
	IsObject() bool
	IsNull() bool
	IsNumeric() bool
	IsString() bool
	IsBool() bool

	Value() (interface{}, error)

	GetNull() (interface{}, error)
	GetNumeric() (float64, error)
	GetString() (string, error)
	GetBool() (bool, error)
	GetArray() ([]Node, error)
	GetObject() (map[string]Node, error)

	MustNull() interface{}
	MustNumeric() float64
	MustString() string
	MustBool() bool
	MustArray() []Node
	MustObject() map[string]Node
}

// Main struct, presents any json node
type node struct {
	parent   *node
	children []*node
	key      *string
	index    *int
	_type    NodeType
	data     *[]byte
	borders  [2]int
	value    atomic.Value
}

type NodeType int

const (
	Null NodeType = iota
	Numeric
	String
	Bool
	Array
	Object
)

func newNode(parent *node, buf *buffer, _type NodeType, key **string) (current *node, err error) {
	current = &node{
		parent:  parent,
		data:    &buf.data,
		borders: [2]int{buf.index, 0},
		_type:   _type,
		key:     *key,
	}
	if parent != nil {
		if parent.IsArray() {
			size := len(parent.children)
			current.index = &size
			parent.children = append(parent.children, current)
		} else if parent.IsObject() {
			parent.children = append(parent.children, current)
			if *key == nil {
				err = errorSymbol(buf)
			} else {
				*key = nil
			}
		} else {
			err = errorSymbol(buf)
		}
	}
	return
}

func (n *node) Source() []byte {
	return (*n.data)[n.borders[0]:n.borders[1]]
}

func (n *node) String() string {
	return string(n.Source())
}

func (n *node) Type() NodeType {
	return n._type
}

func (n *node) Key() string {
	return *n.key
}

func (n *node) Index() int {
	return *n.index
}

func (n *node) Size() int {
	return len(n.children)
}

func (n *node) Keys() (result []string) {
	result = make([]string, 0, len(n.children))
	for _, child := range n.children {
		if child.key != nil {
			result = append(result, *child.key)
		}
	}
	return
}

func (n *node) IsArray() bool {
	return n._type == Array
}

func (n *node) IsObject() bool {
	return n._type == Object
}

func (n *node) IsNull() bool {
	return n._type == Null
}

func (n *node) IsNumeric() bool {
	return n._type == Numeric
}

func (n *node) IsString() bool {
	return n._type == String
}

func (n *node) IsBool() bool {
	return n._type == Bool
}

func (n *node) Value() (value interface{}, err error) {
	value = n.value.Load()
	if value == nil {
		switch n._type {
		case Null:
			return nil, nil
		case Numeric:
			value, err = strconv.ParseFloat(string(n.Source()), 64)
			if err != nil {
				return
			}
			n.value.Store(value)
		case String:
			size := len(n.Source())
			value = string(n.Source()[1 : size-1])
			n.value.Store(value)
		case Bool:
			b := n.Source()[0]
			value = b == 't' || b == 'T'
			n.value.Store(value)
		case Array:
			children := make([]Node, 0, len(n.children))
			for _, child := range n.children {
				children = append(children, child)
			}
			value = children
			n.value.Store(value)
		case Object:
			result := make(map[string]Node)
			for _, child := range n.children {
				result[child.Key()] = child
			}
			value = result
			n.value.Store(value)
		}
	}
	return
}

func (n *node) GetNull() (value interface{}, err error) {
	if n._type != Null {
		return value, errorType()
	}
	return
}

func (n *node) GetNumeric() (value float64, err error) {
	if n._type != Numeric {
		return value, errorType()
	}
	iValue, err := n.Value()
	if err != nil {
		return 0, err
	}
	value = iValue.(float64)
	return
}

func (n *node) GetString() (value string, err error) {
	if n._type != String {
		return value, errorType()
	}
	iValue, err := n.Value()
	if err != nil {
		return "", err
	}
	value = iValue.(string)
	return
}

func (n *node) GetBool() (value bool, err error) {
	if n._type != Bool {
		return value, errorType()
	}
	iValue, err := n.Value()
	if err != nil {
		return false, err
	}
	value = iValue.(bool)
	return
}

func (n *node) GetArray() (value []Node, err error) {
	if n._type != Array {
		return value, errorType()
	}
	iValue, err := n.Value()
	if err != nil {
		return nil, err
	}
	value = iValue.([]Node)
	return
}

func (n *node) GetObject() (value map[string]Node, err error) {
	if n._type != Object {
		return value, errorType()
	}
	iValue, err := n.Value()
	if err != nil {
		return nil, err
	}
	value = iValue.(map[string]Node)
	return
}

func (n *node) MustNull() (value interface{}) {
	value, err := n.GetNull()
	if err != nil {
		panic(err)
	}
	return
}

func (n *node) MustNumeric() (value float64) {
	value, err := n.GetNumeric()
	if err != nil {
		panic(err)
	}
	return
}

func (n *node) MustString() (value string) {
	value, err := n.GetString()
	if err != nil {
		panic(err)
	}
	return
}

func (n *node) MustBool() (value bool) {
	value, err := n.GetBool()
	if err != nil {
		panic(err)
	}
	return
}

func (n *node) MustArray() (value []Node) {
	value, err := n.GetArray()
	if err != nil {
		panic(err)
	}
	return
}

func (n *node) MustObject() (value map[string]Node) {
	value, err := n.GetObject()
	if err != nil {
		panic(err)
	}
	return
}

func (n *node) ready() bool {
	return n.borders[1] != 0
}

func (n *node) isContainer() bool {
	return n._type == Array || n._type == Object
}
