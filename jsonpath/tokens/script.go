package tokens

import (
	"fmt"

	"github.com/spyzhov/ajson/v1/jsonpath/internal"
)

type Script struct {
	RPN
}

var _ Token = (*Script)(nil)

func NewScript() *Script {
	return &Script{
		RPN: *NewRPN(),
	}
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

func (t *Script) GetState(_ internal.State) internal.State {
	return internal.ѢѢ // fixme
}
