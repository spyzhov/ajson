package tokens

import (
	"fmt"

	"github.com/spyzhov/ajson/v1/jerrors"
	"github.com/spyzhov/ajson/v1/jsonpath/internal"
)

type Script struct {
	parent Token
	RPN    *RPN
}

var _ Token = (*Script)(nil)

func NewScript(rpn *RPN) (*Script, error) {
	panic("not implemented")
	//return &Script{
	//	RPN: rpn,
	//}, nil
}

func (t *Script) Type() string {
	return "Script"
}

func (t *Script) String() string {
	if t == nil {
		return "<nil>"
	}
	return fmt.Sprintf("(%s)", t.RPN.String())
}

func (t *Script) Token() string {
	if t == nil {
		return "Script(<nil>)"
	}
	return fmt.Sprintf("Script(%s)", t.RPN.Token())
}

func (t *Script) Path() string {
	if t == nil {
		return "(<nil>)"
	}
	return fmt.Sprintf("(%s)", t.RPN.String())
}

func (t *Script) Parent() Token {
	if t == nil {
		return nil
	}
	return t.parent
}

func (t *Script) SetParent(parent Token) {
	if t == nil {
		return
	}
	t.parent = parent
}

func (t *Script) Append(token Token) error {
	if rpn, ok := token.(*RPN); ok {
		token.SetParent(t)
		t.RPN = rpn
		return nil
	}
	return fmt.Errorf("%w: for Script only RPN is available, %s given", jerrors.ErrUnexpectedStatement, token.Type())
}

func (t *Script) GetState(_ internal.State) internal.State {
	return internal.ѢѢ // fixme
}
