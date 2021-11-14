package ajson

import "github.com/spyzhov/ajson/v1/internal"

func CastFloats(left, right *Node) (lnum, rnum float64, err error) {
	lnum, err = left.GetNumeric()
	if err != nil {
		return
	}
	rnum, err = right.GetNumeric()
	return
}

func CastInts(left, right *Node) (lnum, rnum int, err error) {
	lnum, err = left.getInteger()
	if err != nil {
		return
	}
	rnum, err = right.getInteger()
	return
}

func CastBools(left, right *Node) (lnum, rnum bool, err error) {
	lnum, err = left.GetBool()
	if err != nil {
		return
	}
	rnum, err = right.GetBool()
	return
}

func CastStrings(left, right *Node) (lnum, rnum string, err error) {
	lnum, err = left.GetString()
	if err != nil {
		return
	}
	rnum, err = right.GetString()
	return
}

func CastArrays(left, right *Node) (lnum, rnum []*Node, err error) {
	lnum, err = left.GetArray()
	if err != nil {
		return
	}
	rnum, err = right.GetArray()
	return
}

func CastObjects(left, right *Node) (lnum, rnum map[string]*Node, err error) {
	lnum, err = left.GetObject()
	if err != nil {
		return
	}
	rnum, err = right.GetObject()
	return
}

func CastBoolean(node *Node) (bool, error) {
	switch node.Type() {
	case Bool:
		return node.GetBool()
	case Numeric:
		res, err := node.GetNumeric()
		return res != 0, err
	case String:
		res, err := node.GetString()
		return res != "", err
	case Null:
		return false, nil
	case Array:
		fallthrough
	case Object:
		return !node.Empty(), nil
	}
	return false, nil
}

func CastString(key string) (string, bool) {
	bString := []byte(key)
	from := len(bString)
	if from > 1 && (bString[0] == BQuotes && bString[from-1] == BQuotes) {
		return internal.Unquote(bString, BQuotes)
	}
	if from > 1 && (bString[0] == BQuote && bString[from-1] == BQuote) {
		return internal.Unquote(bString, BQuote)
	}
	return key, true
	// todo quote string and unquote it:
	// {
	// 	bString = append([]byte{quotes}, bString...)
	// 	bString = append(bString, quotes)
	// }
	// return unquote(bString, quotes)
}

func CastNumericToFloat64(value interface{}) (result float64, err error) {
	switch typed := value.(type) {
	case float64:
		result = typed
	case float32:
		result = float64(typed)
	case int:
		result = float64(typed)
	case int8:
		result = float64(typed)
	case int16:
		result = float64(typed)
	case int32:
		result = float64(typed)
	case int64:
		result = float64(typed)
	case uint:
		result = float64(typed)
	case uint8:
		result = float64(typed)
	case uint16:
		result = float64(typed)
	case uint32:
		result = float64(typed)
	case uint64:
		result = float64(typed)
	default:
		err = NewUnsupportedType(value)
	}
	return
}
