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
				// todo: here
			}
			// endregion Change State
		} else {
			// region Action
			switch state {
			case ec: /* empty } */
				if key != nil {
					err = buf.ErrorSymbol()
				}
				fallthrough
			case cc: /* } */
				if current != nil && current.IsObject() && !current.ready() {
					current.borders[1] = buf.Index + 1
					if current.parent != nil {
						current = current.parent
					}
				} else {
					err = buf.ErrorSymbol()
				}
				buf.State = internal.OK
			case bc: /* ] */
				if current != nil && current.IsArray() && !current.ready() {
					current.borders[1] = buf.Index + 1
					if current.parent != nil {
						current = current.parent
					}
				} else {
					err = buf.ErrorSymbol()
				}
				buf.State = internal.OK
			case co: /* { */
				current, err = newNode(current, buf, Object, &key)
				buf.State = internal.OB
			case bo: /* [ */
				current, err = newNode(current, buf, Array, &key)
				buf.State = internal.AR
			case cm: /* , */
				if current == nil {
					return nil, buf.ErrorSymbol()
				}
				if current.IsObject() {
					buf.State = internal.KE
				} else if current.IsArray() {
					buf.State = internal.VA
				} else {
					err = buf.ErrorSymbol()
				}
			case cl: /* : */
				if current == nil || !current.IsObject() || key == nil {
					err = buf.ErrorSymbol()
				} else {
					buf.State = internal.VA
				}
			default: /* syntax error */
				err = buf.ErrorSymbol()
			}
			// endregion Action
		}
		if err != nil {
			return
		}
		if buf.Step() != nil {
			break
		}
		if _, err = buf.FirstNonSpace(); err != nil {
			err = nil
			break
		}
	}

	if current == nil || buf.State != internal.OK {
		err = buf.ErrorEOF()
	} else {
		root = current.root()
		if !root.ready() {
			err = buf.ErrorEOF()
			root = nil
		}
	}

	// todo
	panic("not implemented")
	return nil, nil
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
