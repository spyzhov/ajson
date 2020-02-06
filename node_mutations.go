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
		if err := n.appendNode(nil, val); err != nil {
			return err
		}
	}
	n.mark()
	return nil
}

// AppendObject append current Object node value with key:value
func (n *Node) AppendObject(key string, value *Node) error {
	if !n.IsObject() {
		return errorType()
	}
	err := n.appendNode(&key, value)
	if err != nil {
		return err
	}
	n.mark()
	return nil
}

// DeleteNode removes element child
func (n *Node) DeleteNode(value *Node) error {
	return n.remove(value)
}

// DeleteKey removes element from Object, by it's key
func (n *Node) DeleteKey(key string) error {
	node, err := n.GetKey(key)
	if err != nil {
		return err
	}
	return n.remove(node)
}

// PopKey removes element from Object, by it's key and return it
func (n *Node) PopKey(key string) (node *Node, err error) {
	node, err = n.GetKey(key)
	if err != nil {
		return
	}
	return node, n.remove(node)
}

// DeleteIndex removes element from Array, by it's index
func (n *Node) DeleteIndex(index int) error {
	node, err := n.GetIndex(index)
	if err != nil {
		return err
	}
	return n.remove(node)
}

// PopIndex removes element from Array, by it's index and return it
func (n *Node) PopIndex(index int) (node *Node, err error) {
	node, err = n.GetIndex(index)
	if err != nil {
		return
	}
	return node, n.remove(node)
}

// Delete removes element from parent. For root - do nothing.
func (n *Node) Delete() error {
	if n.parent == nil {
		return nil
	}
	return n.parent.remove(n)
}

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
	value.parent = nil
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

// appendNode append current Node node value with new Node value, by key or index
func (n *Node) appendNode(key *string, value *Node) error {
	if n.isParentNode(value) {
		return errorRequest("try to create infinite loop")
	}
	if value.parent != nil {
		if err := value.parent.remove(value); err != nil {
			return err
		}
	}
	value.parent = n
	value.key = key
	if key != nil {
		if old, ok := n.children[*key]; ok {
			if old != value {
				if err := n.remove(old); err != nil {
					return err
				}
			}
		}
		n.children[*key] = value
	} else {
		index := len(n.children)
		value.index = &index
		n.children[strconv.Itoa(index)] = value
	}
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
