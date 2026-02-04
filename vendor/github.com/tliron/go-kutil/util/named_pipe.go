package util

import (
	"os"
)

//
// NamedPipe
//

type NamedPipe struct {
	path   string
	writer *os.File
	reader *os.File
}

// ([io.Writer] interface)
func (self *NamedPipe) Write(p []byte) (int, error) {
	return self.writer.Write(p)
}

// ([io.Closer] interface)
func (self *NamedPipe) Close() error {
	self.reader.Close()
	self.writer.Close()
	return os.Remove(self.path)
}
