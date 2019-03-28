package ajson

import (
	"errors"
	"math"
)

func ExampleAddFunction() {
	AddFunction("array_sum", func(node *Node) (result *Node, err error) {
		if node.IsArray() {
			var (
				sum, num float64
				array    []*Node
			)
			array, err = node.GetArray()
			if err != nil {
				return nil, err
			}
			for _, child := range array {
				if !child.IsNumeric() {
					return nil, errors.New("wrong type")
				}
				num, err = child.GetNumeric()
				if err != nil {
					return
				}
				sum += num
			}
			return NumericNode("array_sum", sum), nil
		}
		return
	})
}

func ExampleAddConstant() {
	AddConstant("SqrtPi", NumericNode("SqrtPi", math.SqrtPi))
}

func ExampleAddOperation() {
	AddOperation("<>", 3, false, func(left *Node, right *Node) (node *Node, err error) {
		res, err := left.Eq(right)
		if err != nil {
			return nil, err
		}
		return BoolNode("neq", !res), nil
	})
}
