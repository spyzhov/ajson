package internal

import (
	"fmt"
	"io"

	"github.com/spyzhov/ajson/v1/jerrors"
)

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

type Buffer struct {
	Bytes  []byte
	Length int
	Index  int

	Last  State
	State State
	Class Class

	StateTransitionTable SST
}

type SST interface {
	GetState(state State, class Class) State
	GetClass(index byte) Class
}

var (
	C_null  = []byte("null")
	C_true  = []byte("true")
	C_false = []byte("false")

	brackets = map[byte]byte{
		BBracketL:     BBracketR,
		BBracesL:      BBracesR,
		BParenthesesL: BParenthesesR,
	}
)

func NewBuffer(body []byte, sst SST) *Buffer {
	if sst == nil {
		sst = StateTransitionTable
	}
	return &Buffer{
		Bytes:  body,
		Length: len(body),
		Index:  0,

		Last:  GO,
		State: GO,
		Class: C_SPACE,
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

func (b *Buffer) Move(length int) error {
	if b.Index+length >= b.Length {
		return io.EOF
	}
	b.Index += length
	return nil
}

func (b *Buffer) Reset() {
	b.Last = GO
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

// Word method moves Index to the first symbol which does not match [a-zA-Z0-9_]
func (b *Buffer) Word() {
	var c byte
	for ; b.Index < b.Length; b.Index++ {
		c = b.Bytes[b.Index]
		if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') && (c < '0' || c > '9') && c != '_' {
			break
		}
	}
}

// Backslash
// todo: add a lot of tests
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

func (b *Buffer) SkipBrackets() error {
	open := b.Bytes[b.Index]
	closed := brackets[open]
	if closed == 0 {
		return b.ErrorSymbol()
	}
	panic("not implemented")
	for ; b.Index < b.Length; b.Index++ {
		// todo: b.AsString()
	}
	return io.EOF
}

func (b *Buffer) GetClass() Class {
	return b.StateTransitionTable.GetClass(b.Bytes[b.Index])
}

func (b *Buffer) GetState() State {
	return b.State
}

func (b *Buffer) GetNextState() State {
	b.Last = b.State
	b.Class = b.GetClass()
	if b.Class == __ {
		return __
	}
	b.State = b.StateTransitionTable.GetState(b.Last, b.Class)
	return b.State
}

// AsNumeric if token is true - skip error from StateTransitionTable, just stop on unknown state
func (b *Buffer) AsNumeric() error {
	for ; b.Index < b.Length; b.Index++ {
		b.Class = b.GetClass()
		if b.Class == __ {
			return b.ErrorSymbol()
		}
		b.State = b.StateTransitionTable.GetState(b.Last, b.Class)
		if b.State == __ {
			return b.ErrorSymbol()
		}
		if b.State < __ {
			return nil
		}
		if b.State < MI || b.State > E3 {
			return nil
		}
		b.Last = b.State
	}
	// fixme: add tests
	if b.Class == C_MINUS || b.Class == C_PLUS || b.Class == C_POINT {
		//if b.Last == MI || b.Last == DT {
		return b.ErrorSymbol()
	}
	return nil
}

func (b *Buffer) AsString() error {
	if current, eof := b.Current(); eof != nil {
		return b.ErrorEOF()
	} else {
		return b.AsStringBordered(current)
	}
}

func (b *Buffer) AsStringBordered(border byte) error {
	if current, err := b.Current(); err != nil {
		return b.ErrorEOF()
	} else if current != border {
		return b.ErrorSymbol()
	}
	start := b.Index
	for ; b.Index < b.Length; b.Index++ {
		b.Class = b.GetClass()
		if b.Class == C_QUOTE {
			if b.Bytes[b.Index] != border {
				b.Class = b.StateTransitionTable.GetClass(-1) // C_ETC
			}
		}
		if b.Class == __ {
			return b.ErrorSymbol()
		}
		b.State = b.StateTransitionTable.GetState(b.Last, b.Class)
		if b.State == __ {
			return b.ErrorSymbol()
		}
		if b.State < __ { // goto:action
			return nil
		}
		b.Last = b.State
	}
	return b.ErrorUnfinished(start)
}

func (b *Buffer) AsNull() error {
	return b.strict(C_null)
}

func (b *Buffer) AsTrue() error {
	return b.strict(C_true)
}

func (b *Buffer) AsFalse() error {
	return b.strict(C_false)
}

func (b *Buffer) strict(word []byte) error {
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

func (b *Buffer) BytesFrom(start, end int) []byte {
	if start < 0 {
		start = 0
	} else if start > b.Length {
		start = b.Length
	}
	if end > b.Length {
		end = b.Length
	}
	if start > end {
		start = end
	}

	return b.Bytes[start:end]
}

func (b *Buffer) StringFrom(start int) string {
	return string(b.BytesFrom(start, b.Index))
}

func (b *Buffer) ErrorSymbol() error {
	return fmt.Errorf("%w at %d next to %q", jerrors.ErrWrongSymbol, b.Index, b.Bytes[b.Index])
}

func (b *Buffer) ErrorEOF() error {
	return fmt.Errorf("%w at %d", jerrors.ErrUnexpectedEOF, b.Index)
}

func (b *Buffer) ErrorUnfinished(start int) error {
	return fmt.Errorf("%w started from %d", jerrors.ErrUnfinishedToken, start)
}

func (b *Buffer) ErrorIncorrectJSONPath() error {
	return fmt.Errorf("%w started from %d", jerrors.ErrIncorrectJSONPath, b.Index)
}
