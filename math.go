package ajson

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	nameRegex     = regexp.MustCompile(`^[a-z_][a-z0-9_]+$`)
	operRegex     = regexp.MustCompile("^[`~!@#%^&*+\\-/:?\\\\\\|<=>]+$")
	operBlackList = map[string]bool{
		`=`: true,
		`+`: true,
		`-`: true,
		`/`: true,
		`*`: true,
	}
)

// Function - internal left function of JSONPath.
type Function func(node *Node) (result *Node, err error)

// Operation - internal script operation of JSONPath.
type Operation func(left *Node, right *Node) (result *Node, err error)

// SetFunction add/replace a function into the internal JSONPath script.
// Function name should match the regexp: `^[a-z_][a-z0-9_]+$`
// If this function already exists, it will be replaced.
func SetFunction(alias string, function Function) error {
	alias = strings.ToLower(alias)
	if !nameRegex.MatchString(alias) {
		return fmt.Errorf("function name %q should match the regexp: %q", alias, nameRegex.String())
	}
	Functions[alias] = function
	return nil
}

// SetOperation add an operation for internal JSONPath script.
// Operation name should match the regexp: `^[^\w\s'",.:;]+$`
// If this operation already exists, it will be replaced.
// It is forbidden to use/replace some operations as soon they are the reserved identifiers:
//
//   + - / * =
func SetOperation(alias string, operation Operation, prior uint8, right bool) error {
	alias = strings.ToLower(alias)
	if operBlackList[alias] {
		return fmt.Errorf("operation name %q is forbidden to use as soon it is one of the reserved identifiers", alias)
	}
	if !operRegex.MatchString(alias) {
		return fmt.Errorf("operation name %q should match the regexp: %q", alias, operRegex.String())
	}
	Operations[alias] = operation
	OperationsPriority[alias] = prior
	OperationsChar[alias[0]] = true
	if right {
		RightOp[alias] = true
	}
	return nil
}

// SetConstant add a constant for internal JSONPath script.
// Constant name should match the regexp: `^[a-z_][a-z0-9_]+$`
// If this constant already exists, it will be replaced.
func SetConstant(alias string, value *Node) error {
	alias = strings.ToLower(alias)
	if !nameRegex.MatchString(alias) {
		return fmt.Errorf("constant name %q should match the regexp: %q", alias, nameRegex.String())
	}
	Constants[alias] = value
	return nil
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
