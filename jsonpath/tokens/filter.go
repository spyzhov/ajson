package tokens

import (
	"fmt"
)

type Filter struct {
	*RPN
}

var _ Token = (*Filter)(nil)

func NewFilter(rpn *RPN) (*Filter, error) {
	return &Filter{
		RPN: rpn,
	}, nil
}

func (t *Filter) Type() string {
	return "Filter"
}

func (t *Filter) String() string {
	if t == nil {
		return "<nil>"
	}
	return fmt.Sprintf("?(%s)", t.RPN.String())
}

func (t *Filter) Token() string {
	if t == nil {
		return "Filter(<nil>)"
	}
	return fmt.Sprintf("Filter(%s)", t.RPN.Token())
}

func (t *Filter) Path() string {
	if t == nil {
		return "?(<nil>)"
	}
	return fmt.Sprintf("?(%s)", t.RPN.String())
}
