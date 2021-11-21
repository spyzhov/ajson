package tokens

import "fmt"

var (
	ErrUnknownToken     = fmt.Errorf("unknown token")
	ErrIncorrectFormula = fmt.Errorf("incorrect formula format")
)
