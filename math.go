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
			if left.IsString() {
				lnum, rnum, err := _strings(left, right)
				if err != nil {
					return nil, err
				}
				return valueNode(nil, "sum", String, string(lnum+rnum)), nil
			} else {
				lnum, rnum, err := _floats(left, right)
				if err != nil {
					return nil, err
				}
				return valueNode(nil, "sum", Numeric, float64(lnum+rnum)), nil
			}
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
		"abs":         numericFunction("Abs", math.Abs),
		"acos":        numericFunction("Acos", math.Acos),
		"acosh":       numericFunction("Acosh", math.Acosh),
		"asin":        numericFunction("Asin", math.Asin),
		"asinh":       numericFunction("Asinh", math.Asinh),
		"atan":        numericFunction("Atan", math.Atan),
		"atanh":       numericFunction("Atanh", math.Atanh),
		"cbrt":        numericFunction("Cbrt", math.Cbrt),
		"ceil":        numericFunction("Ceil", math.Ceil),
		"cos":         numericFunction("Cos", math.Cos),
		"cosh":        numericFunction("Cosh", math.Cosh),
		"erf":         numericFunction("Erf", math.Erf),
		"erfc":        numericFunction("Erfc", math.Erfc),
		"erfcinv":     numericFunction("Erfcinv", math.Erfcinv),
		"erfinv":      numericFunction("Erfinv", math.Erfinv),
		"exp":         numericFunction("Exp", math.Exp),
		"exp2":        numericFunction("Exp2", math.Exp2),
		"expm1":       numericFunction("Expm1", math.Expm1),
		"floor":       numericFunction("Floor", math.Floor),
		"gamma":       numericFunction("Gamma", math.Gamma),
		"j0":          numericFunction("J0", math.J0),
		"j1":          numericFunction("J1", math.J1),
		"log":         numericFunction("Log", math.Log),
		"log10":       numericFunction("Log10", math.Log10),
		"log1p":       numericFunction("Log1p", math.Log1p),
		"log2":        numericFunction("Log2", math.Log2),
		"logb":        numericFunction("Logb", math.Logb),
		"round":       numericFunction("Round", math.Round),
		"roundtoeven": numericFunction("RoundToEven", math.RoundToEven),
		"sin":         numericFunction("Sin", math.Sin),
		"sinh":        numericFunction("Sinh", math.Sinh),
		"sqrt":        numericFunction("Sqrt", math.Sqrt),
		"tan":         numericFunction("Tan", math.Tan),
		"tanh":        numericFunction("Tanh", math.Tanh),
		"trunc":       numericFunction("Trunc", math.Trunc),
		"y0":          numericFunction("Y0", math.Y0),
		"y1":          numericFunction("Y1", math.Y1),

		"pow10": func(node *Node) (result *Node, err error) {
			num, err := node.getInteger()
			if err != nil {
				return
			}
			return valueNode(nil, "Pow10", Numeric, float64(math.Pow10(num))), nil
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
		"e":   valueNode(nil, "e", Numeric, float64(math.E)),
		"pi":  valueNode(nil, "pi", Numeric, float64(math.Pi)),
		"phi": valueNode(nil, "phi", Numeric, float64(math.Phi)),

		"sqrt2":   valueNode(nil, "sqrt2", Numeric, float64(math.Sqrt2)),
		"sqrte":   valueNode(nil, "sqrte", Numeric, float64(math.SqrtE)),
		"sqrtpi":  valueNode(nil, "sqrtpi", Numeric, float64(math.SqrtPi)),
		"sqrtphi": valueNode(nil, "sqrtphi", Numeric, float64(math.SqrtPhi)),

		"ln2":    valueNode(nil, "ln2", Numeric, float64(math.Ln2)),
		"log2e":  valueNode(nil, "log2e", Numeric, float64(math.Log2E)),
		"ln10":   valueNode(nil, "ln10", Numeric, float64(math.Ln10)),
		"log10e": valueNode(nil, "log10e", Numeric, float64(math.Log10E)),

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

func numericFunction(name string, fn func(float float64) float64) Function {
	return func(node *Node) (result *Node, err error) {
		if node.IsNumeric() {
			num, err := node.GetNumeric()
			if err != nil {
				return nil, err
			}
			return valueNode(nil, name, Numeric, fn(num)), nil
		}
		return nil, errorRequest("function '%s' was called from non numeric node", name)
	}
}