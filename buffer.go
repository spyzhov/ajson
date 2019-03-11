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
	backslash byte = '\\'
	skipS     byte = ' '
	skipN     byte = '\n'
	skipT     byte = '\t'
	bracketL  byte = '['
	bracketR  byte = ']'
	bracesL   byte = '{'
	bracesR   byte = '}'
)

func newBuffer(body []byte, clone bool) (b *buffer) {
	b = &buffer{
		length: len(body),
	}
	if clone {
		copy(body, b.data)
	} else {
		b.data = body
	}
	return
}

func (b *buffer) first() (c byte, err error) {
	for ; b.index < b.length; b.index++ {
		c = b.data[b.index]
		if !(c == skipS || c == skipN || c == skipT) {
			return c, nil
		}
	}
	return 0, io.EOF
}

func (b *buffer) next() (c byte, err error) {
	b.index++
	return b.first()
}

func (b *buffer) scan(s byte, skip bool) (from, to int) {
	var c byte
	find := false
	from = b.index
	to = b.index
	for ; b.index < b.length; b.index++ {
		c = b.data[b.index]
		if c == s && !b.backslash() {
			b.index++
			return from, to
		}
		if skip && (c == skipS || c == skipN || c == skipT) {
			if !find {
				from++
				to++
			}
		} else {
			find = true
			to++
		}
	}
	return -1, -1
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

func (b *buffer) skip(s byte) bool {
	for ; b.index < b.length; b.index++ {
		if b.data[b.index] == s && !b.backslash() {
			b.index++
			return true
		}
	}
	return false
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
				return errorSymbol(c, b.index)
			}
		case c == '+' || c == '-':
			if find == 0 || find == 8 {
				find |= 1
			} else {
				return errorSymbol(c, b.index)
			}
		case c == 'e' || c == 'E':
			if find&8 == 0 {
				find = 8
			} else {
				return errorSymbol(c, b.index)
			}
		default:
			if find&4 != 0 {
				return nil
			}
			return errorSymbol(c, b.index)
		}
	}
	return io.EOF
}

func (b *buffer) string() error {
	if !b.skip(quotes) {
		return io.EOF
	}
	return nil
}
