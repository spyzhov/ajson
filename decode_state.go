package ajson

import (
	. "github.com/spyzhov/ajson/internal"
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
	buf := newBuffer(data)
	var (
		state   States
		key     *string
		current *Node
	)

	_, err = buf.first()
	if err != nil {
		return nil, buf.errorEOF()
	}

	for {
		state = buf.getState()
		if state == __ {
			return nil, buf.errorSymbol()
		}

		if state >= GO {
			// region Change State
			switch buf.state {
			case ST:
				if current != nil && current.IsObject() && key == nil {
					// Detected: Key
					key, err = getString(buf)
					buf.state = CO
				} else {
					// Detected: String
					current, err = newNode(current, buf, String, &key)
					if err != nil {
						break
					}
					err = buf.string(quotes, false)
					current.borders[1] = buf.index + 1
					buf.state = OK
					if current.parent != nil {
						current = current.parent
					}
				}
			case MI, ZE, IN:
				current, err = newNode(current, buf, Numeric, &key)
				if err != nil {
					break
				}
				err = buf.numeric(false)
				current.borders[1] = buf.index
				buf.index -= 1
				buf.state = OK
				if current.parent != nil {
					current = current.parent
				}
			case T1, F1:
				current, err = newNode(current, buf, Bool, &key)
				if err != nil {
					break
				}
				if buf.state == T1 {
					err = buf.true()
				} else {
					err = buf.false()
				}
				current.borders[1] = buf.index + 1
				buf.state = OK
				if current.parent != nil {
					current = current.parent
				}
			case N1:
				current, err = newNode(current, buf, Null, &key)
				if err != nil {
					break
				}
				err = buf.null()
				current.borders[1] = buf.index + 1
				buf.state = OK
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
					err = buf.errorSymbol()
				}
				fallthrough
			case cc: /* } */
				if current != nil && current.IsObject() && !current.ready() {
					current.borders[1] = buf.index + 1
					if current.parent != nil {
						current = current.parent
					}
				} else {
					err = buf.errorSymbol()
				}
				buf.state = OK
			case bc: /* ] */
				if current != nil && current.IsArray() && !current.ready() {
					current.borders[1] = buf.index + 1
					if current.parent != nil {
						current = current.parent
					}
				} else {
					err = buf.errorSymbol()
				}
				buf.state = OK
			case co: /* { */
				current, err = newNode(current, buf, Object, &key)
				buf.state = OB
			case bo: /* [ */
				current, err = newNode(current, buf, Array, &key)
				buf.state = AR
			case cm: /* , */
				if current == nil {
					return nil, buf.errorSymbol()
				}
				if current.IsObject() {
					buf.state = KE
				} else if current.IsArray() {
					buf.state = VA
				} else {
					err = buf.errorSymbol()
				}
			case cl: /* : */
				if current == nil || !current.IsObject() || key == nil {
					err = buf.errorSymbol()
				} else {
					buf.state = VA
				}
			default: /* syntax error */
				err = buf.errorSymbol()
			}
			// endregion Action
		}
		if err != nil {
			return
		}
		if buf.step() != nil {
			break
		}
		if _, err = buf.first(); err != nil {
			err = nil
			break
		}
	}

	if current == nil || buf.state != OK {
		err = buf.errorEOF()
	} else {
		root = current.root()
		if !root.ready() {
			err = buf.errorEOF()
			root = nil
		}
	}

	return
}
