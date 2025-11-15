//go:build !(linux || darwin)

package commonlog

import (
	"errors"
)

//
// LoggerFIFO
//

type LoggerFIFO struct {
}

func NewLoggerFIFO(prefix string, log Logger, level Level) *LoggerFIFO {
	return new(LoggerFIFO)
}

func (self *LoggerFIFO) Start() error {
	return errors.New("not supported on this platform")
}
