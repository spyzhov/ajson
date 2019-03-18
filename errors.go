package ajson

import "fmt"

//Error is common struct to provide internal errors
type Error struct {
	Type  ErrorType
	Index int
	Char  byte
}

//ErrorType is container for reflection type of error
type ErrorType int

const (
	//WrongSymbol means that system found symbol than not allowed to be
	WrongSymbol ErrorType = iota
	//UnexpectedEOF means that data ended, leaving the node undone
	UnexpectedEOF
	//WrongType means that wrong type requested
	WrongType
	//WrongRequest means that wrong range requested
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

//Error interface implementation
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
