package ajson

import (
	"errors"
	"math"
	"testing"
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

type operationTest struct {
	name      string
	operation string
	left      *Node
	right     *Node
	result    *Node
	fail      bool
}

func testNumOperation(operator string, results [3]float64) []*operationTest {
	return []*operationTest{
		{name: "2" + operator + "2", operation: operator, left: NumericNode("", 2), right: NumericNode("", 2), result: NumericNode("", results[0])},
		{name: "3" + operator + "3", operation: operator, left: NumericNode("", 3), right: NumericNode("", 3), result: NumericNode("", results[1])},
		{name: "10" + operator + "3", operation: operator, left: NumericNode("", 10), right: NumericNode("", 3), result: NumericNode("", results[2])},
		{name: "X" + operator + "2", operation: operator, left: StringNode("", "X"), right: NumericNode("", 2), fail: true},
		{name: "2" + operator + "Y", operation: operator, left: NumericNode("", 2), right: StringNode("", "Y"), fail: true},
	}
}

func testBoolOperation(operator string, results [4]bool) []*operationTest {
	return []*operationTest{
		{name: "2" + operator + "2", operation: operator, left: NumericNode("", 2), right: NumericNode("", 2), result: BoolNode("", results[0])},
		{name: "3" + operator + "3", operation: operator, left: NumericNode("", 3), right: NumericNode("", 3), result: BoolNode("", results[1])},
		{name: "10" + operator + "0", operation: operator, left: NumericNode("", 10), right: NumericNode("", 0), result: BoolNode("", results[2])},
		{name: "0" + operator + "10", operation: operator, left: NumericNode("", 0), right: NumericNode("", 10), result: BoolNode("", results[3])},
		{name: "left error: " + operator, operation: operator, left: valueNode(nil, "", Numeric, "foo"), right: NumericNode("", 10), fail: true},
		{name: "right error: " + operator, operation: operator, left: NumericNode("", 10), right: valueNode(nil, "", Numeric, "foo"), fail: true},
	}
}
func testBooleanOperation(operator string, results [4]bool) []*operationTest {
	return []*operationTest{
		{name: "2" + operator + "2", operation: operator, left: NumericNode("", 2), right: NumericNode("", 2), result: BoolNode("", results[0])},
		{name: "3" + operator + "3", operation: operator, left: NumericNode("", 3), right: NumericNode("", 3), result: BoolNode("", results[1])},
		{name: "10" + operator + "0", operation: operator, left: NumericNode("", 10), right: NumericNode("", 0), result: BoolNode("", results[2])},
		{name: "0" + operator + "10", operation: operator, left: NumericNode("", 0), right: NumericNode("", 10), result: BoolNode("", results[3])},
	}
}

func TestOperations(t *testing.T) {
	tests := []*operationTest{
		{name: "0/0", operation: "/", left: NumericNode("", 0), right: NumericNode("", 0), fail: true},
		{name: "1/0", operation: "/", left: NumericNode("", 1), right: NumericNode("", 0), fail: true},
		{name: "X+Y", operation: "+", left: StringNode("", "X"), right: StringNode("", "Y"), result: StringNode("", "XY")},
	}
	tests = append(tests, testNumOperation("**", [3]float64{4, 27, 1000})...)

	tests = append(tests, testNumOperation("*", [3]float64{4, 9, 30})...)
	tests = append(tests, testNumOperation("+", [3]float64{4, 6, 13})...)
	tests = append(tests, testNumOperation("-", [3]float64{0, 0, 7})...)
	tests = append(tests, testNumOperation("/", [3]float64{1, 1, 10. / 3.})...)
	tests = append(tests, testNumOperation("%", [3]float64{0, 0, 1})...)

	tests = append(tests, testNumOperation("<<", [3]float64{8, 24, 80})...)
	tests = append(tests, testNumOperation(">>", [3]float64{0, 0, 1})...)
	tests = append(tests, testNumOperation("&", [3]float64{2, 3, 2})...)
	tests = append(tests, testNumOperation("&^", [3]float64{0, 0, 8})...)
	tests = append(tests, testNumOperation("|", [3]float64{2, 3, 11})...)
	tests = append(tests, testNumOperation("^", [3]float64{0, 0, 9})...)

	tests = append(tests, testBoolOperation("==", [4]bool{true, true, false, false})...)
	tests = append(tests, testBoolOperation("!=", [4]bool{false, false, true, true})...)
	tests = append(tests, testBoolOperation("<", [4]bool{false, false, false, true})...)
	tests = append(tests, testBoolOperation("<=", [4]bool{true, true, false, true})...)
	tests = append(tests, testBoolOperation(">", [4]bool{false, false, true, false})...)
	tests = append(tests, testBoolOperation(">=", [4]bool{true, true, true, false})...)

	tests = append(tests, testBooleanOperation("&&", [4]bool{true, true, false, false})...)
	tests = append(tests, testBooleanOperation("||", [4]bool{true, true, true, true})...)

	_e := valueNode(nil, "", Numeric, "foo")
	_t := NumericNode("", 1)
	_f := NumericNode("", 0)
	_false := BoolNode("", false)
	_true := BoolNode("", true)
	tests = append(
		tests,
		&operationTest{name: "error && true", operation: "&&", left: _e, right: _t, fail: true},
		&operationTest{name: "error && error", operation: "&&", left: _e, right: _e, fail: true},
		&operationTest{name: "error && false", operation: "&&", left: _e, right: _f, fail: true},
		&operationTest{name: "false && error", operation: "&&", left: _f, right: _e, result: _false},
		&operationTest{name: "true && error", operation: "&&", left: _t, right: _e, fail: true},

		&operationTest{name: "error || true", operation: "||", left: _e, right: _t, fail: true},
		&operationTest{name: "error || error", operation: "||", left: _e, right: _e, fail: true},
		&operationTest{name: "error || false", operation: "||", left: _e, right: _f, fail: true},
		&operationTest{name: "false || error", operation: "||", left: _f, right: _e, fail: true},
		&operationTest{name: "true || error", operation: "||", left: _t, right: _e, result: _true},
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := operations[test.operation](test.left, test.right)
			if test.fail {
				if err == nil {
					t.Error("Expected error: nil given")
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
			} else if ok, err := result.Eq(test.result); !ok {
				if err != nil {
					t.Errorf("Unexpected error on comparation: %s", err.Error())
				}
				t.Errorf("Wrong value: %v != %v", result.value.Load(), test.result.value.Load())
			}
		})
	}
}
