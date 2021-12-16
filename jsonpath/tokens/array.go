package tokens

import (
	"fmt"
	"strings"

	"github.com/spyzhov/ajson/v1/internal"
)

type Array struct {
	parent Token
	Tokens []Token
}

var _ Token = (*Array)(nil)

func NewArray(token string) (result *Array, err error) {
	panic("not implemented")
	//return newArray(internal.NewBuffer([]byte(token)))
}

func newArray(b *internal.Buffer) (result *Array, err error) {
	// todo
	panic("not implemented")
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
