package ajson

// Main struct, presents any json node
type Node struct {
	parent   *Node
	children []*Node
	key      []byte
	index    *int
	_type    NodeType
	data     *[]byte
	borders  [2]int
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

func newNode(parent *Node, dec *buffer, _type NodeType, key []byte, index *int) *Node {
	return &Node{
		parent:  parent,
		data:    &dec.data,
		borders: [2]int{dec.index, 0},
		_type:   _type,
		key:     key,
		index:   index,
	}
}

func (n *Node) Value() []byte {
	return (*n.data)[n.borders[0]:n.borders[1]]
}

func (n *Node) Type() NodeType {
	return n._type
}
