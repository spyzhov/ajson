package tokens

import (
	"fmt"

	"github.com/spyzhov/ajson/v1/internal"
)

type String struct {
	Value string
}

var _ Token = (*String)(nil)

func NewString(value []byte) (*String, error) {
	str, ok := internal.Unquote(value, value[0])
	if !ok {
		return nil, fmt.Errorf("value %q can't be parsed as string", value)
	}
	return &String{
		Value: str,
	}, nil
}

func newString(b *internal.Buffer) (*String, error) {
	start := b.Index
	err := b.AsString(b.Bytes[b.Index], true)
	if err != nil {
		return nil, fmt.Errorf("can't parse string value: %w", err)
	}
	return NewString(b.Bytes[start : b.Index+1] /* with quotes */)
}

func (t *String) Type() string {
	return "String"
}

func (t *String) String() string {
	if t == nil {
		return "String(<nil>)"
	}
	return fmt.Sprintf("String(%q)", t.Value)
}

func (t *String) Token() string {
	return t.String()
}
