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
	node, err := Unmarshal(data)
	if err != nil {
		return nil, err
	}
	return deReference(node, commands)
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

func deReference(node *Node, commands []string) (result []*Node, err error) {
	result = make([]*Node, 0)

	var (
		temporary      []*Node
		keys           []string
		from, to, step int
		rfrom, rto     int
		c              byte
		key            string
		ok             bool
		value, temp    *Node
		float          float64
	)
	for i, cmd := range commands {
		switch {
		case cmd == "$": // root element
			if i == 0 {
				result = append(result, root(node))
			}
		case cmd == "@": // current element
			if i == 0 {
				result = append(result, node)
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
		case strings.Contains(cmd, ":"): // fixme:array slice operator
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
					rfrom = from
					if rfrom < 0 {
						rfrom = element.Size() + rfrom
					}
					rto = to
					if rto < 0 {
						rto = element.Size() + rto
					}

					for i := rfrom; i < rto; i += step {
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
		case strings.HasPrefix(cmd, "?(") && strings.HasSuffix(cmd, ")"): // applies a filter (script) expression
			buf := newBuffer([]byte(cmd[2 : len(cmd)-1]))
			rpn, err := buf.rpn()
			if err != nil {
				return nil, errorRequest("wrong request: %s", cmd)
			}
			temporary = make([]*Node, 0)
			for _, element := range result {
				if element.isContainer() {
					for _, temp = range element.inheritors() {
						value, err = eval(temp, rpn, cmd)
						if err != nil {
							return nil, errorRequest("wrong request: %s", cmd)
						}
						if value != nil {
							ok, err = boolean(value)
							if err != nil || !ok {
								continue
							}
							temporary = append(temporary, temp)
						}
					}
				}
			}
			result = temporary
		case strings.HasPrefix(cmd, "(") && strings.HasSuffix(cmd, ")"): // script expression, using the underlying script engine
			buf := newBuffer([]byte(cmd[1 : len(cmd)-1]))
			rpn, err := buf.rpn()
			if err != nil {
				return nil, errorRequest("wrong request: %s", cmd)
			}
			temporary = make([]*Node, 0)
			for _, element := range result {
				if !element.isContainer() {
					continue
				}
				temp, err = eval(element, rpn, cmd)
				if err != nil {
					return nil, errorRequest("wrong request: %s", cmd)
				}
				if temp != nil {
					value = nil
					switch temp.Type() {
					case String:
						key, err = element.GetString()
						if err != nil {
							return nil, errorRequest("wrong type convert: %s", err.Error())
						}
						value, _ = element.children[key]
					case Numeric:
						from, err = temp.getInteger()
						if err == nil { // INTEGER
							if from < 0 {
								key = strconv.Itoa(element.Size() - from)
							} else {
								key = strconv.Itoa(from)
							}
						} else {
							float, err = temp.GetNumeric()
							if err != nil {
								return nil, errorRequest("wrong type convert: %s", err.Error())
							}
							key = strconv.FormatFloat(float, 'g', -1, 64)
						}
						value, _ = element.children[key]
					case Bool:
						ok, err = temp.GetBool()
						if err != nil {
							return nil, errorRequest("wrong type convert: %s", err.Error())
						}
						if ok {
							temporary = append(temporary, element.inheritors()...)
						}
						continue
					}
					if value != nil {
						temporary = append(temporary, value)
					}
				}
			}
			result = temporary
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
					if element.IsArray() {
						if key == "length" {
							value, err = functions["length"](element)
							if err != nil {
								return
							}
							ok = true
						} else {
							from, err = strconv.Atoi(key)
							if err != nil {
								ok = false
								err = nil
							} else {
								if from < 0 {
									key = strconv.Itoa(element.Size() + from)
								}
								value, ok = element.children[key]
							}
						}

					} else if element.IsObject() {
						value, ok = element.children[key]
					}
					if ok {
						temporary = append(temporary, value)
					}
				}
			}
			result = temporary
		}
	}
	return
}

// Evaluate expression `@.price == 19.95 && @.color == 'red'` to the result value i.e. Bool(true), Numeric(3.14), etc.
func eval(node *Node, expression rpn, cmd string) (result *Node, err error) {
	var (
		stack    []*Node
		slice    []*Node
		temp     *Node
		fn       function
		op       operation
		ok       bool
		size     int
		commands []string
		bstr     []byte
	)
	for _, exp := range expression {
		size = len(stack)
		fn, ok = functions[exp]
		if ok {
			if size < 1 {
				return nil, errorRequest("wrong request: %s", cmd)
			}
			stack[size-1], err = fn(stack[size-1])
			if err != nil {
				return
			}
			continue
		}
		op, ok = operations[exp]
		if ok {
			if size < 2 {
				return nil, errorRequest("wrong request: %s", cmd)
			}
			stack[size-2], err = op(stack[size-2], stack[size-1])
			if err != nil {
				return
			}
			stack = stack[:size-1]
			continue
		}
		if len(exp) > 0 {
			if exp[0] == dollar || exp[0] == at {
				commands, err = ParseJSONPath(exp)
				if err != nil {
					return
				}
				slice, err = deReference(node, commands)
				if err != nil {
					return
				}
				if len(slice) == 1 {
					stack = append(stack, slice[0])
				} else { // no data found, or array given
					return nil, nil
				}
			} else {
				bstr = []byte(exp)
				size = len(bstr)
				if size > 2 && bstr[0] == quote && bstr[size-1] == quote {
					bstr[0] = quotes
					bstr[size-1] = quotes
				}
				temp, err = Unmarshal(bstr)
				if err != nil {
					return
				}
				stack = append(stack, temp)
			}
		} else {
			stack = append(stack, varNode(nil, "", String, ""))
		}
	}
	if len(stack) == 1 {
		return stack[0], nil
	}
	if len(stack) == 0 {
		return nil, nil
	}
	return nil, errorRequest("wrong request: %s", cmd)
}

func root(node *Node) (result *Node) {
	for result = node; result.parent != nil; result = result.parent {
	}
	return
}
