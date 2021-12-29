package tokens

import (
	"fmt"

	"github.com/spyzhov/ajson/v1/jsonpath/internal"
)

// Token fixme
type Token interface {
	fmt.Stringer
	Type() string
	Token() string
	Parent() Token
	SetParent(token Token)
}

// Path fixme
// Part of JSONPath
type Path interface {
	Token
	Path() string
}

// Container fixme
type Container interface {
	Token
	Append(Token) error
	IsEmpty() bool
	GetState(internal.State) internal.State
}
