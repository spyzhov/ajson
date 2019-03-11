package ajson

import "fmt"

type Error struct {
	Index int
	Char  byte
}

func errorSymbol(b byte, index int) error {
	return &Error{index, b}
}

func (err *Error) Error() string {
	return fmt.Sprintf("wrong symbol '%s' at %d", []byte{err.Char}, err.Index)
}
