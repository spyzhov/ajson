package tokens

import (
	"fmt"
	"strings"

	"github.com/spyzhov/ajson/v1/internal"
)

type Arguments struct {
	Tokens []Token
}

var _ Token = (*Arguments)(nil)

func NewArguments(token string) (result *Arguments, err error) {
	return newArguments(internal.NewBuffer([]byte(token)))
}

func newArguments(b *internal.Buffer) (result *Arguments, err error) {
	// todo
	panic("not implemented")
}

func (t *Arguments) Type() string {
	return "Arguments"
}

func (t *Arguments) String() string {
	if t == nil {
		return "Arguments(<nil>)"
	}
	parts := make([]string, 0, len(t.Tokens))
	for _, token := range t.Tokens {
		parts = append(parts, token.String())
	}
	return fmt.Sprintf("Arguments(%s)", strings.Join(parts, ", "))
}

func (t *Arguments) Token() string {
	return t.String()
}
