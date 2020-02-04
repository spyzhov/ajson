package ajson

import (
	"strconv"
	"sync/atomic"
)

// IsDirty is the flag that shows, was node changed or not
func (n *Node) IsDirty() bool {
	return n.dirty
}

// SetNull update current node value with Null value
func (n *Node) SetNull() {
	n.update(Null, nil)
}

// SetNumeric update current node value with Numeric value
func (n *Node) SetNumeric(value float64) {
	n.update(Numeric, value)
}

// SetString update current node value with String value
func (n *Node) SetString(value string) {
	n.update(String, value)
}

// SetBool update current node value with Bool value
func (n *Node) SetBool(value bool) {
	n.update(Bool, value)
}

// SetArray update current node value with Array value
func (n *Node) SetArray(value []*Node) {
	n.update(Array, value)
}

// SetObject update current node value with Object value
func (n *Node) SetObject(value map[string]*Node) {
	n.update(Object, value)
}

// AppendArray append current Array node values with Node values
func (n *Node) AppendArray(value ...*Node) error {
	if !n.IsArray() {
		return errorType()
	}
	for _, val := range value {
		if err := n.appendArray(val); err != nil {
			return err
		}
	}
	n.mark()
	return nil
}

// AppendObject append current Object node value with key:value
// TODO
// func (n *Node) AppendObject(key string, value *Node) error {
// 	n.mark()
// 	if !n.IsObject() {
// 		return errorType()
// 	}
// 	if old, ok := n.children[key]; ok {
// 		old.parent = nil
// 	}
// 	value.parent = n
// 	value.key = &key
// 	value.index = nil
// 	n.children[key] = value
// 	return nil
// }

// TODO
// func (n *Node) DeleteNode(value *Node) error {
// }

// TODO
// func (n *Node) DeleteKey(key string) error {
// }

// TODO
// func (n *Node) DeleteIndex(index int) error {
// }

// update stored value, without validations
func (n *Node) update(_type NodeType, value interface{}) {
	n.mark()
	n.clear()

	atomic.StoreInt32((*int32)(&n._type), int32(_type))
	n.value = atomic.Value{}
	if value != nil {
		n.value.Store(value)
	}
}

// update stored value, without validations
func (n *Node) remove(value *Node) error {
	if !n.isContainer() {
		return errorType()
	}
	if value.parent != n {
		return errorRequest("wrong parent")
	}
	n.mark()
	if n.IsArray() {
		delete(n.children, strconv.Itoa(*value.index))
		n.dropindex(*value.index)
	} else {
		delete(n.children, *value.key)
	}
	return nil
}

// dropindex: internal method to reindexing current array value
func (n *Node) dropindex(index int) {
	for i := index + 1; i <= len(n.children); i++ {
		previous := i - 1
		if current, ok := n.children[strconv.Itoa(i)]; ok {
			current.index = &previous
			n.children[strconv.Itoa(previous)] = current
		}
		delete(n.children, strconv.Itoa(i))
	}
}

// appendArray append current Array node value with Node value
func (n *Node) appendArray(value *Node) error {
	if n.isParentNode(value) {
		return errorRequest("try to create infinite loop")
	}
	if value.parent != nil {
		if err := value.parent.remove(value); err != nil {
			return err
		}
	}
	value.parent = n
	size := len(n.children)
	value.key = nil
	value.index = &size
	n.children[strconv.Itoa(size)] = value
	return nil
}

// mark node as dirty, with all parents (up the tree)
func (n *Node) mark() {
	node := n
	for node != nil && !node.dirty {
		node.dirty = true
		node = node.parent
	}
}

// clear current value of node
func (n *Node) clear() {
	n.data = nil
	n.borders[1] = 0
	for key := range n.children {
		n.children[key].parent = nil
	}
	n.children = nil
}

// isParentNode check if current node is one of the parents
func (n *Node) isParentNode(node *Node) bool {
	for current := n; current != nil; current = current.parent {
		if current == node {
			return true
		}
	}
	return false
}
