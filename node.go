package ajson

import (
	"sort"
	"strconv"
	"sync/atomic"
)

// Node is a main struct, presents any type of JSON node.
// Available types are:
//
//	const (
//		Null NodeType = iota
//		Numeric
//		String
//		Bool
//		Array
//		Object
//	)
//
// Every type has its own methods to be called.
// Every Node contains link to a byte data, parent and children, also calculated type of value, atomic value and internal information.
type Node struct {
	parent   *Node
	children map[string]*Node
	key      *string
	index    *int
	_type    NodeType
	data     *[]byte
	borders  [2]int
	value    atomic.Value
}

// NodeType is a kind of reflection of JSON type to a type of golang
type NodeType int

// Reflections:
//
//	Null    = nil.(interface{})
//	Numeric = float64
//	String  = string
//	Bool    = bool
//	Array   = []*Node
//	Object  = map[string]*Node
//
const (
	//Null is reflection of nil.(interface{})
	Null NodeType = iota
	//Numeric is reflection of float64
	Numeric
	//String is reflection of string
	String
	//Bool is reflection of bool
	Bool
	//Array is reflection of []*Node
	Array
	//Object is reflection of map[string]*Node
	Object
)

func newNode(parent *Node, buf *buffer, _type NodeType, key **string) (current *Node, err error) {
	current = &Node{
		parent:  parent,
		data:    &buf.data,
		borders: [2]int{buf.index, 0},
		_type:   _type,
		key:     *key,
	}
	if _type == Object || _type == Array {
		current.children = make(map[string]*Node)
	}
	if parent != nil {
		if parent.IsArray() {
			size := len(parent.children)
			current.index = &size
			parent.children[strconv.Itoa(size)] = current
		} else if parent.IsObject() {
			if *key == nil {
				err = errorSymbol(buf)
			} else {
				parent.children[**key] = current
				*key = nil
			}
		} else {
			err = errorSymbol(buf)
		}
	}
	return
}

//Parent returns link to the parent of current node, nil for root
func (n *Node) Parent() *Node {
	return n.parent
}

//Source returns slice of bytes, which was identified to be current node
func (n *Node) Source() []byte {
	return (*n.data)[n.borders[0]:n.borders[1]]
}

//String is implementation of Stringer interface, returns string based on source part
func (n *Node) String() string {
	return string(n.Source())
}

//Type will return type of current node
func (n *Node) Type() NodeType {
	return n._type
}

//Key will return key of current node, please check, that parent of this node has an Object type
func (n *Node) Key() string {
	return *n.key
}

//Index will return index of current node, please check, that parent of this node has an Array type
func (n *Node) Index() int {
	return *n.index
}

//Size will return count of children of current node, please check, that parent of this node has an Array type
func (n *Node) Size() int {
	return len(n.children)
}

//Keys will return count all keys of children of current node, please check, that parent of this node has an Object type
func (n *Node) Keys() (result []string) {
	result = make([]string, 0, len(n.children))
	for key := range n.children {
		result = append(result, key)
	}
	return
}

//IsArray returns true if current node is Array
func (n *Node) IsArray() bool {
	return n._type == Array
}

//IsObject returns true if current node is Object
func (n *Node) IsObject() bool {
	return n._type == Object
}

//IsNull returns true if current node is Null
func (n *Node) IsNull() bool {
	return n._type == Null
}

//IsNumeric returns true if current node is Numeric
func (n *Node) IsNumeric() bool {
	return n._type == Numeric
}

//IsString returns true if current node is String
func (n *Node) IsString() bool {
	return n._type == String
}

//IsBool returns true if current node is Bool
func (n *Node) IsBool() bool {
	return n._type == Bool
}

//Value is calculating and returns a value of current node.
//
// It returns nil, if current node type is Null.
//
// It returns float64, if current node type is Numeric.
//
// It returns string, if current node type is String.
//
// It returns bool, if current node type is Bool.
//
// It returns []*Node, if current node type is Array.
//
// It returns map[string]*Node, if current node type is Object.
//
// BUT! Current method doesn't calculate underlying nodes (use method Node.Unpack for that).
//
// Value will be calculated only once and saved into atomic.Value.
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
			children := make([]*Node, len(n.children))
			for _, child := range n.children {
				children[*child.index] = child
			}
			value = children
			n.value.Store(value)
		case Object:
			result := make(map[string]*Node)
			for key, child := range n.children {
				result[key] = child
			}
			value = result
			n.value.Store(value)
		}
	}
	return
}

//GetNull returns nil, if current type is Null, else: WrongType error
func (n *Node) GetNull() (value interface{}, err error) {
	if n._type != Null {
		return value, errorType()
	}
	return
}

//GetNumeric returns float64, if current type is Numeric, else: WrongType error
func (n *Node) GetNumeric() (value float64, err error) {
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

//GetString returns string, if current type is String, else: WrongType error
func (n *Node) GetString() (value string, err error) {
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

//GetBool returns bool, if current type is Bool, else: WrongType error
func (n *Node) GetBool() (value bool, err error) {
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

//GetArray returns []*Node, if current type is Array, else: WrongType error
func (n *Node) GetArray() (value []*Node, err error) {
	if n._type != Array {
		return value, errorType()
	}
	iValue, err := n.Value()
	if err != nil {
		return nil, err
	}
	value = iValue.([]*Node)
	return
}

//GetObject returns map[string]*Node, if current type is Object, else: WrongType error
func (n *Node) GetObject() (value map[string]*Node, err error) {
	if n._type != Object {
		return value, errorType()
	}
	iValue, err := n.Value()
	if err != nil {
		return nil, err
	}
	value = iValue.(map[string]*Node)
	return
}

//MustNull returns nil, if current type is Null, else: panic if error happened
func (n *Node) MustNull() (value interface{}) {
	value, err := n.GetNull()
	if err != nil {
		panic(err)
	}
	return
}

//MustNumeric returns float64, if current type is Numeric, else: panic if error happened
func (n *Node) MustNumeric() (value float64) {
	value, err := n.GetNumeric()
	if err != nil {
		panic(err)
	}
	return
}

//MustString returns string, if current type is String, else: panic if error happened
func (n *Node) MustString() (value string) {
	value, err := n.GetString()
	if err != nil {
		panic(err)
	}
	return
}

//MustBool returns bool, if current type is Bool, else: panic if error happened
func (n *Node) MustBool() (value bool) {
	value, err := n.GetBool()
	if err != nil {
		panic(err)
	}
	return
}

//MustArray returns []*Node, if current type is Array, else: panic if error happened
func (n *Node) MustArray() (value []*Node) {
	value, err := n.GetArray()
	if err != nil {
		panic(err)
	}
	return
}

//MustObject returns map[string]*Node, if current type is Object, else: panic if error happened
func (n *Node) MustObject() (value map[string]*Node) {
	value, err := n.GetObject()
	if err != nil {
		panic(err)
	}
	return
}

//Unpack will produce current node to it's interface, recursively with all underlying nodes (in contrast to Node.Value).
func (n *Node) Unpack() (value interface{}, err error) {
	switch n._type {
	case Null:
		return nil, nil
	case Numeric:
		value, err = strconv.ParseFloat(string(n.Source()), 64)
		if err != nil {
			return
		}
	case String:
		size := len(n.Source())
		value = string(n.Source()[1 : size-1])
	case Bool:
		b := n.Source()[0]
		value = b == 't' || b == 'T'
	case Array:
		children := make([]interface{}, len(n.children))
		for _, child := range n.children {
			val, err := child.Unpack()
			if err != nil {
				return nil, err
			}
			children[*child.index] = val
		}
		value = children
	case Object:
		result := make(map[string]interface{})
		for key, child := range n.children {
			result[key], err = child.Unpack()
			if err != nil {
				return nil, err
			}
		}
		value = result
	}
	return
}

//GetIndex will return child node of current array node. If current node is not Array, or index is unavailable, will return error
func (n *Node) GetIndex(index int) (*Node, error) {
	if n._type != Array {
		return nil, errorType()
	}
	child, ok := n.children[strconv.Itoa(index)]
	if !ok {
		return nil, errorRequest("out of index %d", index)
	}
	return child, nil
}

//MustIndex will return child node of current array node. If current node is not Array, or index is unavailable, raise a panic
func (n *Node) MustIndex(index int) (value *Node) {
	value, err := n.GetIndex(index)
	if err != nil {
		panic(err)
	}
	return
}

//GetKey will return child node of current object node. If current node is not Object, or key is unavailable, will return error
func (n *Node) GetKey(key string) (*Node, error) {
	if n._type != Object {
		return nil, errorType()
	}
	value, ok := n.children[key]
	if !ok {
		return nil, errorRequest("wrong key '%s'", key)
	}
	return value, nil
}

//MustKey will return child node of current object node. If current node is not Object, or key is unavailable, raise a panic
func (n *Node) MustKey(key string) (value *Node) {
	value, err := n.GetKey(key)
	if err != nil {
		panic(err)
	}
	return
}

//HasKey will return boolean value, if current object node has custom key
func (n *Node) HasKey(key string) bool {
	_, ok := n.children[key]
	return ok
}

//Empty method check if current container node has no children
func (n *Node) Empty() bool {
	return len(n.children) == 0
}

// Path returns full JsonPath of current Node
func (n *Node) Path() string {
	if n.parent == nil {
		return "$"
	}
	if n.parent.IsObject() {
		return n.parent.Path() + "['" + n.Key() + "']"
	}
	return n.parent.Path() + "[" + strconv.Itoa(n.Index()) + "]"
}

func (n *Node) ready() bool {
	return n.borders[1] != 0
}

func (n *Node) isContainer() bool {
	return n._type == Array || n._type == Object
}

// Return sorted by keys/index slice of children
func (n *Node) inheritors() (result []*Node) {
	size := len(n.children)
	if n.IsObject() {
		result = make([]*Node, size)
		keys := n.Keys()
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})
		for i, key := range keys {
			result[i] = n.children[key]
		}
	} else if n.IsArray() {
		result = make([]*Node, size)
		for _, element := range n.children {
			result[*element.index] = element
		}
	}
	return
}
