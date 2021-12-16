package tokens

import (
	"fmt"

	"github.com/spyzhov/ajson/v1/internal"
)

type String struct {
	parent Token
	Value  string
}

var _ Token = (*String)(nil)

func NewByteString(value []byte, parent Token) (*String, error) {
	if len(value) < 2 {
		return nil, fmt.Errorf("value %q is too short to be a string", value)
	}
	str, ok := internal.Unquote(value, value[0])
	if !ok {
		return nil, fmt.Errorf("value %q can't be parsed as string", value)
	}
	return &String{
		parent: parent,
		Value:  str,
	}, nil
}

func NewString(value string, parent Token) *String {
	return &String{
		parent: parent,
		Value:  value,
	}
}

func (t *String) Type() string {
	return "String"
}

func (t *String) String() string {
	if t == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%q", t.Value)
}

func (t *String) Token() string {
	if t == nil {
		return "String(<nil>)"
	}
	return fmt.Sprintf("String(%q)", t.Value)
}

func (t *String) Parent() Token {
	if t == nil {
		return nil
	}
	return t.parent
}
