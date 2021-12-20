package tokens

type Root struct {
	parent Token
}

var _ Path = (*Root)(nil)

func NewRoot() *Root {
	panic("not implemented")
	return new(Root)
}

func (t *Root) Type() string {
	return "Root"
}

func (t *Root) String() string {
	if t == nil {
		return "<nil>"
	}
	return "$"
}

func (t *Root) Token() string {
	if t == nil {
		return "Root(<nil>)"
	}
	return "Root($)"
}

func (t *Root) Path() string {
	return "$"
}

func (t *Root) Parent() Token {
	if t == nil {
		return nil
	}
	return t.parent
}

func (t *Root) SetParent(parent Token) {
	if t == nil {
		return
	}
	t.parent = parent
}
