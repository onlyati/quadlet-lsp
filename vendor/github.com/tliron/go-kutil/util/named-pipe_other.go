//go:build !linux

package util

import (
	"errors"
	"io"
)

func NewTempNamedPipe(writer io.Writer, mode uint32) (*NamedPipe, error) {
	return nil, errors.New("named pipes are not supported on this platform")
}
