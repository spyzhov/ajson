package tokens

import (
	"fmt"

	"github.com/spyzhov/ajson/v1/jerrors"
	"github.com/spyzhov/ajson/v1/jsonpath/internal"
)

type ObjectElement struct {
	parent Token
	Key    Token
	Value  Token
}

var _ Token = (*ObjectElement)(nil)

func NewObjectElement(key Token, value Token, parent Token) (result *ObjectElement, err error) {
	if _, ok := parent.(*Object); !ok {
		return nil, fmt.Errorf("%w: object element outside of object", jerrors.ErrUnexpectedStatement)
	}
	return &ObjectElement{
		parent: parent,
		Key:    key,
		Value:  value,
	}, nil
}

func (t *ObjectElement) Type() string {
	return "ObjectElement"
}

func (t *ObjectElement) String() string {
	if t == nil {
		return "<nil>"
	}
	return fmt.Sprintf(`%s: %s`, t.Key.String(), t.Value.String())
}

func (t *ObjectElement) Token() string {
	if t == nil {
		return "ObjectElement(<nil>)"
	}
	return fmt.Sprintf(`%s:%s`, t.Key.Token(), t.Value.Token())
}

func (t *ObjectElement) Parent() Token {
	if t == nil {
		return nil
	}
	return t.parent
}

func (t *ObjectElement) SetParent(parent Token) {
	if t == nil {
		return
	}
	t.parent = parent
}

func (t *ObjectElement) Append(token Token) error {
	if t.Value != nil {
		return fmt.Errorf("%w: object element value already filled with %q, new element %q given", jerrors.ErrIncorrectJSONPath, t.Value.Token(), token.Token())
	}
	t.Value = token
	token.SetParent(t)
	return nil
}

func (t *ObjectElement) IsEmpty() bool {
	return t.Value == nil
}

func (t *ObjectElement) GetState(_ internal.State) internal.State {
	return internal.ѢѢ // fixme
}
