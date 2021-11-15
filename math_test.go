package ajson

import (
	"errors"
	"fmt"
	"math"
	"testing"
)

func ExampleSetFunction() {
	_ = SetFunction("array_sum", func(node *Node) (result *Node, err error) {
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
			return NewNumeric(sum), nil
		}
		return
	})
}

func ExampleSetFunction_usage() {
	json := []byte(`{"prices": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]}`)
	root, err := Unmarshal(json)
	if err != nil {
		panic(err)
	}
	result, err := Eval(root, `avg($.prices)`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Avg price: %0.1f", result.MustNumeric())
	// Output:
	// Avg price: 5.5
}

func ExampleSetConstant() {
	_ = SetConstant("SqrtPi", NewNumeric(math.SqrtPi))
}

func ExampleSetConstant_using() {
	json := []byte(`{"foo": [true, null, false, 1, "bar", true, 1e3], "bar": [true, "baz", false]}`)
	result, err := JSONPath(json, `$..[?(@ == true)]`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Count of `true` values: %d", len(result))
	// Output:
	// Count of `true` values: 3
}

func ExampleSetConstant_eval() {
	json := []byte(`{"radius": 50, "position": [56.4772531, 84.9918139]}`)
	root, err := Unmarshal(json)
	if err != nil {
		panic(err)
	}
	result, err := Eval(root, `2 * $.radius * pi`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Circumference: %0.3f m.", result.MustNumeric())
	// Output:
	// Circumference: 314.159 m.
}

func ExampleSetOperation() {
	_ = SetOperation("<>", func(left *Node, right *Node) (node *Node, err error) {
		res, err := left.Eq(right)
		if err != nil {
			return nil, err
		}
		return NewBool(!res), nil
	}, 3, false)
}

func ExampleSetOperation_regex() {
	json := []byte(`[{"name":"Foo","mail":"foo@example.com"},{"name":"bar","mail":"bar@example.org"}]`)
	result, err := JSONPath(json, `$.[?(@.mail =~ '.+@example\\.com')]`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("JSON: %s", result[0].Source())
	// Output:
	// JSON: {"name":"Foo","mail":"foo@example.com"}
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
		{name: "2" + operator + "2", operation: operator, left: NewNumeric(2), right: NewNumeric(2), result: NewNumeric(results[0])},
		{name: "3" + operator + "3", operation: operator, left: NewNumeric(3), right: NewNumeric(3), result: NewNumeric(results[1])},
		{name: "10" + operator + "3", operation: operator, left: NewNumeric(10), right: NewNumeric(3), result: NewNumeric(results[2])},
		{name: "X" + operator + "2", operation: operator, left: NewString("X"), right: NewNumeric(2), fail: true},
		{name: "2" + operator + "Y", operation: operator, left: NewNumeric(2), right: NewString("Y"), fail: true},
	}
}

func testBoolOperation(operator string, results [4]bool) []*operationTest {
	return []*operationTest{
		{name: "2" + operator + "2", operation: operator, left: NewNumeric(2), right: NewNumeric(2), result: NewBool(results[0])},
		{name: "3" + operator + "3", operation: operator, left: NewNumeric(3), right: NewNumeric(3), result: NewBool(results[1])},
		{name: "10" + operator + "0", operation: operator, left: NewNumeric(10), right: NewNumeric(0), result: NewBool(results[2])},
		{name: "0" + operator + "10", operation: operator, left: NewNumeric(0), right: NewNumeric(10), result: NewBool(results[3])},
		{name: "left error: " + operator, operation: operator, left: valueNode(nil, "", Numeric, "foo"), right: NewNumeric(10), fail: true},
		{name: "right error: " + operator, operation: operator, left: NewNumeric(10), right: valueNode(nil, "", Numeric, "foo"), fail: true},
	}
}
func testBooleanOperation(operator string, results [4]bool) []*operationTest {
	return []*operationTest{
		{name: "2" + operator + "2", operation: operator, left: NewNumeric(2), right: NewNumeric(2), result: NewBool(results[0])},
		{name: "3" + operator + "3", operation: operator, left: NewNumeric(3), right: NewNumeric(3), result: NewBool(results[1])},
		{name: "10" + operator + "0", operation: operator, left: NewNumeric(10), right: NewNumeric(0), result: NewBool(results[2])},
		{name: "0" + operator + "10", operation: operator, left: NewNumeric(0), right: NewNumeric(10), result: NewBool(results[3])},
	}
}

func TestOperations(t *testing.T) {
	tests := []*operationTest{
		{name: "0/0", operation: "/", left: NewNumeric(0), right: NewNumeric(0), fail: true},
		{name: "1/0", operation: "/", left: NewNumeric(1), right: NewNumeric(0), fail: true},
		{name: "X+Y", operation: "+", left: NewString("X"), right: NewString("Y"), result: NewString("XY")},
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
	_t := NewNumeric(1)
	_f := NewNumeric(0)
	_false := NewBool(false)
	_true := NewBool(true)
	tests = append(
		tests,
		&operationTest{name: "error && true", operation: "&&", left: _e, right: _t, fail: true},
		&operationTest{name: "error && error", operation: "&&", left: _e, right: _e, fail: true},
		&operationTest{name: "error && false", operation: "&&", left: _e, right: _f, fail: true},
		&operationTest{name: "false && error", operation: "&&", left: _f, right: _e, result: _false},
		&operationTest{name: "true && error", operation: "&&", left: _t, right: _e, fail: true},
		&operationTest{
			name:      "[] && {} == false",
			operation: "&&",
			left:      NewArray([]*Node{}),
			right:     NewObject(map[string]*Node{}),
			result:    _false,
		},
		&operationTest{
			name:      "{} || [] == false",
			operation: "||",
			left:      NewObject(map[string]*Node{}),
			right:     NewArray([]*Node{}),
			result:    _false,
		},
		&operationTest{
			name:      `{"foo":"bar"} || [1] == true`,
			operation: "&&",
			left:      NewObject(map[string]*Node{"foo": NewString("bar")}),
			right:     NewArray([]*Node{NewNumeric(1)}),
			result:    _true,
		},

		&operationTest{name: "error || true", operation: "||", left: _e, right: _t, fail: true},
		&operationTest{name: "error || error", operation: "||", left: _e, right: _e, fail: true},
		&operationTest{name: "error || false", operation: "||", left: _e, right: _f, fail: true},
		&operationTest{name: "false || error", operation: "||", left: _f, right: _e, fail: true},
		&operationTest{name: "true || error", operation: "||", left: _t, right: _e, result: _true},

		&operationTest{name: "regexp true", operation: "=~", left: NewString(`123`), right: NewString(`\d+`), result: _true},
		&operationTest{name: "regexp false", operation: "=~", left: NewString(`1 2 3`), right: NewString(`^\d+$`), result: _false},
		&operationTest{name: "regexp pattern error", operation: "=~", left: NewString(`2`), right: NewString(`\2`), fail: true},
		&operationTest{name: "regexp error 1", operation: "=~", left: _f, right: NewString(`123`), fail: true},
		&operationTest{name: "regexp error 2", operation: "=~", left: NewString(`\d+`), right: _f, fail: true},
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Operations[test.operation](test.left, test.right)
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

func TestAddConstant(t *testing.T) {
	name := "new_constant_name"
	if _, ok := Constants[name]; ok {
		t.Error("test constant already exists")
		return
	}
	err := SetConstant(name, NewNumeric(3.14))
	if err != nil {
		t.Errorf("SetConstant() error: %s", err)
		return
	}
	if _, ok := Constants[name]; !ok {
		t.Error("test constant was not added")
	}
}

func TestSetOperation(t *testing.T) {
	name := "!!!!!!"
	if _, ok := Operations[name]; ok {
		t.Error("test operation already exists")
		return
	}
	err := SetOperation(name, func(left *Node, right *Node) (result *Node, err error) {
		return NewNumeric(1), nil
	}, 1, true)
	if err != nil {
		t.Errorf("SetOperation() error: %s", err)
		return
	}
	if _, ok := Operations[name]; !ok {
		t.Error("test operation was not added")
	}
}

func TestAddFunction(t *testing.T) {
	name := "new_function_name"
	if _, ok := Functions[name]; ok {
		t.Error("test constant already exists")
		return
	}
	err := SetFunction(name, func(node *Node) (result *Node, err error) {
		return NewNumeric(2), nil
	})
	if err != nil {
		t.Errorf("SetFunction() error: %s", err)
		return
	}
	if _, ok := Functions[name]; !ok {
		t.Error("test function was not added")
	}
}

func TestFunctions(t *testing.T) {
	tests := []struct {
		name   string
		fname  string
		value  float64
		result interface{}
	}{
		{name: "abs 1", fname: "abs", value: float64(-100), result: 100},
		{name: "abs 2", fname: "abs", value: float64(100), result: 100},
		{name: "abs 3", fname: "abs", value: float64(0), result: 0},
		{name: "acos 1", fname: "acos", value: float64(0.5), result: math.Acos(0.5)},
		{name: "acosh 1", fname: "acosh", value: float64(100), result: math.Acosh(100)},
		{name: "asin 1", fname: "asin", value: float64(0.5), result: math.Asin(0.5)},
		{name: "asinh 1", fname: "asinh", value: float64(100), result: math.Asinh(100)},
		{name: "atan 1", fname: "atan", value: float64(100), result: math.Atan(100)},
		{name: "atanh 1", fname: "atanh", value: float64(0.5), result: math.Atanh(0.5)},
		{name: "cbrt 1", fname: "cbrt", value: float64(100), result: math.Cbrt(100)},
		{name: "ceil 1", fname: "ceil", value: float64(100), result: math.Ceil(100)},
		{name: "cos 1", fname: "cos", value: float64(100), result: math.Cos(100)},
		{name: "cosh 1", fname: "cosh", value: float64(100), result: math.Cosh(100)},
		{name: "erf 1", fname: "erf", value: float64(100), result: math.Erf(100)},
		{name: "erfc 1", fname: "erfc", value: float64(100), result: math.Erfc(100)},
		{name: "erfcinv 1", fname: "erfcinv", value: float64(0.5), result: math.Erfcinv(0.5)},
		{name: "erfinv 1", fname: "erfinv", value: float64(0.5), result: math.Erfinv(0.5)},
		{name: "exp 1", fname: "exp", value: float64(100), result: math.Exp(100)},
		{name: "exp2 1", fname: "exp2", value: float64(100), result: math.Exp2(100)},
		{name: "expm1 1", fname: "expm1", value: float64(100), result: math.Expm1(100)},
		{name: "floor 1", fname: "floor", value: float64(0), result: math.Floor(0)},
		{name: "floor 2", fname: "floor", value: float64(0.1), result: math.Floor(0.1)},
		{name: "floor 3", fname: "floor", value: float64(0.5), result: math.Floor(0.5)},
		{name: "floor 4", fname: "floor", value: float64(0.9), result: math.Floor(0.9)},
		{name: "floor 5", fname: "floor", value: float64(100), result: math.Floor(100)},
		{name: "gamma 1", fname: "gamma", value: float64(100), result: math.Gamma(100)},
		{name: "j0 1", fname: "j0", value: float64(100), result: math.J0(100)},
		{name: "j1 1", fname: "j1", value: float64(100), result: math.J1(100)},
		{name: "log 1", fname: "log", value: float64(100), result: math.Log(100)},
		{name: "log10 1", fname: "log10", value: float64(100), result: math.Log10(100)},
		{name: "log1p 1", fname: "log1p", value: float64(100), result: math.Log1p(100)},
		{name: "log2 1", fname: "log2", value: float64(100), result: math.Log2(100)},
		{name: "logb 1", fname: "logb", value: float64(100), result: math.Logb(100)},
		{name: "round 1", fname: "round", value: float64(0), result: math.Round(0)},
		{name: "round 2", fname: "round", value: float64(0.1), result: math.Round(0.1)},
		{name: "round 3", fname: "round", value: float64(0.5), result: math.Round(0.5)},
		{name: "round 4", fname: "round", value: float64(0.9), result: math.Round(0.9)},
		{name: "round 5", fname: "round", value: float64(100), result: math.Round(100)},
		{name: "roundtoeven 1", fname: "roundtoeven", value: float64(0), result: math.RoundToEven(0)},
		{name: "roundtoeven 2", fname: "roundtoeven", value: float64(0.5), result: math.RoundToEven(0.5)},
		{name: "roundtoeven 3", fname: "roundtoeven", value: float64(0.1), result: math.RoundToEven(0.1)},
		{name: "roundtoeven 4", fname: "roundtoeven", value: float64(0.9), result: math.RoundToEven(0.9)},
		{name: "roundtoeven 5", fname: "roundtoeven", value: float64(1), result: math.RoundToEven(1)},
		{name: "sin 1", fname: "sin", value: float64(100), result: math.Sin(100)},
		{name: "sinh 1", fname: "sinh", value: float64(100), result: math.Sinh(100)},
		{name: "sqrt 1", fname: "sqrt", value: float64(100), result: math.Sqrt(100)},
		{name: "tan 1", fname: "tan", value: float64(100), result: math.Tan(100)},
		{name: "tanh 1", fname: "tanh", value: float64(100), result: math.Tanh(100)},
		{name: "trunc 1", fname: "trunc", value: float64(100), result: math.Trunc(100)},
		{name: "y0 1", fname: "y0", value: float64(100), result: math.Y0(100)},
		{name: "y1 1", fname: "y1", value: float64(100), result: math.Y1(100)},

		{name: "pow10", fname: "pow10", value: float64(10), result: math.Pow10(10)},
		{name: "factorial", fname: "factorial", value: float64(10), result: 3628800},

		{name: "not_1", fname: "not", value: float64(1), result: false},
		{name: "not_0", fname: "not", value: float64(0), result: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			node := NewNumeric(test.value)
			var expected *Node
			switch test.result.(type) {
			case int:
				expected = NewNumeric(float64(test.result.(int)))
			case float64:
				expected = NewNumeric(test.result.(float64))
			case bool:
				expected = NewBool(test.result.(bool))
			default:
				panic("wrong type")
			}
			result, err := Functions[test.fname](node)
			if err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
			} else if ok, err := result.Eq(expected); !ok {
				if err != nil {
					t.Errorf("Unexpected error on comparation: %s", err.Error())
				}
				t.Errorf("Wrong value: %v != %v", result.value.Load(), expected.value.Load())
			}
		})
	}
}

func TestFunctions2(t *testing.T) {
	_e := valueNode(nil, "", Numeric, "foo")
	_s := valueNode(nil, "", String, true)
	tests := []struct {
		name   string
		fname  string
		value  *Node
		result *Node
		fail   bool
	}{
		{name: "pow10 error", fname: "pow10", value: _e, fail: true},
		{name: "factorial error", fname: "factorial", value: _e, fail: true},
		{name: "abs error 1", fname: "abs", value: _e, fail: true},
		{name: "abs error 2", fname: "abs", value: NewString(""), fail: true},

		{name: "length array", fname: "length", value: NewArray([]*Node{
			valueNode(nil, "", Numeric, "foo"),
			valueNode(nil, "", Numeric, "foo"),
			valueNode(nil, "", Numeric, "foo"),
		}), result: NewNumeric(3)},
		{name: "length blank array", fname: "length", value: NewArray([]*Node{}), result: NewNumeric(0)},
		{name: "length object", fname: "length", value: NewObject(map[string]*Node{
			"foo": NewNumeric(1),
			"bar": NewNumeric(1),
		}), result: NewNumeric(2)},
		{name: "length string", fname: "length", value: NewString("foo_bar"), result: NewNumeric(7)},
		{name: "length string error", fname: "length", value: _s, fail: true},
		{name: "length numeric", fname: "length", value: NewNumeric(123), result: NewNumeric(1)},
		{name: "length bool", fname: "length", value: NewBool(false), result: NewNumeric(1)},
		{name: "length null", fname: "length", value: NewNull(), result: NewNumeric(1)},

		{name: "avg error 1", fname: "avg", value: NewArray([]*Node{
			valueNode(nil, "", Numeric, "foo"),
			valueNode(nil, "", Numeric, "foo"),
			valueNode(nil, "", Numeric, "foo"),
		}), fail: true},
		{name: "avg error 2", fname: "avg", value: _e, fail: false, result: NewNull()},
		{name: "avg array 1", fname: "avg", value: NewArray([]*Node{
			NewNumeric(1),
			NewNumeric(1),
			NewNumeric(1),
			NewNumeric(1),
		}), result: NewNumeric(1)},
		{name: "avg array 2", fname: "avg", value: NewArray([]*Node{
			NewNumeric(1),
			NewNumeric(2),
			NewNumeric(3),
		}), result: NewNumeric(2)},
		{name: "avg object", fname: "avg", value: NewObject(map[string]*Node{
			"q": NewNumeric(1),
			"w": NewNumeric(2),
			"e": NewNumeric(3),
		}), result: NewNumeric(2)},
		{name: "avg array blank", fname: "avg", value: NewArray([]*Node{}), result: NewNumeric(0)},

		{name: "sum error 1", fname: "sum", value: NewArray([]*Node{
			valueNode(nil, "", Numeric, "foo"),
			valueNode(nil, "", Numeric, "foo"),
			valueNode(nil, "", Numeric, "foo"),
		}), fail: true},
		{name: "sum error 2", fname: "sum", value: _e, fail: false, result: NewNull()},
		{name: "sum array 1", fname: "sum", value: NewArray([]*Node{
			NewNumeric(1),
			NewNumeric(1),
			NewNumeric(1),
			NewNumeric(1),
		}), result: NewNumeric(4)},
		{name: "sum array 2", fname: "sum", value: NewArray([]*Node{
			NewNumeric(1),
			NewNumeric(2),
			NewNumeric(3),
		}), result: NewNumeric(6)},
		{name: "sum object", fname: "sum", value: NewObject(map[string]*Node{
			"q": NewNumeric(1),
			"w": NewNumeric(2),
			"e": NewNumeric(3),
		}), result: NewNumeric(6)},
		{name: "sum array blank", fname: "sum", value: NewArray([]*Node{}), result: NewNumeric(0)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Functions[test.fname](test.value)
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

func TestConstants(t *testing.T) {
	tests := []struct {
		name     string
		expected *Node
	}{
		{name: "e", expected: NewNumeric(float64(math.E))},
		{name: "pi", expected: NewNumeric(float64(math.Pi))},
		{name: "phi", expected: NewNumeric(float64(math.Phi))},

		{name: "sqrt2", expected: NewNumeric(float64(math.Sqrt2))},
		{name: "sqrte", expected: NewNumeric(float64(math.SqrtE))},
		{name: "sqrtpi", expected: NewNumeric(float64(math.SqrtPi))},
		{name: "sqrtphi", expected: NewNumeric(float64(math.SqrtPhi))},

		{name: "ln2", expected: NewNumeric(float64(math.Ln2))},
		{name: "log2e", expected: NewNumeric(float64(math.Log2E))},
		{name: "ln10", expected: NewNumeric(float64(math.Ln10))},
		{name: "log10e", expected: NewNumeric(float64(math.Log10E))},

		{name: "true", expected: NewBool(true)},
		{name: "false", expected: NewBool(false)},
		{name: "null", expected: NewNull()},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Constants[test.name]
			if ok, err := result.Eq(test.expected); !ok {
				if err != nil {
					t.Errorf("Unexpected error on comparation: %s", err.Error())
				}
				t.Errorf("Wrong value: %v != %v", result.value.Load(), test.expected.value.Load())
			}
		})
	}
}
