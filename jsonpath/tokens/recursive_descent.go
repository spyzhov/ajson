package tokens

type RecursiveDescent struct {
	parent Token
}

var _ Path = (*RecursiveDescent)(nil)

func NewRecursiveDescent() *RecursiveDescent {
	panic("not implemented")
	return &RecursiveDescent{}
}

func (t *RecursiveDescent) Type() string {
	return "RecursiveDescent"
}

func (t *RecursiveDescent) String() string {
	return t.Token()
}

func (t *RecursiveDescent) Token() string {
	if t == nil {
		return "RecursiveDescent(<nil>)"
	}
	return "RecursiveDescent(..)"
}

func (t *RecursiveDescent) Path() string {
	return ".."
}

func (t *RecursiveDescent) Parent() Token {
	if t == nil {
		return nil
	}
	return t.parent
}

func (t *RecursiveDescent) SetParent(parent Token) {
	if t == nil {
		return
	}
	t.parent = parent
}
