package tokens

import (
	"fmt"
	"strings"

	"github.com/spyzhov/ajson/v1/jsonpath/internal"
)

type Array struct {
	parent Token
	Tokens []Token
}

var _ Token = (*Array)(nil)

func NewArray() *Array {
	return &Array{
		Tokens: make([]Token, 0),
	}
}

func (t *Array) Type() string {
	return "Array"
}

func (t *Array) String() string {
	if t == nil {
		return "<nil>"
	}
	parts := make([]string, 0, len(t.Tokens))
	for _, token := range t.Tokens {
		parts = append(parts, token.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(parts, ", "))
}

func (t *Array) Token() string {
	if t == nil {
		return "Array(<nil>)"
	}
	parts := make([]string, 0, len(t.Tokens))
	for _, token := range t.Tokens {
		parts = append(parts, token.Token())
	}
	return fmt.Sprintf("Array(%s)", strings.Join(parts, ", "))
}

func (t *Array) Parent() Token {
	if t == nil {
		return nil
	}
	return t.parent
}

func (t *Array) SetParent(parent Token) {
	if t == nil {
		return
	}
	t.parent = parent
}

func (t *Array) Append(token Token) error {
	t.Tokens = append(t.Tokens, token)
	token.SetParent(t)
	return nil
}

func (t *Array) IsEmpty() bool {
	return len(t.Tokens) == 0
}

func (t *Array) GetState(_ internal.State) internal.State {
	return internal.ѢѢ // fixme
}
