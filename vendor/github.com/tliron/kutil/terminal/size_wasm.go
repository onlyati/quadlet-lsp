package terminal

import (
	"errors"
)

func GetSize() (int, int, error) {
	return -1, -1, errors.New("terminal size not supported in WASM")
}
