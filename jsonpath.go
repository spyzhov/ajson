package ajson

import (
	"io"
	"math"
	"strconv"
	"strings"
)

// JSONPath returns slice of founded elements in current JSON data, by it's JSONPath.
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
func JSONPath(data []byte, path string) (result []*Node, err error) {
	commands, err := ParseJSONPath(path)
	if err != nil {
		return nil, err
	}
	root, err := Unmarshal(data)
	if err != nil {
		return nil, err
	}
	result = make([]*Node, 0)

	var (
		temporary      []*Node
		keys           []string
		from, to, step int
		c              byte
		key            string
	)
	for i, cmd := range commands {
		switch {
		case cmd == "$": // root element
			if i == 0 {
				result = append(result, root)
			}
		case cmd == "..": // recursive descent
			temporary = make([]*Node, 0)
			for _, element := range result {
				temporary = append(temporary, recursiveChildren(element)...)
			}
			result = append(result, temporary...)
		case cmd == "*": // wildcard
			temporary = make([]*Node, 0)
			for _, element := range result {
				temporary = append(temporary, element.inheritors()...)
			}
			result = temporary
		case strings.Contains(cmd, ":"): // array slice operator
			keys = strings.Split(cmd, ":")
			if len(keys) > 3 {
				return nil, errorRequest("slice must contains no more than 2 colons, got '%s'", cmd)
			}
			if keys[0] == "" {
				from = 0
			} else {
				from, err = strconv.Atoi(keys[0])
				if err != nil {
					return nil, errorRequest("start of slice must be number, got '%s'", keys[0])
				}
			}
			if keys[1] == "" {
				to = math.MaxInt64
			} else {
				to, err = strconv.Atoi(keys[1])
				if err != nil {
					return nil, errorRequest("stop of slice must be number, got '%s'", keys[1])
				}
			}
			step = 1
			if len(keys) == 3 {
				if keys[2] != "" {
					step, err = strconv.Atoi(keys[2])
					if err != nil {
						return nil, errorRequest("step of slice must be number, got '%s'", keys[2])
					}
				}
			}

			temporary = make([]*Node, 0)
			for _, element := range result {
				if element.IsArray() {
					for i := from; i < to; i += step {
						value, ok := element.children[strconv.Itoa(i)]
						if ok {
							temporary = append(temporary, value)
						} else {
							break
						}
					}
				}
			}
			result = temporary
		case strings.HasPrefix(cmd, "?("): // applies a filter (script) expression
		//todo
		//$..[?(@.price == 19.95 && @.color == 'red')].color
		case strings.HasPrefix(cmd, "("): // script expression, using the underlying script engine
		//todo
		default: // try to get by key & Union
			buf := newBuffer([]byte(cmd))
			keys = make([]string, 0)
			for {
				c, err = buf.first()
				if err != nil {
					return nil, errorRequest("blank request")
				}
				if c == coma {
					return nil, errorRequest("wrong request: %s", cmd)
				}
				from = buf.index
				err = buf.token()
				if err != nil && err != io.EOF {
					return nil, errorRequest("wrong request: %s", cmd)
				}
				key = string(buf.data[from:buf.index])
				if len(key) > 2 && key[0] == quote && key[len(key)-1] == quote { // string
					key = key[1 : len(key)-1]
				}
				keys = append(keys, key)
				c, err = buf.first()
				if err != nil {
					err = nil
					break
				}
				if c != coma {
					return nil, errorRequest("wrong request: %s", cmd)
				}
				err = buf.step()
				if err != nil {
					return nil, errorRequest("wrong request: %s", cmd)
				}
			}

			temporary = make([]*Node, 0)
			for _, key = range keys {
				for _, element := range result {
					if element.isContainer() {
						value, ok := element.children[key]
						if ok {
							temporary = append(temporary, value)
						}
					}
				}
			}
			result = temporary
		}
	}
	return
}

//Paths returns calculated paths of underlying nodes
func Paths(array []*Node) []string {
	result := make([]string, 0, len(array))
	for _, element := range array {
		result = append(result, element.Path())
	}
	return result
}

func recursiveChildren(node *Node) (result []*Node) {
	if node.isContainer() {
		for _, element := range node.inheritors() {
			if element.isContainer() {
				result = append(result, element)
			}
		}
	}
	temp := make([]*Node, 0, len(result))
	temp = append(temp, result...)
	for _, element := range result {
		temp = append(temp, recursiveChildren(element)...)
	}
	return temp
}

// ParseJSONPath will parse current path and return all commands tobe run.
// Example:
//
//	result, _ := ParseJSONPath("$.store.book[?(@.price < 10)].title")
//	result == []string{"$", "store", "book", "?(@.price < 10)", "title"}
//
func ParseJSONPath(path string) (result []string, err error) {
	buf := newBuffer([]byte(path))
	result = make([]string, 0)
	var (
		c           byte
		start, stop int
		childEnd    = map[byte]bool{dot: true, bracketL: true}
		str         bool
	)
	for {
		c, err = buf.current()
		if err != nil {
			break
		}
	parseSwitch:
		switch true {
		case c == dollar || c == at:
			result = append(result, string(c))
		case c == dot:
			start = buf.index
			c, err = buf.next()
			if err == io.EOF {
				err = nil
				break
			}
			if err != nil {
				break
			}
			if c == dot {
				result = append(result, "..")
				buf.index--
				break
			}
			err = buf.skipAny(childEnd)
			stop = buf.index
			if err == io.EOF {
				err = nil
				stop = buf.length
			} else {
				buf.index--
			}
			if err != nil {
				break
			}
			if start+1 < stop {
				result = append(result, string(buf.data[start+1:stop]))
			}
		case c == bracketL:
			_, err = buf.next()
			if err != nil {
				return nil, buf.errorEOF()
			}
			start = buf.index
			for ; buf.index < buf.length; buf.index++ {
				c = buf.data[buf.index]
				if c == quote {
					if str {
						str = buf.backslash()
					} else {
						str = true
					}
				} else if c == bracketR {
					if !str {
						result = append(result, string(buf.data[start:buf.index]))
						break parseSwitch
					}
				}
			}
			return nil, buf.errorEOF()
		default:
			return nil, buf.errorSymbol()
		}
		err = buf.step()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
	}
	return
}
