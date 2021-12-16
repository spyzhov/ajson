package tokens

import (
	"fmt"
	"strings"

	"github.com/spyzhov/ajson/v1"
)

type Function struct {
	parent    Token
	Alias     string
	Arguments *Arguments
}

var _ Token = (*Function)(nil)

func NewFunction(alias string, arguments *Arguments) (*Function, error) {
	alias = strings.ToLower(alias)
	if _, ok := ajson.Functions[alias]; !ok {
		return nil, fmt.Errorf("function %q not found", alias)
	}
	return &Function{
		Alias:     alias,
		Arguments: arguments,
	}, nil
}

func (t *Function) Type() string {
	return "Function"
}

func (t *Function) String() string {
	if t == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s(%s)", t.Alias, t.Arguments.String())
}

func (t *Function) Token() string {
	if t == nil {
		return "Function(<nil>, <nil>)"
	}
	return fmt.Sprintf("Function(%s, %s)", t.Alias, t.Arguments.Token())
}

func (t *Function) Function() ajson.Function {
	if t == nil {
		return nil
	}
	return ajson.Functions[t.Alias]
}

func (t *Function) Parent() Token {
	if t == nil {
		return nil
	}
	return t.parent
}
