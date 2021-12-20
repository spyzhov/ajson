package tokens

import (
	"fmt"
	"strings"

	"github.com/spyzhov/ajson/v1/jsonpath/internal"
)

type Arguments struct {
	parent Token
	Tokens []Token
}

var _ Token = (*Arguments)(nil)

func NewArguments(token string) (result *Arguments, err error) {
	panic("not implemented")
}

func (t *Arguments) Type() string {
	return "Arguments"
}

func (t *Arguments) String() string {
	if t == nil {
		return "<nil>"
	}
	parts := make([]string, 0, len(t.Tokens))
	for _, token := range t.Tokens {
		parts = append(parts, token.String())
	}
	return strings.Join(parts, ", ")
}

func (t *Arguments) Token() string {
	if t == nil {
		return "Arguments(<nil>)"
	}
	parts := make([]string, 0, len(t.Tokens))
	for _, token := range t.Tokens {
		parts = append(parts, token.Token())
	}
	return fmt.Sprintf("Arguments(%s)", strings.Join(parts, ", "))
}

func (t *Arguments) Parent() Token {
	if t == nil {
		return nil
	}
	return t.parent
}

func (t *Arguments) SetParent(parent Token) {
	if t == nil {
		return
	}
	t.parent = parent
}

func (t *Arguments) Append(token Token) error {
	t.Tokens = append(t.Tokens, token)
	token.SetParent(t)
	return nil
}

func (t *Arguments) GetState(_ internal.State) internal.State {
	return internal.ѢѢ // fixme
}
