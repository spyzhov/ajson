package ajson

import "io"

func Unmarshal(body []byte, clone bool) (root *Node, err error) {
	buf := newBuffer(body, clone)
	var (
		b       byte
		index   *int
		key     []byte
		current *Node
	)
	// main loop: detect all parts of json struct
base:
	for {
		// detect type of element
		b, err = buf.first()
		if err != nil {
			break
		}
		switch true {
		case b == bracesL:
			// Detected: Object
		case b == bracketL:
			// Detected: Array
		case b == quotes:
			// Detected: String
		case b >= '0' || b <= '9' || b == '+' || b == '-' || b == 'e' || b == 'E':
			// Detected: Numeric
			current = newNode(current, buf, Numeric, key, index)
			err = buf.numeric()
			if err == io.EOF {
				current.borders[1] = buf.index
			}
			if err != nil {
				break base
			}
		case b == 'n' || b == 'N':
			// mb: Null
		case b == 't' || b == 'T' || b == 'f' || b == 'F':
			// mb: Bool
		default:
			return nil, errorSymbol(b, buf.index)
		}
	}
	// outer
	if err == io.EOF {
		err = nil
		root = current
	}

	return
}
