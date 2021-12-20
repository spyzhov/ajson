package tokens

import (
	"fmt"

	"github.com/spyzhov/ajson/v1/jerrors"
	"github.com/spyzhov/ajson/v1/jsonpath/internal"
)

type Child struct {
	parent   Token
	Selector Token
}

var _ Token = (*Child)(nil)

func NewChild(token Token, parent Token) *Child {
	return &Child{
		parent:   parent,
		Selector: token,
	}
}

func (t *Child) Type() string {
	return "Child"
}

func (t *Child) String() string {
	if t == nil {
		return "<nil>"
	}
	return fmt.Sprintf("[%s]", t.Selector.String())
}

func (t *Child) Token() string {
	if t == nil {
		return "Child(<nil>)"
	}
	return fmt.Sprintf("Child(%s)", t.Selector.Token())
}

func (t *Child) Path() string {
	if t == nil {
		return "[<nil>]"
	}
	if path, ok := t.Selector.(Path); ok {
		return fmt.Sprintf("[%s]", path.Path())
	}
	return fmt.Sprintf("[%s]", t.Selector.String())
}

func (t *Child) Parent() Token {
	if t == nil {
		return nil
	}
	return t.parent
}

func (t *Child) SetParent(parent Token) {
	if t == nil {
		return
	}
	t.parent = parent
}

func (t *Child) Append(token Token) error {
	if t.Selector != nil {
		return fmt.Errorf("%w: child selection already filled with %q, new element %q given", jerrors.ErrIncorrectJSONPath, t.Selector.Token(), token.Token())
	}
	t.Selector = token
	token.SetParent(t)
	return nil
}

func (t *Child) GetState(_ internal.State) internal.State {
	return internal.ѢѢ // fixme
}
