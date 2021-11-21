package tokens

import "fmt"

// Token fixme
type Token interface {
	fmt.Stringer
	Type() string
	Token() string
}

// Path fixme
// Part of JSONPath
type Path interface {
	Token
	Path() string
}
