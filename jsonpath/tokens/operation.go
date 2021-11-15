package tokens

import (
	"fmt"
	"strings"

	"github.com/spyzhov/ajson/v1"
)

type Operation struct {
	Alias string
}

var _ Token = (*Operation)(nil)

func NewOperation(alias string) (*Operation, error) {
	alias = strings.ToLower(alias)
	if _, ok := ajson.Operations[alias]; !ok {
		return nil, fmt.Errorf("operation %q not found", alias)
	}
	return &Operation{
		Alias: alias,
	}, nil
}

func (o *Operation) String() string {
	if o == nil {
		return "<nil>"
	}
	return o.Alias
}

func (o *Operation) Operation() ajson.Operation {
	if o == nil {
		return nil
	}
	return ajson.Operations[o.Alias]
}

func (o *Operation) Priority() uint8 {
	if o == nil {
		return 0
	}
	return ajson.OperationsPriority[o.Alias]
}

func (o *Operation) IsRight() bool {
	if o == nil {
		return false
	}
	return ajson.RightOp[o.Alias]
}
