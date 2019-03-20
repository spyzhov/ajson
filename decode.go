//Package ajson implements decoding of JSON as defined in RFC 7159 without predefined mapping to a struct of golang.
//All JSON structs reflects to a custom struct of Node, witch can be presented by it type and value.
package ajson

import "io"

//UnmarshalSafe do the same thing as Unmarshal, but copy data to the local variable, to make it editable.
func UnmarshalSafe(data []byte) (root *Node, err error) {
	var safe []byte
	safe = append(safe, data...)
	return Unmarshal(safe)
}

//Unmarshal parses the JSON-encoded data and return the root node of struct.
//
//Doesn't calculate values, just type of stored value. It will store link to the data, on all life long.
func Unmarshal(data []byte) (root *Node, err error) {
	buf := newBuffer(data)
	var (
		last    byte
		b       byte
		found   bool
		key     *string
		current *Node
	)
	// main loop: detect all parts of json struct
	for {
		// detect type of element
		b, err = buf.first()
		if err != nil {
			break
		}

		if !isCreatable(b, current, last, key) {
			return nil, errorSymbol(buf)
		}
		switch true {
		case b == bracesL:
			// Detected: Object [begin]
			current, err = newNode(current, buf, Object, &key)
			if err == nil {
				err = buf.step()
			}
			found = false
		case b == bracesR:
			// Detected: Object [end]
			if last == coma || key != nil || current == nil || !current.IsObject() || current.ready() {
				return nil, errorSymbol(buf)
			}
			current.borders[1] = buf.index + 1
			err = buf.step()
			found = true
			current = previous(current)
		case b == bracketL:
			// Detected: Array [begin]
			current, err = newNode(current, buf, Array, &key)
			if err == nil {
				err = buf.step()
			}
			found = false
		case b == bracketR:
			// Detected: Array [end]
			if last == coma || current == nil || !current.IsArray() || current.ready() {
				return nil, errorSymbol(buf)
			}
			current.borders[1] = buf.index + 1
			err = buf.step()
			found = true
			current = previous(current)
		case b == quotes:
			// Detected: String OR Key
			if current != nil && current.IsObject() && key == nil { // Detected: Key
				key, err = getString(buf)
				if err == nil {
					err = buf.step()
				}
				found = false
			} else { // Detected: String
				current, err = newNode(current, buf, String, &key)
				if err != nil {
					break
				}
				err = buf.string(quotes)
				current.borders[1] = buf.index + 1
				if err == nil {
					err = buf.step()
				}
				found = true
				current = previous(current)
			}
		case (b >= '0' && b <= '9') || b == '.' || b == '+' || b == '-' || b == 'e' || b == 'E':
			// Detected: Numeric
			current, err = newNode(current, buf, Numeric, &key)
			if err != nil {
				break
			}
			err = buf.numeric()
			current.borders[1] = buf.index
			found = true
			current = previous(current)
		case b == 'n' || b == 'N':
			// Detected: Null
			current, err = newNode(current, buf, Null, &key)
			if err != nil {
				break
			}
			err = buf.null()
			current.borders[1] = buf.index + 1
			if err == nil {
				err = buf.step()
			}
			found = true
			current = previous(current)
		case b == 't' || b == 'T' || b == 'f' || b == 'F':
			// Detected: Bool
			current, err = newNode(current, buf, Bool, &key)
			if err != nil {
				break
			}
			if b == 't' || b == 'T' {
				err = buf.true()
			} else {
				err = buf.false()
			}
			current.borders[1] = buf.index + 1
			if err == nil {
				err = buf.step()
			}
			found = true
			current = previous(current)
		case b == coma:
			if last == coma || current == nil || current.Empty() || !found {
				return nil, errorSymbol(buf)
			}
			found = false
			err = buf.step()
		case b == colon:
			if last != quotes || key == nil || found {
				return nil, errorSymbol(buf)
			}
			found = false
			err = buf.step()
		default:
			return nil, errorSymbol(buf)
		}
		if err != nil {
			break
		}
		last = b
	}

	// outer
	if err == io.EOF {
		if current == nil || current.parent != nil || !current.ready() {
			return nil, errorEOF(buf)
		}
		err = nil
		root = current
	}

	return
}

func previous(current *Node) *Node {
	if current.parent != nil {
		return current.parent
	}
	return current
}

func isCreatable(b byte, current *Node, last byte, key *string) bool {
	if b == bracketL || b == bracesL || b == quotes || (b >= '0' && b <= '9') || b == '.' || b == '+' || b == '-' || b == 'e' || b == 'E' || b == 't' || b == 'T' || b == 'f' || b == 'F' || b == 'n' || b == 'N' {
		if current == nil {
			return key == nil
		}
		if key != nil && !current.IsObject() {
			return false
		}
		if current.isContainer() && current.ready() {
			return false
		}
		if current.IsArray() {
			if len(current.children) == 0 {
				return last != coma
			}
			return last == coma
		}
	}
	return true
}

func getString(b *buffer) (*string, error) {
	start := b.index
	err := b.string(quotes)
	if err != nil {
		return nil, err
	}
	value := string(b.data[start+1 : b.index])
	return &value, nil
}
