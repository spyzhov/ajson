package tokens

import (
	"fmt"
	"io"
	"strings"

	"github.com/spyzhov/ajson/v1/internal"
)

type JSONPath struct {
	Tokens []Path
}

var _ Token = (*JSONPath)(nil)

func NewJSONPath(token string) (result *JSONPath, err error) {
	return newJSONPath(internal.NewBuffer([]byte(token)))
}

// todo
func newJSONPath(b *internal.Buffer) (result *JSONPath, err error) {
	result = &JSONPath{Tokens: make([]Path, 0)}
	const (
		fQuote  = 1 << 0
		fQuotes = 1 << 1
	)
	var (
		c              byte
		eof            error
		start, stop    int
		entrance       = b.Index
		childEnd       = map[byte]bool{internal.BDot: true, internal.BBracketL: true}
		flag           int
		brackets       int
		jsonpathErrFmt = fmt.Sprintf("JSONPath(starts from %d) found an error at %%d (%%q): %%w", entrance)
	)
bufferLoop:
	for {
		if c, eof = b.Current(); eof != nil {
			break
		}
	parseSwitch:
		switch true {
		// $ : get root
		case c == internal.BDollar:
			if len(result.Tokens) != 0 {
				return nil, fmt.Errorf(jsonpathErrFmt, b.Index, string(c), ErrIncorrectJSONPath)
			}
			result.Tokens = append(result.Tokens, NewRoot())
		// @ : get current; works only inside filters
		case c == internal.BAt:
			if len(result.Tokens) != 0 {
				return nil, fmt.Errorf(jsonpathErrFmt, b.Index, string(c), ErrIncorrectJSONPath)
			}
			result.Tokens = append(result.Tokens, NewCurrent())
		// .  : get child
		// .. : get all children
		case c == internal.BDot:
			if len(result.Tokens) == 0 {
				return nil, fmt.Errorf(jsonpathErrFmt, b.Index, string(c), ErrIncorrectJSONPath)
			}
			start = b.Index
			if c, eof = b.Next(); eof != nil {
				// ends with dot, like `$.foo.bar.` - do nothing.
				break bufferLoop
			}
			if c == internal.BDot {
				result.Tokens = append(result.Tokens, NewRecursiveDescent())
				b.Index--
				break
			}
			// todo: tbd
			err = b.SkipAny(childEnd)
			stop = b.Index
			if err == io.EOF {
				err = nil
				stop = b.Length
			} else {
				b.Index--
			}
			if err != nil {
				break
			}
			if start+1 < stop {
				result = append(result, string(b.Bytes[start+1:stop]))
			}
		case c == internal.BBracketL:
			_, err = b.Next()
			if err != nil {
				return nil, b.ErrorEOF()
			}
			brackets = 1
			start = b.Index
			for ; b.Index < b.Length; b.Index++ {
				c = b.Bytes[b.Index]
				switch c {
				case internal.BQuote:
					if flag&fQuotes == 0 {
						if flag&fQuote == 0 {
							flag |= fQuote
						} else if !b.Backslash() {
							flag ^= fQuote
						}
					}
				case internal.BQuotes:
					if flag&fQuote == 0 {
						if flag&fQuotes == 0 {
							flag |= fQuotes
						} else if !b.Backslash() {
							flag ^= fQuotes
						}
					}
				case internal.BBracketL:
					if flag == 0 && !b.Backslash() {
						brackets++
					}
				case internal.BBracketR:
					if flag == 0 && !b.Backslash() {
						brackets--
					}
					if brackets == 0 {
						result = append(result, string(b.Bytes[start:b.Index]))
						break parseSwitch
					}
				}
			}
			return nil, b.ErrorEOF()
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
