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
	WrongType
	WrongRequest
)

func errorSymbol(b *buffer) error {
	return &Error{Type: WrongSymbol, Index: b.index, Char: b.data[b.index]}
}

func errorEOF(b *buffer) error {
	return &Error{Type: UnexpectedEOF, Index: b.index}
}

func errorType() error {
	return &Error{Type: WrongType}
}

func errorRequest() error {
	return &Error{Type: WrongRequest}
}

func (err *Error) Error() string {
	switch err.Type {
	case WrongSymbol:
		return fmt.Sprintf("wrong symbol '%s' at %d", []byte{err.Char}, err.Index)
	case UnexpectedEOF:
		return fmt.Sprintf("unexpected end of file")
	case WrongType:
		return fmt.Sprintf("wrong type of Node")
	case WrongRequest:
		return fmt.Sprintf("wrong request")
	}
	return fmt.Sprintf("unknown error: '%s' at %d", []byte{err.Char}, err.Index)
}
