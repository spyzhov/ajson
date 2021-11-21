package tokens

import (
	"fmt"

	"github.com/spyzhov/ajson/v1/internal"
)

type parentheses bool

var _ Token = (parentheses)(false)

func newParentheses(char byte) (result parentheses, err error) {
	switch char {
	case internal.BParenthesesL:
		result = true
	case internal.BParenthesesR:
		result = false
	default:
		err = fmt.Errorf("given value is not `parentheses`, char index: %d, value: %q", char, char)
	}
	return result, nil
}

func (t parentheses) Type() string {
	return "parentheses"
}

func (t parentheses) IsOpen() bool {
	return bool(t)
}

func (t parentheses) String() string {
	if t.IsOpen() {
		return "("
	}
	return ")"
}

func (t parentheses) Token() string {
	if t.IsOpen() {
		return "parentheses(Open)"
	}
	return "parentheses(Close)"
}
