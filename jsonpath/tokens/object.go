package tokens

import (
	"fmt"
	"strings"

	"github.com/spyzhov/ajson/v1/internal"
)

type Object struct {
	Tokens map[Token]Token
}

var _ Token = (*Object)(nil)

func NewObject(token string) (result *Object, err error) {
	return newObject(internal.NewBuffer([]byte(token)))
}

func newObject(_ *internal.Buffer) (result *Object, err error) {
	// todo
	panic("not implemented")
}

func (t *Object) Type() string {
	return "Object"
}

func (t *Object) String() string {
	if t == nil {
		return "<nil>"
	}
	parts := make([]string, 0, len(t.Tokens))
	for key, value := range t.Tokens {
		parts = append(parts, fmt.Sprintf(`%s: %s`, key.String(), value.String()))
	}
	return fmt.Sprintf("{%s}", strings.Join(parts, ", "))
}

func (t *Object) Token() string {
	if t == nil {
		return "Object(<nil>)"
	}
	parts := make([]string, 0, len(t.Tokens))
	for key, value := range t.Tokens {
		parts = append(parts, fmt.Sprintf(`%s:%s`, key.Token(), value.Token()))
	}
	return fmt.Sprintf("Object(%s)", strings.Join(parts, ","))
}
