package ajson

import (
	"strings"
)

// Function - internal left function of JSONPath
type Function func(node *Node) (result *Node, err error)

// Operation - internal script operation of JSONPath
type Operation func(left *Node, right *Node) (result *Node, err error)

// AddFunction add a function for internal JSONPath script
func AddFunction(alias string, function Function) {
	Functions[strings.ToLower(alias)] = function
}

// AddOperation add an operation for internal JSONPath script
func AddOperation(alias string, prior uint8, right bool, operation Operation) {
	alias = strings.ToLower(alias)
	Operations[alias] = operation
	Priority[alias] = prior
	PriorityChar[alias[0]] = true
	if right {
		RightOp[alias] = true
	}
}

// AddConstant add a constant for internal JSONPath script
func AddConstant(alias string, value *Node) {
	Constants[strings.ToLower(alias)] = value
}

func numericFunction(name string, fn func(float float64) float64) Function {
	return func(node *Node) (result *Node, err error) {
		if node.IsNumeric() {
			num, err := node.GetNumeric()
			if err != nil {
				return nil, err
			}
			return valueNode(nil, name, Numeric, fn(num)), nil
		}
		return nil, NewErrorRequest("function '%s' was called from non numeric node", name)
	}
}

func mathFactorial(x uint) uint {
	if x == 0 {
		return 1
	}
	return x * mathFactorial(x-1)
}
