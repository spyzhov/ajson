package tokens

import "fmt"

type Token interface {
	fmt.Stringer
	Type() string
	Token() string
}
