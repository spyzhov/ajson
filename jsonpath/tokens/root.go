package tokens

type Root struct{}

var _ Path = (*Root)(nil)

func NewRoot() *Root {
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
