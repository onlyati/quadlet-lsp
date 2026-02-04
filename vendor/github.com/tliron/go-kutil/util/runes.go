package util

import (
	"io"
	"unicode/utf8"
)

//
// RuneReader
//

type RuneReader struct {
	runes  []rune
	length int
	index  int
}

func NewRuneReader(runes []rune) *RuneReader {
	return &RuneReader{
		runes:  runes,
		length: len(runes),
	}
}

// ([io.RuneReader] interface)
func (self *RuneReader) ReadRune() (rune, int, error) {
	if self.index >= self.length {
		return 0, 0, io.EOF
	}

	rune_ := self.runes[self.index]
	self.index++
	return rune_, utf8.RuneLen(rune_), nil
}
