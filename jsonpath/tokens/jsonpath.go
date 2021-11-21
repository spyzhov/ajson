package tokens

import (
	"fmt"
	"strings"

	"github.com/spyzhov/ajson/v1/internal"
)

type JSONPath struct {
	Tokens []Token
}

var _ Token = (*JSONPath)(nil)

func NewJSONPath(token string) (result *JSONPath, err error) {
	return newJSONPath(internal.NewBuffer([]byte(token)))
}

func newJSONPath(b *internal.Buffer) (result *JSONPath, err error) {
	// todo
	panic("not implemented")
}

func (t *JSONPath) Type() string {
	return "JSONPath"
}

func (t *JSONPath) String() string {
	if t == nil {
		return "JSONPath(<nil>)"
	}
	parts := make([]string, 0, len(t.Tokens))
	for _, token := range t.Tokens {
		parts = append(parts, token.String())
	}
	return fmt.Sprintf("JSONPath(%s)", strings.Join(parts, ", "))
}

func (t *JSONPath) Token() string {
	return t.String()
}
