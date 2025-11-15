package commonlog

import (
	"fmt"
)

//
// BackendLogger
//

// Default [Logger] implementation that logs to the current backend set with
// [SetBackend].
type BackendLogger struct {
	name []string
}

func NewBackendLogger(name ...string) BackendLogger {
	return BackendLogger{name: name}
}

// ([Logger] interface)
func (self BackendLogger) AllowLevel(level Level) bool {
	return AllowLevel(level, self.name...)
}

// ([Logger] interface)
func (self BackendLogger) SetMaxLevel(level Level) {
	SetMaxLevel(level, self.name...)
}

// ([Logger] interface)
func (self BackendLogger) GetMaxLevel() Level {
	return GetMaxLevel(self.name...)
}

// ([Logger] interface)
func (self BackendLogger) NewMessage(level Level, depth int, keysAndValues ...any) Message {
	if message := NewMessage(level, depth+1, self.name...); message != nil {
		SetMessageKeysAndValues(message, keysAndValues...)
		return message
	} else {
		return nil
	}
}

// ([Logger] interface)
func (self BackendLogger) Log(level Level, depth int, message string, keysAndValues ...any) {
	if message_ := self.NewMessage(level, depth+1, keysAndValues...); message_ != nil {
		message_.Set(MESSAGE, message)
		message_.Send()
	}
}

// ([Logger] interface)
func (self BackendLogger) Logf(level Level, depth int, format string, args ...any) {
	if message := self.NewMessage(level, depth+1); message != nil {
		message.Set(MESSAGE, fmt.Sprintf(format, args...))
		message.Send()
	}
}

// ([Logger] interface)
func (self BackendLogger) Critical(message string, keysAndValues ...any) {
	self.Log(Critical, 1, message, keysAndValues...)
}

// ([Logger] interface)
func (self BackendLogger) Criticalf(format string, args ...any) {
	self.Logf(Critical, 1, format, args...)
}

// ([Logger] interface)
func (self BackendLogger) Error(message string, keysAndValues ...any) {
	self.Log(Error, 1, message, keysAndValues...)
}

// ([Logger] interface)
func (self BackendLogger) Errorf(format string, args ...any) {
	self.Logf(Error, 1, format, args...)
}

// ([Logger] interface)
func (self BackendLogger) Warning(message string, keysAndValues ...any) {
	self.Log(Warning, 1, message, keysAndValues...)
}

// ([Logger] interface)
func (self BackendLogger) Warningf(format string, args ...any) {
	self.Logf(Warning, 1, format, args...)
}

// ([Logger] interface)
func (self BackendLogger) Notice(message string, keysAndValues ...any) {
	self.Log(Notice, 1, message, keysAndValues...)
}

// ([Logger] interface)
func (self BackendLogger) Noticef(format string, args ...any) {
	self.Logf(Notice, 1, format, args...)
}

// ([Logger] interface)
func (self BackendLogger) Info(message string, keysAndValues ...any) {
	self.Log(Info, 1, message, keysAndValues...)
}

// ([Logger] interface)
func (self BackendLogger) Infof(format string, args ...any) {
	self.Logf(Info, 1, format, args...)
}

// ([Logger] interface)
func (self BackendLogger) Debug(message string, keysAndValues ...any) {
	self.Log(Debug, 1, message, keysAndValues...)
}

// ([Logger] interface)
func (self BackendLogger) Debugf(format string, args ...any) {
	self.Logf(Debug, 1, format, args...)
}
