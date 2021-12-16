package tokens

import (
	"fmt"
	"strings"
)

type Slice struct {
	parent Token
	Start  Token
	Stop   Token
	Step   Token
}

var _ Path = (*Slice)(nil)

func NewSlice(start, stop, step Token) (*Slice, error) {
	panic("not implemented")
	//for key, token := range map[string]Token{"start": start, "stop": stop, "step": step} {
	//	if token == nil {
	//		continue
	//	}
	//	switch token.(type) {
	//	case *Number:
	//	case *Script:
	//		break
	//	default:
	//		return nil, fmt.Errorf("slice argument %q has wrong type: %s", key, token.Type())
	//	}
	//}
	//
	//return &Slice{
	//	Start: start,
	//	Stop:  stop,
	//	Step:  step,
	//}, nil
}

func (t *Slice) Type() string {
	return "Slice"
}

func (t *Slice) String() string {
	if t == nil {
		return "<nil>"
	}
	parts := make([]string, 3)
	for i, token := range []Token{t.Start, t.Stop, t.Step} {
		if token != nil {
			parts[i] = token.String()
		}
	}

	return fmt.Sprintf("[%s]", strings.Join(parts, ":"))
}

func (t *Slice) Token() string {
	if t == nil {
		return "Slice(<nil>)"
	}
	parts := make([]string, 3)
	for i, token := range []Token{t.Start, t.Stop, t.Step} {
		if token != nil {
			parts[i] = token.Token()
		} else {
			parts[i] = "<nil>"
		}
	}

	return fmt.Sprintf("Slice(%s)", strings.Join(parts, ", "))
}

func (t *Slice) Path() string {
	if t == nil {
		return "<nil>"
	}
	parts := make([]string, 3)
	for i, token := range []Token{t.Start, t.Stop, t.Step} {
		if token != nil {
			if path, ok := token.(Path); ok {
				parts[i] = path.String()
			} else {
				parts[i] = token.String()
			}
		}
	}

	return fmt.Sprintf("[%s]", strings.Join(parts, ":"))
}

func (t *Slice) Parent() Token {
	if t == nil {
		return nil
	}
	return t.parent
}
