package util

import (
	"archive/tar"
	"io"
	"strings"
)

//
// TarballReader
//

type TarballReader struct {
	TarReader         *tar.Reader
	ArchiveReader     io.ReadCloser
	CompressionReader io.ReadCloser
}

func NewTarballReader(reader *tar.Reader, archiveReader io.ReadCloser, compressionReader io.ReadCloser) *TarballReader {
	return &TarballReader{reader, archiveReader, compressionReader}
}

// ([io.Closer] interface)
func (self *TarballReader) Close() error {
	var err1 error
	if self.CompressionReader != nil {
		err1 = self.CompressionReader.Close()
	}
	err2 := self.ArchiveReader.Close()
	if err1 != nil {
		return err1
	} else {
		return err2
	}
}

func (self *TarballReader) Open(path string) (*TarballEntryReader, error) {
	for {
		if header, err := self.TarReader.Next(); err == nil {
			if path == FixTarballEntryPath(header.Name) {
				return NewTarballEntryReader(self), nil
			}
		} else if err == io.EOF {
			break
		} else {
			return nil, err
		}
	}
	return nil, nil
}

func (self *TarballReader) Has(path string) (bool, error) {
	for {
		if header, err := self.TarReader.Next(); err == nil {
			if path == FixTarballEntryPath(header.Name) {
				return true, nil
			}
		} else if err == io.EOF {
			break
		} else {
			return false, err
		}
	}
	return false, nil
}

func (self *TarballReader) Iterate(f func(*tar.Header) bool) error {
	for {
		if header, err := self.TarReader.Next(); err == nil {
			if !f(header) {
				return nil
			}
		} else if err == io.EOF {
			break
		} else {
			return err
		}
	}
	return nil
}

//
// TarballEntryReader
//

type TarballEntryReader struct {
	TarballReader *TarballReader
}

func NewTarballEntryReader(tarballReader *TarballReader) *TarballEntryReader {
	return &TarballEntryReader{tarballReader}
}

// ([io.Reader] interface)
func (self *TarballEntryReader) Read(p []byte) (n int, err error) {
	return self.TarballReader.TarReader.Read(p)
}

// ([io.Closer] interface)
func (self *TarballEntryReader) Close() error {
	return self.TarballReader.Close()
}

// Utils

func FixTarballEntryPath(path string) string {
	if strings.HasPrefix(path, "./") {
		return path[3:]
	}
	return path
}
