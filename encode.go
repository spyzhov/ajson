package ajson

import (
	"strconv"

	"github.com/spyzhov/ajson/v1/internal"
)

// Marshal returns slice of bytes, marshaled from current value
func Marshal(node *Node) (result []byte, err error) {
	result = make([]byte, 0)
	var (
		sValue string
		bValue bool
		nValue float64
		oValue []byte
	)

	if node == nil {
		return nil, NewErrorUnparsed()
	} else if node.dirty {
		switch node._type {
		case Null:
			result = append(result, C_null...)
		case Numeric:
			nValue, err = node.GetNumeric()
			if err != nil {
				return nil, err
			}
			result = append(result, strconv.FormatFloat(nValue, 'g', -1, 64)...)
		case String:
			sValue, err = node.GetString()
			if err != nil {
				return nil, err
			}
			result = append(result, BQuotes)
			result = append(result, internal.QuoteString(sValue, true)...)
			result = append(result, BQuotes)
		case Bool:
			bValue, err = node.GetBool()
			if err != nil {
				return nil, err
			} else if bValue {
				result = append(result, C_true...)
			} else {
				result = append(result, C_false...)
			}
		case Array:
			result = append(result, BBracketL)
			for i := 0; i < len(node.children); i++ {
				if i != 0 {
					result = append(result, BComa)
				}
				child, ok := node.children[strconv.Itoa(i)]
				if !ok {
					return nil, NewErrorRequest("wrong length of array")
				}
				oValue, err = Marshal(child)
				if err != nil {
					return nil, err
				}
				result = append(result, oValue...)
			}
			result = append(result, BBracketR)
		case Object:
			result = append(result, BBracesL)
			bValue = false
			for key, child := range node.children {
				if bValue {
					result = append(result, BComa)
				} else {
					bValue = true
				}
				result = append(result, BQuotes)
				result = append(result, internal.QuoteString(key, true)...)
				result = append(result, BQuotes, BColon)
				oValue, err = Marshal(child)
				if err != nil {
					return nil, err
				}
				result = append(result, oValue...)
			}
			result = append(result, BBracesR)
		}
	} else if node.ready() {
		result = append(result, node.Source()...)
	} else {
		return nil, NewErrorUnparsed()
	}

	return
}
