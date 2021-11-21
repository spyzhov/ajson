package internal

import "fmt"

var (
	ErrWrongSymbol   = fmt.Errorf("wrong symbol")
	ErrUnexpectedEof = fmt.Errorf("unexpected end of file")
)
