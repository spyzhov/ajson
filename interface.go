package ajson

//import (
//	"fmt"
//)
//
//var (
//	_ DefaultNode = new(Node)
//)
//
//type DefaultNode interface {
//	fmt.Stringer
//	// Type will return type of current node
//	Type() NodeType
//	Parent() *Node
//	Interface() (value interface{}, err error)
//	Set(value interface{}) error
//}
//
//type NullNode interface {
//	DefaultNode
//	GetNull() (interface{}, error)
//	MustNull() interface{}
//}
//
//type NumericNode interface {
//	DefaultNode
//
//}
//
//type StringNode interface {
//	DefaultNode
//
//}
//
//type BoolNode interface {
//	DefaultNode
//
//}
//
//type ArrayNode interface {
//	DefaultNode
//
//}
//
//type ObjectNode interface {
//	DefaultNode
//
//}
