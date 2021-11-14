package ajson

import (
	. "github.com/spyzhov/ajson/v1/internal"
)

/*
	The action codes.
	Copy from `internal/state.go:144`
*/
const (
	cl States = -2 /* colon           */
	cm States = -3 /* comma           */
	//qt States = -4 /* quote           */
	bo States = -5 /* bracket open    */
	co States = -6 /* curly br. open  */
	bc States = -7 /* bracket close   */
	cc States = -8 /* curly br. close */
	ec States = -9 /* curly br. empty */
)

// Unmarshal parses the JSON-encoded data and return the root node of struct.
//
// Doesn't calculate values, just type of stored value. It will store link to the data, on all life long.
func Unmarshal(data []byte) (root *Node, err error) {
	buf := NewBuffer(data)
	var (
		state   States
		key     *string
		current *Node
	)

	_, err = buf.First()
	if err != nil {
		return nil, buf.ErrorEOF()
	}

	for {
		state = buf.GetState()
		if state == __ {
			return nil, buf.ErrorSymbol()
		}

		if state >= GO {
			// region Change State
			switch buf.State {
			case ST:
				if current != nil && current.IsObject() && key == nil {
					// Detected: Key
					key, err = getString(buf)
					buf.State = CO
				} else {
					// Detected: String
					current, err = newNode(current, buf, String, &key)
					if err != nil {
						break
					}
					err = buf.AsString(BQuotes, false)
					current.borders[1] = buf.Index + 1
					buf.State = OK
					if current.parent != nil {
						current = current.parent
					}
				}
			case MI, ZE, IN:
				current, err = newNode(current, buf, Numeric, &key)
				if err != nil {
					break
				}
				err = buf.AsNumeric(false)
				current.borders[1] = buf.Index
				buf.Index -= 1
				buf.State = OK
				if current.parent != nil {
					current = current.parent
				}
			case T1, F1:
				current, err = newNode(current, buf, Bool, &key)
				if err != nil {
					break
				}
				if buf.State == T1 {
					err = buf.AsTrue()
				} else {
					err = buf.AsFalse()
				}
				current.borders[1] = buf.Index + 1
				buf.State = OK
				if current.parent != nil {
					current = current.parent
				}
			case N1:
				current, err = newNode(current, buf, Null, &key)
				if err != nil {
					break
				}
				err = buf.AsNull()
				current.borders[1] = buf.Index + 1
				buf.State = OK
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
				buf.State = OK
			case bc: /* ] */
				if current != nil && current.IsArray() && !current.ready() {
					current.borders[1] = buf.Index + 1
					if current.parent != nil {
						current = current.parent
					}
				} else {
					err = buf.ErrorSymbol()
				}
				buf.State = OK
			case co: /* { */
				current, err = newNode(current, buf, Object, &key)
				buf.State = OB
			case bo: /* [ */
				current, err = newNode(current, buf, Array, &key)
				buf.State = AR
			case cm: /* , */
				if current == nil {
					return nil, buf.ErrorSymbol()
				}
				if current.IsObject() {
					buf.State = KE
				} else if current.IsArray() {
					buf.State = VA
				} else {
					err = buf.ErrorSymbol()
				}
			case cl: /* : */
				if current == nil || !current.IsObject() || key == nil {
					err = buf.ErrorSymbol()
				} else {
					buf.State = VA
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
		if _, err = buf.First(); err != nil {
			err = nil
			break
		}
	}

	if current == nil || buf.State != OK {
		err = buf.ErrorEOF()
	} else {
		root = current.root()
		if !root.ready() {
			err = buf.ErrorEOF()
			root = nil
		}
	}

	return
}

// UnmarshalSafe do the same thing as Unmarshal, but copy data to the local variable, to make it editable.
func UnmarshalSafe(data []byte) (root *Node, err error) {
	var safe []byte
	safe = append(safe, data...)
	return Unmarshal(safe)
}

// Must returns a Node if there was no error. Else - panic with error as the value.
func Must(root *Node, err error) *Node {
	if err != nil {
		panic(err)
	}
	return root
}

func getString(b *Buffer) (*string, error) {
	start := b.Index
	err := b.AsString(BQuotes, false)
	if err != nil {
		return nil, err
	}
	value, ok := Unquote(b.Bytes[start:b.Index+1], BQuotes)
	if !ok {
		return nil, NewErrorSymbol(b)
	}
	return &value, nil
}
