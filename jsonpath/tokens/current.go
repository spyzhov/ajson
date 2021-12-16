package tokens

type Current struct {
	parent Token
}

var _ Path = (*Current)(nil)

func NewCurrent() *Current {
	panic("not implemented")
	//return new(Current)
}

func (t *Current) Type() string {
	return "Current"
}

func (t *Current) String() string {
	if t == nil {
		return "<nil>"
	}
	return "@"
}

func (t *Current) Token() string {
	if t == nil {
		return "Current(<nil>)"
	}
	return "Current(@)"
}

func (t *Current) Path() string {
	return "@"
}

func (t *Current) Parent() Token {
	if t == nil {
		return nil
	}
	return t.parent
}
