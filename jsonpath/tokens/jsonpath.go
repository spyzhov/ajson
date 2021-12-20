package tokens

import (
	"fmt"
	"strings"

	"github.com/spyzhov/ajson/v1/internal"
	"github.com/spyzhov/ajson/v1/jerrors"
)

type JSONPath struct {
	parent Token
	Tokens []Path
}

var _ Token = (*JSONPath)(nil)

func NewJSONPath(token string) (result *JSONPath, err error) {
	panic("not implemented")
	//return newJSONPath(internal.NewBuffer([]byte(token)))
}

// fixme: found the way how to stop when part of the other Buffer is given
func newJSONPath(b *internal.Buffer) (result *JSONPath, err error) {
	result = &JSONPath{Tokens: make([]Path, 0)}
	const (
		fQuote  = 1 << 0
		fQuotes = 1 << 1
	)
	var (
		c           byte
		eof         error
		start, stop int
		entrance    = b.Index
		childEnd    = map[byte]bool{internal.BDot: true, internal.BBracketL: true}
		//flag           int
		//brackets       int
		jsonpathErrFmt = fmt.Sprintf("JSONPath(starts from %d) found an error at %%d (%%q): %%w", entrance)
	)
bufferLoop:
	for {
		if c, eof = b.Current(); eof != nil {
			break
		}
		//parseSwitch:
		switch true {
		// $ : get root
		case c == internal.BDollar:
			if len(result.Tokens) != 0 {
				return nil, fmt.Errorf(jsonpathErrFmt, b.Index, string(c), jerrors.ErrIncorrectJSONPath)
			}
			result.append(NewRoot())
		// @ : get current; works only inside filters
		case c == internal.BAt:
			if len(result.Tokens) != 0 {
				return nil, fmt.Errorf(jsonpathErrFmt, b.Index, string(c), jerrors.ErrIncorrectJSONPath)
			}
			result.append(NewCurrent())
		// .  : get child with raw string key
		// .. : get all children
		case c == internal.BDot:
			if len(result.Tokens) == 0 {
				return nil, fmt.Errorf(jsonpathErrFmt, b.Index, string(c), jerrors.ErrIncorrectJSONPath)
			}
			start = b.Index
			if c, eof = b.Next(); eof != nil {
				// ends with dot, like `$.foo.bar.` - do nothing.
				break bufferLoop
			}
			if c == internal.BDot {
				result.append(NewRecursiveDescent())
				b.Index--
				break
			}
			if eof = b.SkipAny(childEnd); eof != nil { // fixme: here could be space in the string `$.foo.bar baz.biz`
				stop = b.Length
			} else {
				stop = b.Index
				b.Index--
			}
			if start+1 < stop {
				result.append(NewChild(NewRawString(string(b.Bytes[start+1 : stop]))))
			}
		// [ : get child ...
		// todo
		case c == internal.BBracketL:
			panic("not implemented")
			//_, err = b.Next()
			//if err != nil {
			//	return nil, b.ErrorEOF()
			//}
			//brackets = 1
			//start = b.Index
			//for ; b.Index < b.Length; b.Index++ {
			//	c = b.Bytes[b.Index]
			//	switch c {
			//	case internal.BQuote:
			//		if flag&fQuotes == 0 {
			//			if flag&fQuote == 0 {
			//				flag |= fQuote
			//			} else if !b.Backslash() {
			//				flag ^= fQuote
			//			}
			//		}
			//	case internal.BQuotes:
			//		if flag&fQuote == 0 {
			//			if flag&fQuotes == 0 {
			//				flag |= fQuotes
			//			} else if !b.Backslash() {
			//				flag ^= fQuotes
			//			}
			//		}
			//	case internal.BBracketL:
			//		if flag == 0 && !b.Backslash() {
			//			brackets++
			//		}
			//	case internal.BBracketR:
			//		if flag == 0 && !b.Backslash() {
			//			brackets--
			//		}
			//		if brackets == 0 {
			//			result = append(result, string(b.Bytes[start:b.Index]))
			//			break parseSwitch
			//		}
			//	}
			//}
			//return nil, b.ErrorEOF()
		default:
			return nil, b.ErrorSymbol()
		}
		if eof = b.Step(); eof != nil {
			break
		}
	}
	return
}

func (t *JSONPath) Type() string {
	return "JSONPath"
}

func (t *JSONPath) String() string {
	if t == nil {
		return "<nil>"
	}
	parts := make([]string, 0, len(t.Tokens))
	for _, token := range t.Tokens {
		parts = append(parts, token.String())
	}
	return strings.Join(parts, ", ")
}

func (t *JSONPath) Token() string {
	if t == nil {
		return "JSONPath(<nil>)"
	}
	parts := make([]string, 0, len(t.Tokens))
	for _, token := range t.Tokens {
		parts = append(parts, token.Token())
	}
	return fmt.Sprintf("JSONPath(%s)", strings.Join(parts, ", "))
}

func (t *JSONPath) Path() string {
	if t == nil {
		return "<nil>"
	}
	parts := make([]string, 0, len(t.Tokens))
	for _, token := range t.Tokens {
		parts = append(parts, token.Path())
	}
	return strings.Join(parts, "")
}

func (t *JSONPath) append(path Path) {
	t.Tokens = append(t.Tokens, path)
}

func (t *JSONPath) Parent() Token {
	if t == nil {
		return nil
	}
	return t.parent
}

func (t *JSONPath) SetParent(parent Token) {
	if t == nil {
		return
	}
	t.parent = parent
}

func (t *JSONPath) Append(token Token) error {
	if path, ok := token.(Path); ok {
		token.SetParent(t)
		t.Tokens = append(t.Tokens, path)
		return nil
	}
	return fmt.Errorf("%w: for JSONPath only Path is available, %s given", jerrors.ErrUnexpectedStatement, token.Type())
}

func (t *JSONPath) GetState(_ internal.State) internal.State {
	return -1 // fixme
}
