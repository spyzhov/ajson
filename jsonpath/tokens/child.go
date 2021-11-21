package tokens

import (
	"fmt"
	"strings"
)

type Child struct {
	Tokens []Token
}

var _ Token = (*Child)(nil)

func NewChild(tokens []Token) (*Child, error) {
	return &Child{
		Tokens: tokens,
	}, nil
}

func (t *Child) Type() string {
	return "Child"
}

func (t *Child) String() string {
	if t == nil {
		return "<nil>"
	}
	parts := make([]string, 0, len(t.Tokens))
	for _, filter := range t.Tokens {
		parts = append(parts, filter.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(parts, ","))
}

func (t *Child) Token() string {
	if t == nil {
		return "Child(<nil>)"
	}
	parts := make([]string, 0, len(t.Tokens))
	for _, token := range t.Tokens {
		parts = append(parts, token.Token())
	}
	return fmt.Sprintf("Child(%s)", strings.Join(parts, ","))
}

func (t *Child) Path() string {
	if t == nil {
		return "[<nil>]"
	}
	parts := make([]string, 0, len(t.Tokens))
	for _, token := range t.Tokens {
		if path, ok := token.(Path); ok {
			parts = append(parts, path.Path())
		} else {
			parts = append(parts, token.String())
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(parts, ","))
}
