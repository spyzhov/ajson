package jsonpath

import (
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
			case internal.ST:
				switch current.(type) {
				case *tokens.Object:
					// Detected: Key
					start := buf.Index
					err := buf.AsString(ainternal.BQuotes)
					if err != nil {
						return nil, err
					}
					value, ok := ainternal.Unquote(buf.BytesFrom(start, buf.Index+1), ainternal.BQuotes)
					if !ok {
						return nil, buf.ErrorSymbol()
					}

					current, err = tokens.NewObjectElement(nil, nil, current)
					if err != nil {
						return nil, err
					}
					current.(*tokens.ObjectElement).Key = tokens.NewString(value, current)
					if buf.State == qt {
						buf.State = internal.CO
					}
				default: // fixme:here
					// Detected: String
					current, err = newNode(current, buf, String, &key)
					if err != nil {
						break
					}
					err = buf.AsString(BQuotes, false)
					current.borders[1] = buf.Index + 1
					buf.State = internal.OK
					if current.Parent() != nil {
						current = current.Parent()
					}
				}
			case internal.MI, internal.ZE, internal.IN:
				current, err = newNode(current, buf, Numeric, &key)
				if err != nil {
					break
				}
				err = buf.AsNumeric(false)
				current.borders[1] = buf.Index
				buf.Index -= 1
				buf.State = internal.OK
				if current.parent != nil {
					current = current.parent
				}
			case internal.T1, internal.F1:
				current, err = newNode(current, buf, Bool, &key)
				if err != nil {
					break
				}
				if buf.State == internal.T1 {
					err = buf.AsTrue()
				} else {
					err = buf.AsFalse()
				}
				current.borders[1] = buf.Index + 1
				buf.State = internal.OK
				if current.parent != nil {
					current = current.parent
				}
			case internal.N1:
				current, err = newNode(current, buf, Null, &key)
				if err != nil {
					break
				}
				err = buf.AsNull()
				current.borders[1] = buf.Index + 1
				buf.State = internal.OK
				if current.parent != nil {
					current = current.parent
				}
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
