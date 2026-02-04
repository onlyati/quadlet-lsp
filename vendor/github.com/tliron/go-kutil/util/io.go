package util

import (
	contextpkg "context"
	"errors"
	"io"
	"sync"
	"testing"
)

const BUFFER_SIZE = 65536

func WriteNewline(writer io.Writer) error {
	_, err := io.WriteString(writer, "\n")
	return err
}

func ReaderSize(reader io.Reader) (int64, error) {
	var size int64 = 0

	buffer := make([]byte, BUFFER_SIZE)
	for {
		count, err := reader.Read(buffer)
		size += int64(count)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return 0, err
			}
		}
	}

	return size, nil
}

func ContextualRead(context contextpkg.Context, reader io.Reader, p []byte) (int, error) {
	if err := context.Err(); err == nil {
		return reader.Read(p)
	} else {
		return 0, err
	}
}

// Sends the byte slice to a channel and returns its length.
//
// When copy is true will copy the byte slice before sending it to the
// channel. This is necessary for ensuring that the submitted data is
// indeed the data that will be received, even if the underlying array
// changes after this call.
func WriteBytesToChannel(ch chan []byte, p []byte, copy bool) (int, error) {
	if copy {
		p = append(p[:0:0], p...)
	}

	length := len(p)

	if length > 0 {
		select {
		case ch <- p:
		default:
			return 0, errors.New("channel full")
		}
	}

	return length, nil
}

//
// BufferedWriter
//

// https://www.reddit.com/r/golang/comments/6irpt1/is_there_a_golang_logging_library_around_that/dj8yicz?utm_source=share&utm_medium=web2x&context=3
// https://gobyexample.com/closing-channels

type BufferedWriter struct {
	writer      io.Writer
	copy        bool
	submissions chan []byte
	close       chan struct{}
	closed      chan struct{}
}

// Creates a thread-safe [io.WriterCloser] that separates the acceptance
// of byte slices via [BufferedWriter.Write] from the actual writing of
// the data to the underlying [io.Writer], which occurs on a separate
// goroutine.
//
// The size argument is for the internal channel, referring to the maximum
// number of write submissions to be buffered (not the number of bytes). If
// the buffering channel is full then [BufferedWriter.Write] will return an
// error. The actual write errors are ignored.
//
// When copy is true the implementation copies the byte slice before sending
// it to the separate goroutine. This ensures that the submitted data is
// indeed the data that will be written, even if the underlying array
// changes after submission.
func NewBufferedWriter(writer io.Writer, size int, copy bool) *BufferedWriter {
	self := BufferedWriter{
		writer:      writer,
		copy:        copy,
		submissions: make(chan []byte, size),
		close:       make(chan struct{}, 1),
		closed:      make(chan struct{}, 1),
	}

	go self.run()

	return &self
}

// ([io.Writer] interface)
func (self *BufferedWriter) Write(p []byte) (int, error) {
	defer func() {
		if recover() != nil {
			// The channel was closed
			//fmt.Println("closed!!!!")
			self.writer.Write(p)
		}
	}()

	return WriteBytesToChannel(self.submissions, p, self.copy)
}

// ([io.Closer] interface)
func (self *BufferedWriter) Close() error {
	defer func() {
		recover()
	}()

	close(self.submissions)
	<-self.closed
	return nil
}

func (self *BufferedWriter) run() {
	for {
		select {
		case submission, ok := <-self.submissions:
			if ok {
				self.writer.Write(submission) // ignores errors!
			} else {
				self.closed <- struct{}{}
				return
			}
		}
	}
}

//
// SyncedWriter
//

type SyncedWriter struct {
	Writer io.Writer

	lock sync.Mutex
}

// Creates a thread-safe [io.WriterCloser] that synchronizes all
// the [SyncedWriter.Write] calls via a mutex.
//
// If the underlying writer does not support [io.Closer] then
// [SyncedWriter.Close] will do nothing.
func NewSyncedWriter(writer io.Writer) *SyncedWriter {
	return &SyncedWriter{
		Writer: writer,
	}
}

// ([io.Writer] interface)
func (self *SyncedWriter) Write(p []byte) (int, error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	return self.Writer.Write(p)
}

// ([io.Closer] interface)
func (self *SyncedWriter) Close() error {
	self.lock.Lock()
	defer self.lock.Unlock()

	if closer, ok := self.Writer.(io.Closer); ok {
		return closer.Close()
	} else {
		return nil
	}
}

//
// ChannelWriter
//

type ChannelWriter struct {
	ch   chan []byte
	copy bool
}

// Creates an [io.Writer] that writes bytes to a channel. The expectation is that
// something else will be receiving from the channel and processing the bytes.
// Writing is non-blocking: if the channel is full then [ChannelWriter.Write] will
// return an error.
//
// When copy is true the implementation copies the byte slice before sending it to
// the channel. This ensures that the submitted data is indeed the data that will be
// received, even if the underlying array changes after submission.
func NewChannelWriter(ch chan []byte, copy bool) *ChannelWriter {
	return &ChannelWriter{ch, copy}
}

// ([io.Writer] interface)
func (self *ChannelWriter) Write(p []byte) (int, error) {
	return WriteBytesToChannel(self.ch, p, self.copy)
}

//
// ChannelReader
//

type ChannelReader struct {
	reader *io.PipeReader
}

// Creates an [io.Reader] that reads bytes from a channel. This works via
// a separate goroutine that pipes (via [io.Pipe]) the channel data.
func NewChannelReader(ch chan []byte) *ChannelReader {
	reader, writer := io.Pipe()

	go func() {
		defer writer.Close() // Note: always returns nil
		for p := range ch {
			if _, err := writer.Write(p); err != nil {
				return
			}
		}
	}()

	return &ChannelReader{reader: reader}
}

// ([io.Reader] interface)
func (self *ChannelReader) Read(p []byte) (n int, err error) {
	return self.reader.Read(p)
}

//
// ContextualReader
//

// https://pace.dev/blog/2020/02/03/context-aware-ioreader-for-golang-by-mat-ryer.html
// https://medium.com/@zombiezen/canceling-i-o-in-go-capn-proto-5ae8c09c5b29
// https://benjamincongdon.me/blog/2020/04/23/Cancelable-Reads-in-Go/
// https://github.com/muesli/cancelreader

// Creates a [io.Reader] that allows cancellation of reading via a provided
// context. Note that cancellation is not guaranteed because the underlying reader
// is called in-thread. At best we can guarantee that subsequent calls to
// [ContextualReader.Read] will fail after cancellation.
type ContextualReader struct {
	reader  io.Reader
	context contextpkg.Context
}

func NewContextualReader(context contextpkg.Context, reader io.Reader) io.Reader {
	return &ContextualReader{reader: reader, context: context}
}

// ([io.Reader] interface)
func (self *ContextualReader) Read(p []byte) (int, error) {
	return ContextualRead(self.context, self.reader, p)
}

//
// ContextualReadCloser
//

// Creates a [io.ReadCloser] that allows cancellation of reading via a provided
// context. Note that cancellation is not guaranteed because the underlying reader
// is called in-thread. At best we can guarantee that subsequent calls to
// [ContextualReadCloser.Read] will fail after cancellation.
type ContextualReadCloser struct {
	reader  io.ReadCloser
	context contextpkg.Context
}

// Creates a [io.ReadCloser]
func NewContextualReadCloser(context contextpkg.Context, reader io.ReadCloser) io.ReadCloser {
	return &ContextualReadCloser{reader: reader, context: context}
}

// ([io.Reader] interface)
func (self *ContextualReadCloser) Read(p []byte) (int, error) {
	return ContextualRead(self.context, self.reader, p)
}

// ([io.Closer] interface)
func (self ContextualReadCloser) Close() error {
	return self.reader.Close()
}

//
// TestLogWriter
//

type TestLogWriter struct {
	t *testing.T
}

func NewTestLogWriter(t *testing.T) *TestLogWriter {
	return &TestLogWriter{t}
}

// ([io.Writer] interface)
func (self *TestLogWriter) Write(p []byte) (n int, err error) {
	self.t.Helper()
	self.t.Log(p)
	return len(p), nil
}
