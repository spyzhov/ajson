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

func (o *String) String() string {
	if o == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%q", o.Value)
}
