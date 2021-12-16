package jerrors

import "fmt"

var (
	ErrUnknownToken        = fmt.Errorf("unknown token")
	ErrWrongSymbol         = fmt.Errorf("wrong symbol")
	ErrIncorrectFormula    = fmt.Errorf("incorrect formula format")
	ErrIncorrectJSONPath   = fmt.Errorf("incorrect JSONPath format")
	ErrUnexpectedEOF       = fmt.Errorf("unexpected end of file or string")
	ErrUnfinishedToken     = fmt.Errorf("unfinished token")
	ErrUnexpectedStatement = fmt.Errorf("unexpected statement")
	ErrBlankRequest        = fmt.Errorf("blank request")
)
