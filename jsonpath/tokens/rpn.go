package tokens

import (
	"fmt"
	"strings"

	"github.com/spyzhov/ajson/v1"
	"github.com/spyzhov/ajson/v1/internal"
	"github.com/spyzhov/ajson/v1/jerrors"
)

// RPN is a `Reverse Polish notation`
type RPN struct {
	parent Token
	Tokens []Token
}

var _ Token = (*RPN)(nil)

func NewRPN() *RPN {
	return &RPN{
		Tokens: make([]Token, 0),
	}
}

// fixme: found the way how to stop when part of the other Buffer is given
func newRPN(b *internal.Buffer) (result *RPN, err error) {
	var (
		c          byte
		start      int
		eof        error
		current    string
		found      bool
		isVariable bool
		entrance   = b.Index
		stack      = make([]Token, 0)
		pop        = func() Token {
			last := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			return last
		}
		foundVariable = func() error {
			if isVariable {
				return fmt.Errorf("wrong value position")
			}
			isVariable = true
			return nil
		}
		rpnErrFmt = fmt.Sprintf("RPN(starts from %d) found an error at %%d (%%q): %%w", entrance)
	)
	result = &RPN{Tokens: make([]Token, 0)}
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
			if isVariable {
				var cToken *Operation
				isVariable = false
				current = ""

				// Read the complete operation into the variable `current`: `+`, `!=`, `<=>`
				for operation := range ajson.Operations {
					if bytes, ok := b.Slice(len(operation) - 1); ok == nil {
						if string(c)+string(bytes) == operation {
							current = operation
							_ = b.Move(len(operation) - 1) // error can't occupy here because of b.Slice result
							break
						}
					}
				}
				if current == "" {
					return nil, fmt.Errorf(rpnErrFmt, start, string(c), jerrors.ErrUnknownToken)
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
						result.Tokens = append(result.Tokens, pop())
					case *Operation:
						if (last.Priority() > cToken.Priority()) || (last.Priority() == cToken.Priority() && !last.IsRight()) {
							result.Tokens = append(result.Tokens, pop())
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
			if c != internal.BMinus && c != internal.BPlus {
				return nil, fmt.Errorf(rpnErrFmt, start, string(c), jerrors.ErrUnknownToken)
			}
			fallthrough // for numbers like `-1e6`, `+100500`
		// numbers
		case (c >= '0' && c <= '9') || c == '.':
			var cToken *Number
			if err = foundVariable(); err != nil {
				return nil, fmt.Errorf(rpnErrFmt, b.Index, string(b.Bytes[b.Index]), err)
			}
			start = b.Index
			if cToken, err = newNumber(b, true); err != nil {
				return nil, fmt.Errorf(rpnErrFmt, start, b.StringFrom(start), err)
			}
			result.Tokens = append(result.Tokens, cToken)
			b.Index--
		// string, with double quotes
		case c == internal.BQuotes:
			fallthrough
		// string, with a single quote
		case c == internal.BQuote:
			if err = foundVariable(); err != nil {
				return nil, fmt.Errorf(rpnErrFmt, b.Index, string(b.Bytes[b.Index]), err)
			}
			var cToken *String
			start = b.Index
			if cToken, err = newString(b); err != nil {
				return nil, fmt.Errorf(rpnErrFmt, start, string(c), err)
			}
			result.Tokens = append(result.Tokens, cToken)
		// variable : like @.length , $.expensive, etc.
		case c == internal.BDollar || c == internal.BAt:
			if err = foundVariable(); err != nil {
				return nil, fmt.Errorf(rpnErrFmt, b.Index, string(b.Bytes[b.Index]), err)
			}
			var cToken *JSONPath
			start = b.Index
			cToken, err = newJSONPath(b)
			if err != nil {
				return nil, fmt.Errorf(rpnErrFmt, start, string(b.Bytes[start]), err)
			}
			result.Tokens = append(result.Tokens, cToken)
		// ( : Parenthesis open
		case c == internal.BParenthesesL:
			var cToken *parentheses
			isVariable = false
			cToken, err = newParentheses(c)
			if err != nil {
				return nil, fmt.Errorf(rpnErrFmt, start, string(b.Bytes[start]), err)
			}
			stack = append(stack, cToken)
		// ) : Parenthesis close
		case c == internal.BParenthesesR:
			if err = foundVariable(); err != nil {
				return nil, fmt.Errorf(rpnErrFmt, b.Index, string(b.Bytes[b.Index]), err)
			}
			var cToken Token
			found = false
		stackCycle:
			for len(stack) > 0 {
				cToken = pop()
				switch token := cToken.(type) {
				case *parentheses:
					if token.IsOpen() {
						found = true
						break stackCycle
					}
					err = fmt.Errorf("wrong position for closing parenthesis")
					return nil, fmt.Errorf(rpnErrFmt, start, string(b.Bytes[start]), err)
				default:
					result.Tokens = append(result.Tokens, cToken)
				}
			}
			if !found { // have no BParenthesesL
				err = fmt.Errorf("formula has no open parenthesis")
				return nil, fmt.Errorf(rpnErrFmt, start, string(b.Bytes[start]), err)
			}
		// [ : start the manual array definition
		case c == internal.BBracketL:
			if err = foundVariable(); err != nil {
				return nil, fmt.Errorf(rpnErrFmt, b.Index, string(b.Bytes[b.Index]), err)
			}
			var cToken *Array
			start = b.Index
			cToken, err = newArray(b)
			if err != nil {
				return nil, fmt.Errorf(rpnErrFmt, start, string(b.Bytes[start]), err)
			}
			result.Tokens = append(result.Tokens, cToken)
		// { : start the manual object definition
		case c == internal.BBracesL:
			if err = foundVariable(); err != nil {
				return nil, fmt.Errorf(rpnErrFmt, b.Index, string(b.Bytes[b.Index]), err)
			}
			var cToken *Object
			start = b.Index
			cToken, err = newObject(b)
			if err != nil {
				return nil, fmt.Errorf(rpnErrFmt, start, string(b.Bytes[start]), err)
			}
			result.Tokens = append(result.Tokens, cToken)
		// function or constant
		// todo: add an ability to use letters in operations (# `in`).
		// todo: There could be intersections between operations and function or constant (# `x in [a,b,c]` and `in(1, [1,2,3])`)
		default:
			if err = foundVariable(); err != nil {
				return nil, fmt.Errorf(rpnErrFmt, start, string(b.Bytes[start]), err)
			}
			start = b.Index
			b.Word()
			if start == b.Index {
				return nil, fmt.Errorf(rpnErrFmt, start, string(c), jerrors.ErrUnknownToken)
			}
			current = strings.ToLower(string(b.Bytes[start:b.Index]))
			c, eof = b.FirstNonSpace()
			if c == internal.BParenthesesL {
				// function detection
				var args *Arguments
				args, err = newArguments(b)
				if err != nil {
					return nil, fmt.Errorf(rpnErrFmt, b.Index, string(b.Bytes[b.Index]), err)
				}
				var cToken *Function
				cToken, err = NewFunction(current, args)
				if err != nil {
					return nil, fmt.Errorf(rpnErrFmt, start, current, err)
				}
				result.Tokens = append(result.Tokens, cToken)
			} else {
				// constant found
				var cToken *Constant
				cToken, err = NewConstant(current)
				if err != nil {
					return nil, fmt.Errorf(rpnErrFmt, start, current, err)
				}
				result.Tokens = append(result.Tokens, cToken)
			}
			b.Index-- // fixme: looks weird, test this twice
		}
		eof = b.Step()
		if eof != nil {
			break
		}
	}

	for len(stack) > 0 {
		switch cToken := pop().(type) {
		case *Operation:
			result.Tokens = append(result.Tokens, cToken)
		case *Function: // fixme: validate this twice. check for constants and variables
			result.Tokens = append(result.Tokens, cToken)
		default:
			return nil, fmt.Errorf(rpnErrFmt, entrance, string(b.Bytes[entrance:b.Index]), jerrors.ErrIncorrectFormula)
		}
	}

	if len(result.Tokens) == 0 {
		return nil, fmt.Errorf(rpnErrFmt, entrance, string(b.Bytes[entrance:b.Index]), jerrors.ErrIncorrectFormula)
	}

	return
}

func (t *RPN) Type() string {
	return "RPN"
}

// todo: display it as applied line # ((4 + 6) / sin(pi))
func (t *RPN) String() string {
	if t == nil {
		return "<nil>"
	}
	parts := make([]string, 0, len(t.Tokens))
	for _, token := range t.Tokens {
		parts = append(parts, token.String())
	}
	return fmt.Sprintf("%s", strings.Join(parts, " "))
}

func (t *RPN) Token() string {
	if t == nil {
		return "RPN(<nil>)"
	}
	parts := make([]string, 0, len(t.Tokens))
	for _, token := range t.Tokens {
		parts = append(parts, token.String())
	}
	return fmt.Sprintf("RPN(%s)", strings.Join(parts, ", "))
}

func (t *RPN) Parent() Token {
	if t == nil {
		return nil
	}
	return t.parent
}

func (t *RPN) SetParent(parent Token) {
	if t == nil {
		return
	}
	t.parent = parent
}

// fixme: m.b. implement from ?newRPN?
func (t *RPN) Append(token Token) error {
	t.Tokens = append(t.Tokens, token)
	token.SetParent(t)
	return nil
}

func (t *RPN) IsEmpty() bool {
	return len(t.Tokens) == 0
}

func (t *RPN) GetState(_ internal.State) internal.State {
	return -1 // fixme
}
