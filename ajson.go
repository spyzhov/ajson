package ajson

import "io"

func Unmarshal(body []byte, clone bool) (root *Node, err error) {
	buf := newBuffer(body, clone)
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

		if !isCreatable(b, current, last) {
			return nil, errorSymbol(buf)
		}
		switch true {
		case b == bracesL:
			// Detected: Object [begin]
		case b == bracesR:
			// Detected: Object [end]
		case b == bracketL:
			// Detected: Array [begin]
			current = newNode(current, buf, Array, key)
			err = buf.step()
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
			// Detected: String
			current = newNode(current, buf, String, key)
			err = buf.string()
			current.borders[1] = buf.index + 1
			if err == nil {
				err = buf.step()
			}
			found = true
			current = previous(current)
		case (b >= '0' && b <= '9') || b == '.' || b == '+' || b == '-' || b == 'e' || b == 'E':
			// Detected: Numeric
			current = newNode(current, buf, Numeric, key)
			err = buf.numeric()
			current.borders[1] = buf.index
			found = true
			current = previous(current)
		case b == 'n' || b == 'N':
			// mb: Null
			current = newNode(current, buf, Null, key)
			err = buf.null()
			current.borders[1] = buf.index + 1
			if err == nil {
				err = buf.step()
			}
			found = true
			current = previous(current)
		case b == 't' || b == 'T' || b == 'f' || b == 'F':
			// mb: Bool
			current = newNode(current, buf, Bool, key)
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
			if last == coma || last == bracesL || last == bracketL || !found {
				return nil, errorSymbol(buf)
			} else {
				found = false
				err = buf.step()
			}
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

func isCreatable(b byte, current *Node, last byte) bool {
	if b == bracketL || b == bracesL || b == quotes || (b >= '0' && b <= '9') || b == '.' || b == '+' || b == '-' || b == 'e' || b == 'E' || b == 't' || b == 'T' || b == 'f' || b == 'F' || b == 'n' || b == 'N' {
		return current == nil || (current.isContainer() && !current.ready() && (len(current.children) == 0 || last == coma))
	}
	return true
}
