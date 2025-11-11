//go:build linux || darwin

package commonlog

import (
	"bufio"
	"os"
	"path/filepath"
	"syscall"

	"github.com/segmentio/ksuid"
)

//
// LoggerFIFO
//

// A Linux FIFO file that forwards all lines written to it to a [Logger].
type LoggerFIFO struct {
	Path  string
	Log   Logger
	Level Level
}

func NewLoggerFIFO(prefix string, log Logger, level Level) *LoggerFIFO {
	path := filepath.Join(os.TempDir(), prefix+ksuid.New().String())
	return &LoggerFIFO{
		Path:  path,
		Log:   NewKeyValueLogger(log, "fifo", path),
		Level: level,
	}
}

func (self *LoggerFIFO) Start() error {
	if err := self.create(); err == nil {
		go self.start()
		return nil
	} else {
		return err
	}
}

func (self *LoggerFIFO) create() error {
	if err := os.Remove(self.Path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	self.Log.Debug("creating logger FIFO")
	return syscall.Mkfifo(self.Path, 0600)
}

func (self *LoggerFIFO) start() {
	// Note: os.Open will block until the FIFO will be opened for write
	if file, err := os.Open(self.Path); err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			self.Log.Log(self.Level, 0, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			self.Log.Error(err.Error())
		}

		self.Log.Debug("closing logger FIFO")
		if err := file.Close(); err != nil {
			self.Log.Error(err.Error())
		}
		if err := os.Remove(self.Path); err != nil {
			self.Log.Error(err.Error())
		}
	} else {
		self.Log.Critical(err.Error())
	}
}
