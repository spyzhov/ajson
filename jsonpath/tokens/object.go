package tokens

import (
	"fmt"
	"strings"

	"github.com/spyzhov/ajson/v1/jerrors"
)

type Object struct {
	parent Token
	Tokens []*ObjectElement
}

var _ Token = (*Object)(nil)

func NewObject(parent Token) (*Object, error) {
	return &Object{
		parent: parent,
		Tokens: make([]*ObjectElement, 0),
	}, nil
}

func (t *Object) Append(token Token) error {
	if element, ok := token.(*ObjectElement); ok {
		t.Tokens = append(t.Tokens, element)
		return nil
	}
	return fmt.Errorf("%w: for Object only ObjectElement is available, %s given", jerrors.ErrUnexpectedStatement, token.Type())
}

func (t *Object) Type() string {
	return "Object"
}

func (t *Object) String() string {
	if t == nil {
		return "<nil>"
	}
	parts := make([]string, 0, len(t.Tokens))
	for _, value := range t.Tokens {
		parts = append(parts, value.String())
	}
	return fmt.Sprintf("{%s}", strings.Join(parts, ", "))
}

func (t *Object) Token() string {
	if t == nil {
		return "Object(<nil>)"
	}
	parts := make([]string, 0, len(t.Tokens))
	for _, value := range t.Tokens {
		parts = append(parts, value.Token())
	}
	return fmt.Sprintf("Object(%s)", strings.Join(parts, ","))
}

func (t *Object) Parent() Token {
	if t == nil {
		return nil
	}
	return t.parent
}
