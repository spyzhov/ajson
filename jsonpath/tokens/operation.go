package tokens

import (
	"fmt"
	"strings"

	"github.com/spyzhov/ajson/v1"
)

type Operation struct {
	parent Token
	Alias  string
}

var _ Token = (*Operation)(nil)

func NewOperation(alias string) (*Operation, error) {
	alias = strings.ToLower(alias)
	// fixme: m.b. remove this dependency
	if _, ok := ajson.Operations[alias]; !ok {
		return nil, fmt.Errorf("operation %q not found", alias)
	}
	return &Operation{
		Alias: alias,
	}, nil
}

func (t *Operation) Type() string {
	return "Operation"
}

func (t *Operation) String() string {
	if t == nil {
		return "<nil>"
	}
	return t.Alias
}

func (t *Operation) Token() string {
	if t == nil {
		return "Operation(<nil>)"
	}
	return fmt.Sprintf("Operation(%s)", t.Alias)
}

func (t *Operation) Operation() ajson.Operation {
	if t == nil {
		return nil
	}
	return ajson.Operations[t.Alias]
}

func (t *Operation) Priority() uint8 {
	if t == nil {
		return 0
	}
	return ajson.OperationsPriority[t.Alias]
}

func (t *Operation) IsRight() bool {
	if t == nil {
		return false
	}
	return ajson.RightOp[t.Alias]
}

func (t *Operation) Parent() Token {
	if t == nil {
		return nil
	}
	return t.parent
}
