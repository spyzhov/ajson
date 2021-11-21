package tokens

import (
	"fmt"
)

type Script struct {
	*RPN
}

var _ Token = (*Script)(nil)

func NewScript(rpn *RPN) (*Script, error) {
	return &Script{
		RPN: rpn,
	}, nil
}

func (t *Script) Type() string {
	return "Script"
}

func (t *Script) String() string {
	if t == nil {
		return "<nil>"
	}
	return fmt.Sprintf("(%s)", t.RPN.String())
}

func (t *Script) Token() string {
	if t == nil {
		return "Script(<nil>)"
	}
	return fmt.Sprintf("Script(%s)", t.RPN.Token())
}

func (t *Script) Path() string {
	if t == nil {
		return "(<nil>)"
	}
	return fmt.Sprintf("(%s)", t.RPN.String())
}
