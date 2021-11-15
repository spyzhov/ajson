package tokens

import (
	"fmt"
	"io"
	"strings"

	"github.com/spyzhov/ajson/v1"
)

// RPN is a `Reverse Polish notation`
type RPN []Token

const (
	rpnErrFmt = "RPN builder found an error starting from the %d next to %q: %w"
)

func CompileRPN(token string) (result RPN, err error) {
	return NewRPN([]byte(token))
}

func NewRPN(token []byte) (result RPN, err error) {
	return newRPN(ajson.NewBuffer(token))
}

func newRPN(b *ajson.Buffer) (result RPN, err error) {
	var (
		c             byte
		start         int
		eof           error
		temp          string
		current       string
		found         bool
		foundVariable bool
		stack         = make([]Token, 0)
		pop           = func() Token {
			last := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			return last
		}
	)
	for {
		b.Reset()
		c, eof = b.FirstNonSpace()
		if eof != nil {
			break
		}
		switch true {
		// operations
		case ajson.OperationsChar[c]:
			start = b.Index
			if foundVariable {
				var cToken *Operation
				foundVariable = false
				current = ""

				// Read the complete operation into the variable `current`: `+`, `!=`, `<=>`
				for operation, _ := range ajson.Operations {
					if bytes, ok := b.Slice(len(operation) - 1); ok == nil {
						if string(c)+string(bytes) == operation {
							current = operation
						}
					}
				}
				if current == "" {
					return nil, fmt.Errorf(rpnErrFmt, start, string(c), ErrUnknownToken)
				}

				// Create and validate a new operation instance.
				if cToken, err = NewOperation(current); err != nil {
					return nil, fmt.Errorf(rpnErrFmt, start, current, err)
				}
			stackLoop:
				for len(stack) > 0 {
					found = false
					switch last := stack[len(stack)-1].(type) {
					case *Function:
						result = append(result, pop())
					case *Operation:
						if (last.Priority() > cToken.Priority()) || (last.Priority() == cToken.Priority() && !last.IsRight()) {
							result = append(result, pop())
						} else {
							break stackLoop
						}
					default:
						break stackLoop
					}
				}
				stack = append(stack, cToken)
				break
			}
			if c != ajson.BMinus && c != ajson.BPlus {
				return nil, fmt.Errorf(rpnErrFmt, start, string(c), ErrUnknownToken)
			}
			fallthrough // for numbers like `-1e6`, `+100500`
		// numbers
		case (c >= '0' && c <= '9') || c == '.':
			var cToken *Number
			foundVariable = true
			start = b.Index
			err = b.AsNumeric(true)
			if err != nil {
				return nil, fmt.Errorf(rpnErrFmt, start, string(c), err)
			}
			if cToken, err = NewNumber(string(b.Bytes[start:b.Index])); err != nil {
				return nil, fmt.Errorf(rpnErrFmt, start, string(c), err)
			}
			result = append(result, cToken)
			b.Index--
		// string, with double quotes
		case c == ajson.BQuotes:
			fallthrough
		// string, with a single quote
		case c == ajson.BQuote:
			var cToken *String
			foundVariable = true
			start = b.Index
			err = b.AsString(c, true)
			if err != nil {
				return nil, b.ErrorEOF()
			}
			if cToken, err = NewString(b.Bytes[start : b.Index+1] /* with quotes */); err != nil {
				return nil, fmt.Errorf(rpnErrFmt, start, string(c), err)
			}
			result = append(result, cToken)
		// todo ...
		case c == ajson.BDollar || c == ajson.BAt: // foundVariable : like @.length , $.expensive, etc.
			foundVariable = true
			start = b.Index
			err = b.Token()
			if err != nil {
				if err != io.EOF {
					return nil, err
				}
			}
			current = string(b.Bytes[start:b.Index])
			result = append(result, current)
			if err != nil {
				err = nil
			} else {
				b.Index--
			}
		case c == ajson.BParenthesesL: // (
			foundVariable = false
			current = string(c)
			stack = append(stack, current)
		case c == ajson.BParenthesesR: // )
			foundVariable = true
			found = false
			for len(stack) > 0 {
				temp = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if temp == "(" {
					found = true
					break
				}
				result = append(result, temp)
			}
			if !found { // have no BParenthesesL
				return nil, ajson.NewErrorRequest("formula has no left parentheses")
			}
		default: // prefix functions or etc.
			start = b.Index
			foundVariable = true
			for ; b.Index < b.Length; b.Index++ {
				c = b.Bytes[b.Index]
				if c == ajson.BParenthesesL { // function detection, example: sin(...), round(...), etc.
					foundVariable = false
					break
				}
				if c < 'A' || c > 'z' { // fixme
					if !(c >= '0' && c <= '9') && c != '_' { // constants detection, example: true, false, null, PI, e, etc.
						break
					}
				}
			}
			current = strings.ToLower(string(b.Bytes[start:b.Index]))
			b.Index--
			if !foundVariable {
				if _, found = ajson.Functions[current]; !found {
					return nil, ajson.NewErrorRequest("wrong formula, '%s' is not a function", current)
				}
				stack = append(stack, current)
			} else {
				if _, found = ajson.Constants[current]; !found {
					return nil, ajson.NewErrorRequest("wrong formula, '%s' is not a constant", current)
				}
				result = append(result, current)
			}
		}
		err = b.Step()
		if err != nil {
			break
		}
	}
	if err == io.EOF {
		err = nil // only io.EOF can be here
	}

	for len(stack) > 0 {
		temp = stack[len(stack)-1]
		_, ok := ajson.Functions[temp]
		if ajson.OperationsPriority[temp] == 0 && !ok { // operations only
			return nil, ajson.NewErrorRequest("wrong formula, '%s' is not an operation or function", temp)
		}
		result = append(result, temp)
		stack = stack[:len(stack)-1]
	}

	if len(result) == 0 {
		return nil, b.ErrorEOF()
	}

	return
}

func (r RPN) String() string {
	if r == nil {
		return "<nil>"
	}
	parts := make([]string, 0, len(r))
	for _, token := range r {
		parts = append(parts, token.String())
	}
	return fmt.Sprintf("RPN(%s)", strings.Join(parts, " "))
}
