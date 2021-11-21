package tokens

type Wildcard struct{}

var _ Path = (*Wildcard)(nil)

func NewWildcard() *Wildcard {
	return new(Wildcard)
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
