package ajson

import (
	"fmt"
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

func newNode(parent *Node, dec *buffer, _type NodeType, key *string, index *int) *Node {
	return &Node{
		parent:  parent,
		data:    &dec.data,
		borders: [2]int{dec.index, 0},
		_type:   _type,
		key:     key,
		index:   index,
	}
}

func (n *Node) Source() []byte {
	return (*n.data)[n.borders[0]:n.borders[1]]
}

func (n *Node) Type() NodeType {
	return n._type
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
			return nil, fmt.Errorf("not implemented")
		case Object:
			return nil, fmt.Errorf("not implemented")
		}
	}
	return
}

func (n *Node) set(value interface{}) {
	n.value.Store(value)
}
