package commonlog

import (
	"fmt"

	"github.com/tliron/kutil/util"
)

//
// KeyValueLogger
//

// Wrapping [Logger] that calls [Message.Set] with keys and values
// on all messages.
//
// If we're wrapping another [KeyValueLogger] then our keys and values
// will be merged into the wrapped keys and values.
type KeyValueLogger struct {
	logger        Logger
	keysAndValues []any
}

func NewKeyValueLogger(logger Logger, keysAndValues ...any) KeyValueLogger {
	if keyValueLogger, ok := logger.(KeyValueLogger); ok {
		logger = keyValueLogger.logger
		keysAndValues, _ = MergeKeysAndValues(keyValueLogger.keysAndValues, keysAndValues)
	}

	return KeyValueLogger{
		logger:        logger,
		keysAndValues: keysAndValues,
	}
}

// Wrapping [Logger] that calls [Message.Set] with a "_scope" key
// on all messages.
//
// If we're wrapping another [KeyValueLogger] then our "_scope" key
// will be appended to the wrapped "_scope" key with a "." notation.
func NewScopeLogger(logger Logger, scope string) KeyValueLogger {
	var keysAndValues []any

	if keyValueLogger, ok := logger.(KeyValueLogger); ok {
		if scope_, ok := GetKeyValue(SCOPE, keyValueLogger.keysAndValues...); ok {
			if scope_ != "" {
				scope = util.ToString(scope_) + "." + scope
			}
		}

		logger = keyValueLogger.logger
		keysAndValues, _ = MergeKeysAndValues(keyValueLogger.keysAndValues, []any{SCOPE, scope})
	} else {
		keysAndValues = []any{SCOPE, scope}
	}

	return KeyValueLogger{
		logger:        logger,
		keysAndValues: keysAndValues,
	}
}

// ([Logger] interface)
func (self KeyValueLogger) AllowLevel(level Level) bool {
	return self.logger.AllowLevel(level)
}

// ([Logger] interface)
func (self KeyValueLogger) SetMaxLevel(level Level) {
	self.logger.SetMaxLevel(level)
}

// ([Logger] interface)
func (self KeyValueLogger) GetMaxLevel() Level {
	return self.logger.GetMaxLevel()
}

// ([Logger] interface)
func (self KeyValueLogger) NewMessage(level Level, depth int, keysAndValues ...any) Message {
	if message := self.logger.NewMessage(level, depth+1, self.keysAndValues...); message != nil {
		SetMessageKeysAndValues(message, keysAndValues...)
		return message
	} else {
		return nil
	}
}

// ([Logger] interface)
func (self KeyValueLogger) Log(level Level, depth int, message string, keysAndValues ...any) {
	if message_ := self.NewMessage(level, depth+1, keysAndValues...); message_ != nil {
		message_.Set(MESSAGE, message)
		message_.Send()
	}
}

// ([Logger] interface)
func (self KeyValueLogger) Logf(level Level, depth int, format string, args ...any) {
	if message := self.NewMessage(level, depth+1); message != nil {
		message.Set(MESSAGE, fmt.Sprintf(format, args...))
		message.Send()
	}
}

// ([Logger] interface)
func (self KeyValueLogger) Critical(message string, keysAndValues ...any) {
	self.Log(Critical, 1, message, keysAndValues...)
}

// ([Logger] interface)
func (self KeyValueLogger) Criticalf(format string, args ...any) {
	self.Logf(Critical, 1, format, args...)
}

// ([Logger] interface)
func (self KeyValueLogger) Error(message string, keysAndValues ...any) {
	self.Log(Error, 1, message, keysAndValues...)
}

// ([Logger] interface)
func (self KeyValueLogger) Errorf(format string, args ...any) {
	self.Logf(Error, 1, format, args...)
}

// ([Logger] interface)
func (self KeyValueLogger) Warning(message string, keysAndValues ...any) {
	self.Log(Warning, 1, message, keysAndValues...)
}

// ([Logger] interface)
func (self KeyValueLogger) Warningf(format string, args ...any) {
	self.Logf(Warning, 1, format, args...)
}

// ([Logger] interface)
func (self KeyValueLogger) Notice(message string, keysAndValues ...any) {
	self.Log(Notice, 1, message, keysAndValues...)
}

// ([Logger] interface)
func (self KeyValueLogger) Noticef(format string, args ...any) {
	self.Logf(Notice, 1, format, args...)
}

// ([Logger] interface)
func (self KeyValueLogger) Info(message string, keysAndValues ...any) {
	self.Log(Info, 1, message, keysAndValues...)
}

// ([Logger] interface)
func (self KeyValueLogger) Infof(format string, args ...any) {
	self.Logf(Info, 1, format, args...)
}

// ([Logger] interface)
func (self KeyValueLogger) Debug(message string, keysAndValues ...any) {
	self.Log(Debug, 1, message, keysAndValues...)
}

// ([Logger] interface)
func (self KeyValueLogger) Debugf(format string, args ...any) {
	self.Logf(Debug, 1, format, args...)
}
