package ajson

import (
	"io"
	"strings"

	"github.com/spyzhov/ajson/v1/internal"
)

type Buffer struct {
	Bytes  []byte
	Length int
	Index  int

	Last  internal.States
	State internal.States
	Class internal.Classes
}

const __ = -1

const (
	// USpace is a code for symbol `Space` (taken from www.json.org)
	USpace byte = '\u0020'
	// UNewLine is a code for symbol `New Line` or `\n` (taken from www.json.org)
	UNewLine byte = '\u000A'
	// UCarriageReturn is a code for symbol `Carriage Return` or `\r` (taken from www.json.org)
	UCarriageReturn byte = '\u000D'
	// UTab is a code for symbol `Tab` or `\t` (taken from www.json.org)
	UTab byte = '\u0009'

	BQuotes       byte = '"'
	BQuote        byte = '\''
	BComa         byte = ','
	BColon        byte = ':'
	BBackslash    byte = '\\'
	BSkipS             = USpace
	BSkipN             = UNewLine
	BSkipR             = UCarriageReturn
	BSkipT             = UTab
	BBracketL     byte = '['
	BBracketR     byte = ']'
	BBracesL      byte = '{'
	BBracesR      byte = '}'
	BParenthesesL byte = '('
	BParenthesesR byte = ')'
	BDollar       byte = '$'
	BAt           byte = '@'
	BDot          byte = '.'
	BAsterisk     byte = '*'
	BPlus         byte = '+'
	BMinus        byte = '-'
	BDivision     byte = '/'
	BExclamation  byte = '!'
	BCaret        byte = '^'
	BSignL        byte = '<'
	BSignG        byte = '>'
	BSignE        byte = '='
	BAmpersand    byte = '&'
	BPipe         byte = '|'
	BQuestion     byte = '?'
)

// RPN is a `Reverse Polish notation`
type RPN []string

var (
	C_null  = []byte("null")
	C_true  = []byte("true")
	C_false = []byte("false")
)

func NewBuffer(body []byte) *Buffer {
	return &Buffer{
		Bytes:  body,
		Length: len(body),
		Index:  0,

		Last:  internal.GO,
		State: internal.GO,
		Class: internal.C_SPACE,
	}
}

func (b *Buffer) Current() (c byte, err error) {
	if b.Index < b.Length {
		return b.Bytes[b.Index], nil
	}
	return 0, io.EOF
}

func (b *Buffer) Next() (c byte, err error) {
	err = b.Step()
	if err != nil {
		return 0, err
	}
	return b.Bytes[b.Index], nil
}

func (b *Buffer) Slice(length int) ([]byte, error) {
	if b.Index+length >= b.Length {
		return nil, io.EOF
	}
	return b.Bytes[b.Index : b.Index+length], nil
}

func (b *Buffer) Reset() {
	b.Last = internal.GO
}

func (b *Buffer) FirstNonSpace() (c byte, err error) {
	for ; b.Index < b.Length; b.Index++ {
		c = b.Bytes[b.Index]
		if !(c == BSkipS || c == BSkipR || c == BSkipN || c == BSkipT) {
			return c, nil
		}
	}
	return 0, io.EOF
}

func (b *Buffer) Backslash() (result bool) {
	for i := b.Index - 1; i >= 0; i-- {
		if b.Bytes[i] == BBackslash {
			result = !result
		} else {
			break
		}
	}
	return
}

func (b *Buffer) Skip(s byte) error {
	for ; b.Index < b.Length; b.Index++ {
		if b.Bytes[b.Index] == s && !b.Backslash() {
			return nil
		}
	}
	return io.EOF
}

func (b *Buffer) SkipAny(s map[byte]bool) error {
	for ; b.Index < b.Length; b.Index++ {
		if s[b.Bytes[b.Index]] && !b.Backslash() {
			return nil
		}
	}
	return io.EOF
}

// AsNumeric if token is true - skip error from StateTransitionTable, just stop on unknown state
func (b *Buffer) AsNumeric(token bool) error {
	if token {
		b.Last = internal.GO
	}
	for ; b.Index < b.Length; b.Index++ {
		b.Class = b.GetClasses(BQuotes)
		if b.Class == __ {
			return b.ErrorSymbol()
		}
		b.State = internal.StateTransitionTable[b.Last][b.Class]
		if b.State == __ {
			if token {
				break
			}
			return b.ErrorSymbol()
		}
		if b.State < __ {
			return nil
		}
		if b.State < internal.MI || b.State > internal.E3 {
			return nil
		}
		b.Last = b.State
	}
	if b.Last != internal.ZE && b.Last != internal.IN && b.Last != internal.FR && b.Last != internal.E3 {
		return b.ErrorSymbol()
	}
	return nil
}

func (b *Buffer) GetClasses(search byte) internal.Classes {
	if b.Bytes[b.Index] >= 128 {
		return internal.C_ETC
	}
	if search == BQuote {
		return internal.QuoteAsciiClasses[b.Bytes[b.Index]]
	}
	return internal.AsciiClasses[b.Bytes[b.Index]]
}

func (b *Buffer) GetState() internal.States {
	b.Last = b.State
	b.Class = b.GetClasses(BQuotes)
	if b.Class == __ {
		return __
	}
	b.State = internal.StateTransitionTable[b.Last][b.Class]
	return b.State
}

func (b *Buffer) AsString(search byte, token bool) error {
	if token {
		b.Last = internal.GO
	}
	for ; b.Index < b.Length; b.Index++ {
		b.Class = b.GetClasses(search)

		if b.Class == __ {
			return b.ErrorSymbol()
		}
		b.State = internal.StateTransitionTable[b.Last][b.Class]
		if b.State == __ {
			return b.ErrorSymbol()
		}
		if b.State < __ {
			return nil
		}
		b.Last = b.State
	}
	return b.ErrorSymbol()
}

func (b *Buffer) AsNull() error {
	return b.word(C_null)
}

func (b *Buffer) AsTrue() error {
	return b.word(C_true)
}

func (b *Buffer) AsFalse() error {
	return b.word(C_false)
}

func (b *Buffer) word(word []byte) error {
	var c byte
	max := len(word)
	index := 0
	for ; b.Index < b.Length; b.Index++ {
		c = b.Bytes[b.Index]
		// if c != word[index] && c != (word[index]-32) {
		if c != word[index] {
			return b.ErrorSymbol()
		}
		index++
		if index >= max {
			break
		}
	}
	if index != max {
		return b.ErrorEOF()
	}
	return nil
}

func (b *Buffer) Step() error {
	if b.Index+1 < b.Length {
		b.Index++
		return nil
	}
	return io.EOF
}

func (b *Buffer) Back() {
	if b.Index > 0 {
		b.Index--
	}
}

// Token reads until the end of the token e.g.: `@.length`, `@['foo'].bar[(@.length - 1)].baz`
func (b *Buffer) Token() (err error) {
	var (
		c     byte
		stack = make([]byte, 0)
		first = b.Index
		start int
		find  bool
	)
tokenLoop:
	for ; b.Index < b.Length; b.Index++ {
		c = b.Bytes[b.Index]
		switch {
		case c == BQuotes:
			fallthrough
		case c == BQuote:
			find = true
			err = b.Step()
			if err != nil {
				return b.ErrorEOF()
			}
			err = b.Skip(c)
			if err == io.EOF {
				return b.ErrorEOF()
			}
		case c == BBracketL:
			find = true
			stack = append(stack, c)
		case c == BBracketR:
			find = true
			if len(stack) == 0 {
				if first == b.Index {
					return b.ErrorSymbol()
				}
				break tokenLoop
			}
			if stack[len(stack)-1] != BBracketL {
				return b.ErrorSymbol()
			}
			stack = stack[:len(stack)-1]
		case c == BParenthesesL:
			find = true
			stack = append(stack, c)
		case c == BParenthesesR:
			find = true
			if len(stack) == 0 {
				if first == b.Index {
					return b.ErrorSymbol()
				}
				break tokenLoop
			}
			if stack[len(stack)-1] != BParenthesesL {
				return b.ErrorSymbol()
			}
			stack = stack[:len(stack)-1]
		case c == BDot || c == BAt || c == BDollar || c == BQuestion || c == BAsterisk || (c >= 'A' && c <= 'z') || (c >= '0' && c <= '9'): // standard Token name
			find = true
			continue
		case len(stack) != 0:
			find = true
			continue
		case c == BMinus || c == BPlus:
			if !find {
				find = true
				start = b.Index
				err = b.AsNumeric(true)
				if err == nil || err == io.EOF {
					b.Index--
					continue
				}
				b.Index = start
			}
			fallthrough
		default:
			break tokenLoop
		}
	}
	if len(stack) != 0 {
		return b.ErrorEOF()
	}
	if first == b.Index {
		return b.Step()
	}
	if b.Index >= b.Length {
		return io.EOF
	}
	return nil
}

// RPN is a builder for `Reverse Polish notation`
// fixme: remove
// Deprecated
func (b *Buffer) RPN() (result RPN, err error) {
	var (
		c        byte
		start    int
		temp     string
		current  string
		found    bool
		variable bool
		stack    = make([]string, 0)
	)
	for {
		b.Reset()
		c, err = b.FirstNonSpace()
		if err != nil {
			break
		}
		switch true {
		case c == BAsterisk || c == BDivision || c == BMinus || c == BPlus || c == BCaret || c == BAmpersand || c == BPipe || c == BSignL || c == BSignG || c == BSignE || c == BExclamation: // operations
			if variable {
				variable = false
				current = string(c)

				c, err = b.Next()
				if err == nil {
					temp = current + string(c)
					if OperationsPriority[temp] != 0 {
						current = temp
					} else {
						b.Index--
					}
				} else {
					err = nil
				}

				for len(stack) > 0 {
					temp = stack[len(stack)-1]
					found = false
					if temp[0] >= 'A' && temp[0] <= 'z' { // function
						found = true
					} else if OperationsPriority[temp] != 0 { // operation
						if OperationsPriority[temp] > OperationsPriority[current] {
							found = true
						} else if OperationsPriority[temp] == OperationsPriority[current] && !RightOp[temp] {
							found = true
						}
					}

					if found {
						stack = stack[:len(stack)-1]
						result = append(result, temp)
					} else {
						break
					}
				}
				stack = append(stack, current)
				break
			}
			if c != BMinus && c != BPlus {
				return nil, b.ErrorSymbol()
			}
			fallthrough // for numbers like `-1e6`
		case (c >= '0' && c <= '9') || c == '.': // numbers
			variable = true
			start = b.Index
			err = b.AsNumeric(true)
			if err != nil {
				return nil, err
			}
			current = string(b.Bytes[start:b.Index])
			result = append(result, current)
			b.Index--
		case c == BQuotes: // string
			fallthrough
		case c == BQuote: // string
			variable = true
			start = b.Index
			err = b.AsString(c, true)
			if err != nil {
				return nil, b.ErrorEOF()
			}
			current = string(b.Bytes[start : b.Index+1])
			result = append(result, current)
		case c == BDollar || c == BAt: // variable : like @.length , $.expensive, etc.
			variable = true
			start = b.Index
			err = b.Token()
			if err != nil {
				if err != io.EOF {
					return nil, err
				}
			}
			current = string(b.Bytes[start:b.Index])
			result = append(result, current)
			if err != nil {
				err = nil
			} else {
				b.Index--
			}
		case c == BParenthesesL: // (
			variable = false
			current = string(c)
			stack = append(stack, current)
		case c == BParenthesesR: // )
			variable = true
			found = false
			for len(stack) > 0 {
				temp = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if temp == "(" {
					found = true
					break
				}
				result = append(result, temp)
			}
			if !found { // have no BParenthesesL
				return nil, NewErrorRequest("formula has no left parentheses")
			}
		default: // prefix functions or etc.
			start = b.Index
			variable = true
			for ; b.Index < b.Length; b.Index++ {
				c = b.Bytes[b.Index]
				if c == BParenthesesL { // function detection, example: sin(...), round(...), etc.
					variable = false
					break
				}
				if c < 'A' || c > 'z' {
					if !(c >= '0' && c <= '9') && c != '_' { // constants detection, example: true, false, null, PI, e, etc.
						break
					}
				}
			}
			current = strings.ToLower(string(b.Bytes[start:b.Index]))
			b.Index--
			if !variable {
				if _, found = Functions[current]; !found {
					return nil, NewErrorRequest("wrong formula, '%s' is not a function", current)
				}
				stack = append(stack, current)
			} else {
				if _, found = Constants[current]; !found {
					return nil, NewErrorRequest("wrong formula, '%s' is not a constant", current)
				}
				result = append(result, current)
			}
		}
		err = b.Step()
		if err != nil {
			break
		}
	}
	if err == io.EOF {
		err = nil // only io.EOF can be here
	}

	for len(stack) > 0 {
		temp = stack[len(stack)-1]
		_, ok := Functions[temp]
		if OperationsPriority[temp] == 0 && !ok { // operations only
			return nil, NewErrorRequest("wrong formula, '%s' is not an operation or function", temp)
		}
		result = append(result, temp)
		stack = stack[:len(stack)-1]
	}

	if len(result) == 0 {
		return nil, b.ErrorEOF()
	}

	return
}

func (b *Buffer) GetTokens() (result Tokens, err error) {
	var (
		c        byte
		start    int
		temp     string
		current  string
		variable bool
	)
	for {
		b.Reset()
		c, err = b.FirstNonSpace()
		if err != nil {
			break
		}
		switch true {
		case OperationsChar[c]: // operations
			if variable || (c != BMinus && c != BPlus) {
				variable = false
				current = string(c)

				c, err = b.Next()
				if err == nil {
					temp = current + string(c)
					if OperationsPriority[temp] != 0 {
						current = temp
					} else {
						b.Index--
					}
				} else {
					err = nil
				}

				result = append(result, current)
				break
			}
			fallthrough // for numbers like `-1e6`
		case (c >= '0' && c <= '9') || c == BDot: // numbers
			variable = true
			start = b.Index
			err = b.AsNumeric(true)
			if err != nil && err != io.EOF {
				if c == BDot {
					err = nil
					result = append(result, ".")
					b.Index = start
					break
				}
				return nil, err
			}
			current = string(b.Bytes[start:b.Index])
			result = append(result, current)
			b.Index--
		case c == BQuotes: // string
			fallthrough
		case c == BQuote: // string
			variable = true
			start = b.Index
			err = b.AsString(c, true)
			if err != nil {
				return nil, b.ErrorEOF()
			}
			current = string(b.Bytes[start : b.Index+1])
			result = append(result, current)
		case c == BDollar || c == BAt: // variable : like @.length , $.expensive, etc.
			variable = true
			start = b.Index
			err = b.Token()
			if err != nil {
				if err != io.EOF {
					return nil, err
				}
			}
			current = string(b.Bytes[start:b.Index])
			result = append(result, current)
			if err != nil {
				err = nil
			} else {
				b.Index--
			}
		case c == BParenthesesL: // (
			variable = false
			current = string(c)
			result = append(result, current)
		case c == BParenthesesR: // )
			variable = true
			current = string(c)
			result = append(result, current)
		default: // prefix functions or etc.
			start = b.Index
			variable = true
			for ; b.Index < b.Length; b.Index++ {
				c = b.Bytes[b.Index]
				if c == BParenthesesL { // function detection, example: sin(...), round(...), etc.
					variable = false
					break
				}
				if c < 'A' || c > 'z' {
					if !(c >= '0' && c <= '9') && c != '_' { // constants detection, example: true, false, null, PI, e, etc.
						break
					}
				}
			}
			if start == b.Index {
				err = b.Step()
				if err != nil {
					err = nil
					current = strings.ToLower(string(b.Bytes[start : b.Index+1]))
				} else {
					current = strings.ToLower(string(b.Bytes[start:b.Index]))
					b.Index--
				}
			} else {
				current = strings.ToLower(string(b.Bytes[start:b.Index]))
				b.Index--
			}
			result = append(result, current)
		}
		err = b.Step()
		if err != nil {
			break
		}
	}

	if err == io.EOF {
		err = nil
	}

	return
}

func (b *Buffer) ErrorEOF() error {
	return NewErrorEOF(b)
}

func (b *Buffer) ErrorSymbol() error {
	return NewErrorSymbol(b)
}
