package tokens

type Current struct{}

var _ Path = (*Current)(nil)

func NewCurrent() *Current {
	return new(Current)
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
