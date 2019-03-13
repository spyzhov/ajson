package ajson

import (
	"strconv"
	"sync/atomic"
)

// Main struct, presents any json node
type Node struct {
	parent   *Node
	children []*Node
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

func newNode(parent *Node, dec *buffer, _type NodeType, key *string) *Node {
	node := &Node{
		parent:  parent,
		data:    &dec.data,
		borders: [2]int{dec.index, 0},
		_type:   _type,
		key:     key,
	}
	if parent != nil && parent.IsArray() {
		size := len(parent.children)
		node.index = &size
		parent.children = append(parent.children, node)
	}
	return node
}

func (n *Node) Source() []byte {
	return (*n.data)[n.borders[0]:n.borders[1]]
}

func (n *Node) String() string {
	return string(n.Source())
}

func (n *Node) Type() NodeType {
	return n._type
}

func (n *Node) Key() string {
	return *n.key
}

func (n *Node) Index() int {
	return *n.index
}

func (n *Node) IsArray() bool {
	return n._type == Array
}

func (n *Node) IsObject() bool {
	return n._type == Object
}

func (n *Node) IsNull() bool {
	return n._type == Null
}

func (n *Node) IsNumeric() bool {
	return n._type == Numeric
}

func (n *Node) IsString() bool {
	return n._type == String
}

func (n *Node) IsBool() bool {
	return n._type == Bool
}

func (n *Node) Value() (value interface{}, err error) {
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
			n.set(value)
		case String:
			size := len(n.Source())
			value = string(n.Source()[1 : size-1])
			n.set(value)
		case Bool:
			b := n.Source()[0]
			value = b == 't' || b == 'T'
			n.set(value)
		case Array:
			value = n.children
			n.set(value)
		case Object:
			result := make(map[string]*Node)
			for _, child := range n.children {
				result[child.Key()] = child
			}
			value = result
			n.set(value)
		}
	}
	return
}

func (n *Node) GetNumeric() (value float64, err error) {
	iValue, err := n.Value()
	if err != nil {
		return 0, err
	}
	value = iValue.(float64)
	return
}

func (n *Node) GetString() (value string, err error) {
	iValue, err := n.Value()
	if err != nil {
		return "", err
	}
	value = iValue.(string)
	return
}

func (n *Node) GetBool() (value bool, err error) {
	iValue, err := n.Value()
	if err != nil {
		return false, err
	}
	value = iValue.(bool)
	return
}

func (n *Node) GetArray() (value []*Node, err error) {
	iValue, err := n.Value()
	if err != nil {
		return nil, err
	}
	value = iValue.([]*Node)
	return
}

func (n *Node) GetObject() (value map[string]*Node, err error) {
	iValue, err := n.Value()
	if err != nil {
		return nil, err
	}
	value = iValue.(map[string]*Node)
	return
}

func (n *Node) set(value interface{}) {
	n.value.Store(value)
}

func (n *Node) ready() bool {
	return n.borders[1] != 0
}

func (n *Node) isContainer() bool {
	return n._type == Array || n._type == Object
}
