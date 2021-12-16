package tokens

type Wildcard struct {
	parent Token
}

var _ Path = (*Wildcard)(nil)

func NewWildcard() *Wildcard {
	panic("not implemented")
	//return new(Wildcard)
}

func (t *Wildcard) Type() string {
	return "Wildcard"
}

func (t *Wildcard) String() string {
	if t == nil {
		return "<nil>"
	}
	return "*"
}

func (t *Wildcard) Token() string {
	if t == nil {
		return "Wildcard(<nil>)"
	}
	return "Wildcard(*)"
}

func (t *Wildcard) Path() string {
	return "*"
}

func (t *Wildcard) Parent() Token {
	if t == nil {
		return nil
	}
	return t.parent
}
