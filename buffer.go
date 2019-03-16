package ajson

import (
	"io"
)

type buffer struct {
	data   []byte
	length int
	index  int
}

const (
	quotes    byte = '"'
	coma      byte = ','
	colon     byte = ':'
	backslash byte = '\\'
	skipS     byte = ' '
	skipN     byte = '\n'
	skipR     byte = '\r'
	skipT     byte = '\t'
	bracketL  byte = '['
	bracketR  byte = ']'
	bracesL   byte = '{'
	bracesR   byte = '}'
)

var (
	_null  = []byte("null")
	_true  = []byte("true")
	_false = []byte("false")
)

func newBuffer(body []byte) (b *buffer) {
	b = &buffer{
		length: len(body),
		data:   body,
	}
	return
}

func (b *buffer) first() (c byte, err error) {
	for ; b.index < b.length; b.index++ {
		c = b.data[b.index]
		if !(c == skipS || c == skipR || c == skipN || c == skipT) {
			return c, nil
		}
	}
	return 0, io.EOF
}

func (b *buffer) backslash() (result bool) {
	for i := b.index - 1; i >= 0; i-- {
		if b.data[i] == backslash {
			result = !result
		} else {
			break
		}
	}
	return
}

func (b *buffer) skip(s byte) error {
	for ; b.index < b.length; b.index++ {
		if b.data[b.index] == s && !b.backslash() {
			return nil
		}
	}
	return io.EOF
}

func (b *buffer) numeric() error {
	var c byte
	find := 0
	for ; b.index < b.length; b.index++ {
		c = b.data[b.index]
		switch true {
		case c >= '0' && c <= '9':
			find |= 4
		case c == '.':
			if find&2 == 0 {
				find &= 2
			} else {
				return errorSymbol(b)
			}
		case c == '+' || c == '-':
			if find == 0 || find == 8 {
				find |= 1
			} else {
				return errorSymbol(b)
			}
		case c == 'e' || c == 'E':
			if find&8 == 0 {
				find = 8
			} else {
				return errorSymbol(b)
			}
		default:
			if find&4 != 0 {
				return nil
			}
			return errorSymbol(b)
		}
	}
	if find&4 != 0 {
		return io.EOF
	}
	return errorEOF(b)
}

func (b *buffer) string() error {
	err := b.step()
	if err != nil {
		return errorEOF(b)
	}
	if b.skip(quotes) != nil {
		return errorEOF(b)
	}
	return nil
}

func (b *buffer) null() error {
	return b.word(_null)
}

func (b *buffer) true() error {
	return b.word(_true)
}

func (b *buffer) false() error {
	return b.word(_false)
}

func (b *buffer) word(word []byte) error {
	var c byte
	max := len(word)
	index := 0
	for ; b.index < b.length; b.index++ {
		c = b.data[b.index]
		if c != word[index] && c != (word[index]-32) {
			return errorSymbol(b)
		}
		index++
		if index >= max {
			break
		}
	}
	if index != max {
		return errorEOF(b)
	}
	return nil
}

func (b *buffer) step() error {
	if b.index+1 < b.length {
		b.index++
		return nil
	}
	return io.EOF
}
