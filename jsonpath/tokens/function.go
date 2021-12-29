package tokens

import (
	"fmt"
	"strings"

	"github.com/spyzhov/ajson/v1"
	"github.com/spyzhov/ajson/v1/jerrors"
	"github.com/spyzhov/ajson/v1/jsonpath/internal"
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

func (t *Function) SetParent(parent Token) {
	if t == nil {
		return
	}
	t.parent = parent
}

func (t *Function) Append(token Token) error {
	if arguments, ok := token.(*Arguments); ok {
		token.SetParent(t)
		t.Arguments = arguments
		return nil
	}
	return fmt.Errorf("%w: for Function only Arguments is available, %s given", jerrors.ErrUnexpectedStatement, token.Type())
}

func (t *Function) IsEmpty() bool {
	return t.Arguments == nil
}

func (t *Function) GetState(_ internal.State) internal.State {
	return internal.ѢѢ // fixme
}
