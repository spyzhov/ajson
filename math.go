package ajson

import (
	"math"
	"strings"
)

// Function - internal left function of JSONPath
type Function func(node *Node) (result *Node, err error)

// Function - internal script operation of JSONPath
type Operation func(left *Node, right *Node) (result *Node, err error)

var (
	// Operator precedence
	// From https://golang.org/ref/spec#Operator_precedence
	//
	//	Precedence    Operator
	//	    5             *  /  %  <<  >>  &  &^
	//	    4             +  -  |  ^
	//	    3             ==  !=  <  <=  >  >=
	//	    2             &&
	//	    1             ||
	//
	// Arithmetic operators
	// From https://golang.org/ref/spec#Arithmetic_operators
	//
	//	+    sum                    integers, floats, complex values, strings
	//	-    difference             integers, floats, complex values
	//	*    product                integers, floats, complex values
	//	/    quotient               integers, floats, complex values
	//	%    remainder              integers
	//
	//	&    bitwise AND            integers
	//	|    bitwise OR             integers
	//	^    bitwise XOR            integers
	//	&^   bit clear (AND NOT)    integers
	//
	//	<<   left shift             integer << unsigned integer
	//	>>   right shift            integer >> unsigned integer
	//
	priority = map[string]uint8{
		"**": 6, // additional: power
		"*":  5,
		"/":  5,
		"%":  5,
		"<<": 5,
		">>": 5,
		"&":  5,
		"&^": 5,
		"+":  4,
		"-":  4,
		"|":  4,
		"^":  4,
		"==": 3,
		"!=": 3,
		"<":  3,
		"<=": 3,
		">":  3,
		">=": 3,
		"&&": 2,
		"||": 1,
	}

	rightOp = map[string]bool{
		"**": true,
	}

	operations = map[string]Operation{
		"**": func(left *Node, right *Node) (result *Node, err error) {
			lnum, rnum, err := _floats(left, right)
			if err != nil {
				return
			}
			return valueNode(nil, "power", Numeric, math.Pow(lnum, rnum)), nil
		},
		"*": func(left *Node, right *Node) (result *Node, err error) {
			lnum, rnum, err := _floats(left, right)
			if err != nil {
				return
			}
			return valueNode(nil, "multiply", Numeric, float64(lnum*rnum)), nil
		},
		"/": func(left *Node, right *Node) (result *Node, err error) {
			lnum, rnum, err := _floats(left, right)
			if err != nil {
				return
			}
			if rnum == 0 {
				return nil, errorRequest("division by zero")
			}
			return valueNode(nil, "division", Numeric, float64(lnum/rnum)), nil
		},
		"%": func(left *Node, right *Node) (result *Node, err error) {
			lnum, err := left.getInteger()
			if err != nil {
				return
			}
			rnum, err := left.getInteger()
			if err != nil {
				return
			}
			return valueNode(nil, "remainder", Numeric, float64(lnum%rnum)), nil
		},
		"<<": func(left *Node, right *Node) (result *Node, err error) {
			lnum, err := left.getInteger()
			if err != nil {
				return
			}
			rnum, err := left.getUInteger()
			if err != nil {
				return
			}
			return valueNode(nil, "left shift", Numeric, float64(lnum<<rnum)), nil
		},
		">>": func(left *Node, right *Node) (result *Node, err error) {
			lnum, err := left.getInteger()
			if err != nil {
				return
			}
			rnum, err := left.getUInteger()
			if err != nil {
				return
			}
			return valueNode(nil, "right shift", Numeric, float64(lnum>>rnum)), nil
		},
		"&": func(left *Node, right *Node) (result *Node, err error) {
			lnum, rnum, err := _ints(left, right)
			if err != nil {
				return
			}
			return valueNode(nil, "bitwise AND", Numeric, float64(lnum&rnum)), nil
		},
		"&^": func(left *Node, right *Node) (result *Node, err error) {
			lnum, rnum, err := _ints(left, right)
			if err != nil {
				return
			}
			return valueNode(nil, "bit clear (AND NOT)", Numeric, float64(lnum&rnum)), nil
		},
		"+": func(left *Node, right *Node) (result *Node, err error) {
			lnum, rnum, err := _floats(left, right)
			if err != nil {
				return
			}
			return valueNode(nil, "sum", Numeric, float64(lnum+rnum)), nil
		},
		"-": func(left *Node, right *Node) (result *Node, err error) {
			lnum, rnum, err := _floats(left, right)
			if err != nil {
				return
			}
			return valueNode(nil, "sub", Numeric, float64(lnum-rnum)), nil
		},
		"|": func(left *Node, right *Node) (result *Node, err error) {
			if left.IsNumeric() && right.IsNumeric() {
				lnum, rnum, err := _ints(left, right)
				if err != nil {
					return nil, err
				}
				return valueNode(nil, "bitwise OR", Numeric, float64(lnum|rnum)), nil
			}
			return nil, errorRequest("function 'bitwise OR' was called from non numeric node")
		},
		"^": func(left *Node, right *Node) (result *Node, err error) {
			if left.IsNumeric() && right.IsNumeric() {
				lnum, rnum, err := _ints(left, right)
				if err != nil {
					return nil, err
				}
				return valueNode(nil, "bitwise XOR", Numeric, float64(lnum^rnum)), nil
			}
			return nil, errorRequest("function 'bitwise XOR' was called from non numeric node")
		},
		"==": func(left *Node, right *Node) (result *Node, err error) {
			res, err := left.Eq(right)
			if err != nil {
				return nil, err
			}
			return valueNode(nil, "eq", Bool, bool(res)), nil
		},
		"!=": func(left *Node, right *Node) (result *Node, err error) {
			res, err := left.Eq(right)
			if err != nil {
				return nil, err
			}
			return valueNode(nil, "neq", Bool, bool(!res)), nil
		},
		"<": func(left *Node, right *Node) (result *Node, err error) {
			res, err := left.Le(right)
			if err != nil {
				return nil, err
			}
			return valueNode(nil, "le", Bool, bool(res)), nil
		},
		"<=": func(left *Node, right *Node) (result *Node, err error) {
			res, err := left.Leq(right)
			if err != nil {
				return nil, err
			}
			return valueNode(nil, "leq", Bool, bool(res)), nil
		},
		">": func(left *Node, right *Node) (result *Node, err error) {
			res, err := left.Ge(right)
			if err != nil {
				return nil, err
			}
			return valueNode(nil, "ge", Bool, bool(res)), nil
		},
		">=": func(left *Node, right *Node) (result *Node, err error) {
			res, err := left.Geq(right)
			if err != nil {
				return nil, err
			}
			return valueNode(nil, "geq", Bool, bool(res)), nil
		},
		"&&": func(left *Node, right *Node) (result *Node, err error) {
			res := false
			lval, err := boolean(left)
			if err != nil {
				return nil, err
			}
			if lval {
				rval, err := boolean(right)
				if err != nil {
					return nil, err
				}
				res = rval
			}
			return valueNode(nil, "AND", Bool, bool(res)), nil
		},
		"||": func(left *Node, right *Node) (result *Node, err error) {
			res := true
			lval, err := boolean(left)
			if err != nil {
				return nil, err
			}
			if !lval {
				rval, err := boolean(right)
				if err != nil {
					return nil, err
				}
				res = rval
			}
			return valueNode(nil, "OR", Bool, bool(res)), nil
		},
	}

	functions = map[string]Function{
		"sin": func(node *Node) (result *Node, err error) {
			if node.IsNumeric() {
				num, err := node.GetNumeric()
				if err != nil {
					return nil, err
				}
				return valueNode(nil, "sin", Numeric, math.Sin(num)), nil
			}
			return nil, errorRequest("function 'sin' was called from non numeric node")
		},
		"cos": func(node *Node) (result *Node, err error) {
			if node.IsNumeric() {
				num, err := node.GetNumeric()
				if err != nil {
					return nil, err
				}
				return valueNode(nil, "cos", Numeric, math.Cos(num)), nil
			}
			return nil, errorRequest("function 'cos' was called from non numeric node")
		},
		"length": func(node *Node) (result *Node, err error) {
			if node.IsArray() {
				return valueNode(node, "length", Numeric, float64(node.Size())), nil
			}
			return nil, errorRequest("function 'length' was called from non array node")
		},
		"factorial": func(node *Node) (result *Node, err error) {
			num, err := node.getUInteger()
			if err != nil {
				return
			}
			return valueNode(nil, "factorial", Numeric, float64(mathFactorial(num))), nil
		},
	}
	constants = map[string]*Node{
		"pi":    valueNode(nil, "pi", Numeric, float64(math.Pi)),
		"e":     valueNode(nil, "e", Numeric, float64(math.E)),
		"true":  valueNode(nil, "true", Bool, true),
		"false": valueNode(nil, "false", Bool, false),
		"null":  valueNode(nil, "null", Null, nil),
	}
)

// AddFunction - add a function for internal JSONPath script
func AddFunction(alias string, function Function) {
	functions[strings.ToLower(alias)] = function
}

// AddOperation - add an operation for internal JSONPath script
func AddOperation(alias string, prior uint8, right bool, operation Operation) {
	alias = strings.ToLower(alias)
	operations[alias] = operation
	priority[alias] = prior
	if right {
		rightOp[alias] = true
	}
}

// AddOperation - add a constant for internal JSONPath script
func AddConstant(alias string, value *Node) {
	constants[strings.ToLower(alias)] = value
}
