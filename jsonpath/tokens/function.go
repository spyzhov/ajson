package tokens

import (
	"fmt"
	"strings"

	"github.com/spyzhov/ajson/v1"
)

type Function struct {
	Alias string
}

var _ Token = (*Function)(nil)

func NewFunction(alias string) (*Function, error) {
	alias = strings.ToLower(alias)
	if _, ok := ajson.Functions[alias]; !ok {
		return nil, fmt.Errorf("function %q not found", alias)
	}
	return &Function{
		Alias: alias,
	}, nil
}

func (o *Function) String() string {
	if o == nil {
		return "<nil>"
	}
	return o.Alias
}

func (o *Function) Function() ajson.Function {
	if o == nil {
		return nil
	}
	return ajson.Functions[o.Alias]
}
