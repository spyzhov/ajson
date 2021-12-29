package jsonpath

import (
	"fmt"

	ainternal "github.com/spyzhov/ajson/v1/internal"
	"github.com/spyzhov/ajson/v1/jerrors"
	"github.com/spyzhov/ajson/v1/jsonpath/internal"
	"github.com/spyzhov/ajson/v1/jsonpath/tokens"
)

func New(jsonpath []byte) (*JSONPath, error) {
	buf := ainternal.NewBuffer(jsonpath, internal.StateTransitionTable)

	var (
		state   internal.State
		current tokens.Token
	)

	if _, eof := buf.FirstNonSpace(); eof != nil {
		return nil, jerrors.ErrBlankRequest
	}

	for {
		state = buf.GetNextState()
		if state == __ {
			return nil, buf.ErrorSymbol()
		}

		if state >= internal.GO {
			// region Change State
			switch buf.State {
			// String
			case internal.ST:
				switch token := current.(type) {
				case *tokens.Object:
					// Detected: Key
					start := buf.Index
					err := buf.AsStringBordered(ainternal.BQuotes)
					if err != nil {
						return nil, err
					}
					value, ok := ainternal.Unquote(buf.BytesFrom(start, buf.Index+1), ainternal.BQuotes)
					if !ok {
						return nil, buf.ErrorSymbol()
					}
					current = token.NewObjectElement()
					current.(*tokens.ObjectElement).Key = tokens.NewString(value, current)
					if buf.State == qt {
						buf.State = internal.CO
					}
				case tokens.Container:
					// Detected: String
					start := buf.Index
					err := buf.AsString()
					if err != nil {
						return nil, err
					}
					value, ok := ainternal.Unquote(buf.BytesFrom(start, buf.Index+1), ainternal.BQuotes)
					if !ok {
						return nil, buf.ErrorSymbol()
					}
					err = token.Append(tokens.NewString(value, current))
					if err != nil {
						return nil, err
					}
					buf.State = token.GetState(buf.State)
				default:
					return nil, buf.ErrorIncorrectJSONPath()
				}
			// Number
			case internal.MI, internal.ZE, internal.IN:
				switch token := current.(type) {
				case tokens.Container:
					start := buf.Index
					err := buf.AsNumeric()
					if err != nil {
						return nil, err
					}
					value, err := tokens.NewNumber(buf.StringFrom(start))
					if err != nil {
						return nil, fmt.Errorf("%w started from %d: %s", jerrors.ErrIncorrectJSONPath, start, err)
					}
					err = token.Append(value)
					if err != nil {
						return nil, err
					}
				default:
					return nil, buf.ErrorIncorrectJSONPath()
				}
			// Bool and Null
			case internal.T1, internal.F1, internal.N1:
				switch token := current.(type) {
				case tokens.Container:
					var (
						value *tokens.Constant
						err   error
					)
					switch buf.State {
					case internal.T1:
						value, err = tokens.NewConstant("true")
					case internal.F1:
						value, err = tokens.NewConstant("false")
					case internal.N1:
						value, err = tokens.NewConstant("null")
					}
					if err != nil {
						return nil, err
					}
					err = token.Append(value)
					if err != nil {
						return nil, err
					}

					buf.State = token.GetState(buf.State)
				default:
					return nil, buf.ErrorIncorrectJSONPath()
				}
			case internal.SS:
				switch token := current.(type) {
				case tokens.Container:
					err := token.Append(tokens.NewScript())
					if err != nil {
						return nil, err
					}
				default:
					return nil, buf.ErrorIncorrectJSONPath()
				}
			}
			// endregion Change State
		} else {
			// region Action
			switch state {
			case co: /* { */
				switch token := current.(type) {
				case tokens.Container:
					err := token.Append(tokens.NewObject())
					if err != nil {
						return nil, err
					}
				default:
					return nil, buf.ErrorIncorrectJSONPath()
				}
				buf.State = internal.OB
			case ec: /* empty } */
				fallthrough
			case cc: /* } */
				_, ok := current.(*tokens.Object)
				if !ok {
					return nil, buf.ErrorSymbol()
				}
				current = current.Parent()
				if token, ok := current.(tokens.Container); ok {
					buf.State = token.GetState(buf.State)
				} else {
					return nil, buf.ErrorIncorrectJSONPath()
				}
				// todo: here
			case bo: /* [ */
				switch token := current.(type) {
				case tokens.Container:
					err := token.Append(tokens.NewArray())
					if err != nil {
						return nil, err
					}
				default:
					return nil, buf.ErrorIncorrectJSONPath()
				}
				buf.State = internal.AR
			case bc: /* ] */
				switch token := current.(type) {
				case *tokens.Array:
					break
				case *tokens.Slice:
					if token.Start == nil {
						return nil, buf.ErrorIncorrectJSONPath()
					}
				case *tokens.Child:
					if token.IsEmpty() {
						return nil, buf.ErrorIncorrectJSONPath()
					}
				default:
					return nil, buf.ErrorIncorrectJSONPath()
				}

				// todo: here

				current = current.Parent()
				if token, ok := current.(tokens.Container); ok {
					buf.State = token.GetState(buf.State)
				} else {
					return nil, buf.ErrorIncorrectJSONPath()
				}
			case cm: /* , */
				switch current.(type) {
				case *tokens.Array:
					buf.State = internal.VA
				case *tokens.Object:
					buf.State = internal.KE
				default:
					return nil, buf.ErrorSymbol()
				}
			case cl: /* : */
				switch token := current.(type) {
				case *tokens.ObjectElement:
					if token.Key == nil {
						return nil, buf.ErrorSymbol()
					}
					buf.State = internal.VA
				default:
					return nil, buf.ErrorSymbol()
				}
			default: /* syntax error */
				return nil, buf.ErrorIncorrectJSONPath()
			}
			// endregion Action
		}
		if buf.Step() != nil {
			break
		}
	}

	if current == nil {
		return nil, buf.ErrorEOF()
	}
	if buf.State != internal.OK {
		return nil, buf.ErrorEOF()
	}

	return &JSONPath{
		Root: current,
	}, nil
}

func Compile(jsonpath string) (*JSONPath, error) {
	return New([]byte(jsonpath))
}

func MustCompile(jsonpath string) *JSONPath {
	result, err := Compile(jsonpath)
	if err != nil {
		panic(err)
	}
	return result
}
