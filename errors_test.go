package ajson

import "testing"

func TestError_Error(t *testing.T) {
	tests := []struct {
		name    string
		_type   ErrorType
		message string
	}{
		{name: "WrongSymbol", _type: WrongSymbol, message: "wrong symbol 'S' at 10"},
		{name: "UnexpectedEOF", _type: UnexpectedEOF, message: "unexpected end of file"},
		{name: "WrongType", _type: WrongType, message: "wrong type of Node"},
		{name: "WrongRequest", _type: WrongRequest, message: "wrong request: example error"},
		{name: "unknown", _type: -666, message: "unknown error: 'S' at 10"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := &Error{
				Type:    test._type,
				Index:   10,
				Char:    'S',
				Message: "example error",
			}
			if result.Error() != test.message {
				t.Errorf("Wrong error message: %s", result.Error())
			}
		})
	}
}
