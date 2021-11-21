package tokens

type RecursiveDescent struct{}

var _ Path = (*RecursiveDescent)(nil)

func NewRecursiveDescent() *RecursiveDescent {
	return new(RecursiveDescent)
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
