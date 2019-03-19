package ajson

import "io"

// JsonPath returns slice of founded elements in current JSON data, by it's JSONPath.
//
// JSONPath expressions:
//
//	$	the root object/element
//	@	the current object/element
//	. or []	child operator
//	..	recursive descent. JSONPath borrows this syntax from E4X.
//	*	wildcard. All objects/elements regardless their names.
//	[]	subscript operator. XPath uses it to iterate over element collections and for predicates. In Javascript and JSON it is the native array operator.
//	[,]	Union operator in XPath results in a combination of node sets. JSONPath allows alternate names or array indices as a set.
//	[start:end:step]	array slice operator borrowed from ES4.
//	?()	applies a filter (script) expression.
//	()	script expression, using the underlying script engine.
func JsonPath(data []byte, path string) (result []*Node, err error) {
	buf := newBuffer([]byte(path))
	root, err := Unmarshal(data)
	if err != nil {
		return nil, err
	}
	result = make([]*Node, 0)
	// validate for first symbol
	b, err := buf.current()
	if err != nil {
		return nil, errorEOF(buf)
	}
	if b != dollar {
		return nil, errorSymbol(buf)
	}
	var (
		//last  byte
		start     int
		found     bool
		childEnd  = map[byte]bool{dot: true, bracketL: true}
		temporary []*Node
	)
	for {
		b, err := buf.current()
		if err != nil {
			break
		}
		switch true {
		case b == dollar:
			result = append(result, root)
			err = buf.step()
			found = true
		case b == dot: // child operator or recursive descent
			err = buf.step()
			if err != nil {
				break
			}
			start = buf.index
			err = buf.skipAny(childEnd)
			if err == io.EOF {
				err = nil
			}
			if err != nil {
				break
			}

			if buf.index-start == 0 { // recursive descent '..'
				temporary = make([]*Node, 0)
				for _, element := range result {
					temporary = append(temporary, recursiveChildren(element)...)
				}
				result = append(result, temporary...)
			} else if buf.index-start == 1 && buf.data[start] == asterisk { // todo:child
				temporary = make([]*Node, 0)
				for _, element := range result {
					temporary = append(temporary, element.inheritors()...)
				}
				result = temporary
			}
			found = true
		}
		//last = b
		if err != nil {
			break
		}
	}

	if err == io.EOF {
		if found {
			err = nil
		} else {
			err = errorEOF(buf)
		}
	}
	return
}

func recursiveChildren(node *Node) (result []*Node) {
	if node.isContainer() {
		for _, element := range node.inheritors() {
			if element.isContainer() {
				result = append(result, element)
			}
		}
	}
	for _, element := range result {
		result = append(result, recursiveChildren(element)...)
	}
	return
}
