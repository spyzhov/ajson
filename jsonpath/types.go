package jsonpath

import "fmt"

type JSONPath struct {
	Tokens []Token
}

type Token interface {
	fmt.Stringer
}

type (
	Tokens []string
)
