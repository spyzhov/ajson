package ajson

import "fmt"

type Error struct {
	Type  ErrorType
	Index int
	Char  byte
}

type ErrorType int

const (
	WrongSymbol ErrorType = iota
	UnexpectedEOF
)

func errorSymbol(b *buffer) error {
	return &Error{Type: WrongSymbol, Index: b.index, Char: b.data[b.index]}
}

func errorEOF(b *buffer) error {
	return &Error{Type: UnexpectedEOF, Index: b.index}
}

func (err *Error) Error() string {
	switch err.Type {
	case WrongSymbol:
		return fmt.Sprintf("wrong symbol '%s' at %d", []byte{err.Char}, err.Index)
	case UnexpectedEOF:
		return fmt.Sprintf("unexpected end of file")
	}
	return fmt.Sprintf("unknown error: '%s' at %d", []byte{err.Char}, err.Index)
}
