package jsonpath

func New(jsonpath []byte) (*JSONPath, error) {
	// todo
	return nil, nil
}

func Compile(jsonpath string) (*JSONPath, error) {
	return New([]byte(jsonpath))
}

func MustCompile(jsonpath string) *JSONPath {
	result, err := Compile(jsonpath)
	if err != nil {
		panic(err)
	}
	return result
}
