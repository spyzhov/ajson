package tokens

import (
	"fmt"

	"github.com/spyzhov/ajson/v1/jerrors"
	"github.com/spyzhov/ajson/v1/jsonpath/internal"
)

type Filter struct {
	parent Token
	RPN    *RPN
}

var _ Token = (*Filter)(nil)

func NewFilter(rpn *RPN, parent Token) (*Filter, error) {
	return &Filter{
		parent: parent,
		RPN:    rpn,
	}, nil
}

func (t *Filter) Type() string {
	return "Filter"
}

func (t *Filter) String() string {
	if t == nil {
		return "<nil>"
	}
	return fmt.Sprintf("?(%s)", t.RPN.String())
}

func (t *Filter) Token() string {
	if t == nil {
		return "Filter(<nil>)"
	}
	return fmt.Sprintf("Filter(%s)", t.RPN.Token())
}

func (t *Filter) Path() string {
	if t == nil {
		return "?(<nil>)"
	}
	return fmt.Sprintf("?(%s)", t.RPN.String())
}

func (t *Filter) Parent() Token {
	if t == nil {
		return nil
	}
	return t.parent
}

func (t *Filter) SetParent(parent Token) {
	if t == nil {
		return
	}
	t.parent = parent
}

func (t *Filter) Append(token Token) error {
	if rpn, ok := token.(*RPN); ok {
		token.SetParent(t)
		t.RPN = rpn
		return nil
	}
	return fmt.Errorf("%w: for Filter only RPN is available, %s given", jerrors.ErrUnexpectedStatement, token.Type())
}

func (t *Filter) GetState(_ internal.State) internal.State {
	return internal.ѢѢ // fixme
}
