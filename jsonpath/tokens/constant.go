package tokens

import (
	"fmt"
	"strings"

	"github.com/spyzhov/ajson/v1"
)

type Constant struct {
	parent Token
	Alias  string
}

var _ Token = (*Constant)(nil)

func NewConstant(alias string) (*Constant, error) {
	alias = strings.ToLower(alias)
	// fixme: m.b. remove this dependency?
	if _, ok := ajson.Constants[alias]; !ok {
		return nil, fmt.Errorf("constant %q not found", alias)
	}
	return &Constant{
		Alias: alias,
	}, nil
}

func (t *Constant) Type() string {
	return "Constant"
}

func (t *Constant) String() string {
	if t == nil {
		return "<nil>"
	}
	return t.Alias
}

func (t *Constant) Token() string {
	if t == nil {
		return "Constant(<nil>)"
	}
	return fmt.Sprintf("Constant(%s)", t.Alias)
}

func (t *Constant) Parent() Token {
	if t == nil {
		return nil
	}
	return t.parent
}
