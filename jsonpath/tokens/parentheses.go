package tokens

type parentheses struct {
	parent Token
	Open   bool
}

var _ Token = (*parentheses)(nil)

func newParentheses(char byte) (result *parentheses, err error) {
	panic("not implemented")
	//switch char {
	//case internal.BParenthesesL:
	//	result = &parentheses{
	//		parent: parent,
	//		Open:   true,
	//	}
	//case internal.BParenthesesR:
	//	result = &parentheses{
	//		parent: parent,
	//		Open:   false,
	//	}
	//default:
	//	err = fmt.Errorf("given value is not `parentheses`, char index: %d, value: %q", char, char)
	//}
	//return
}

func (t *parentheses) Type() string {
	return "parentheses"
}

func (t *parentheses) IsOpen() bool {
	if t == nil {
		return false
	}
	return t.Open
}

func (t *parentheses) String() string {
	if t.IsOpen() {
		return "("
	}
	return ")"
}

func (t *parentheses) Token() string {
	if t.IsOpen() {
		return "parentheses(Open)"
	}
	return "parentheses(Close)"
}

func (t *parentheses) Parent() Token {
	if t == nil {
		return nil
	}
	return t.parent
}

func (t *parentheses) SetParent(parent Token) {
	if t == nil {
		return
	}
	t.parent = parent
}
