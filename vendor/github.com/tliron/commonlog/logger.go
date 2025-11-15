package commonlog

//
// Logger
//

// While [NewMessage] is the "true" API entry point, this interface enables
// a familiar logger API. Because it's an interface, references can easily
// replace the implementation, for example setting a reference to
// [MOCK_LOGGER] will disable the logger.
//
// See [GetLogger].
type Logger interface {
	// Returns true if a level is loggable for this logger.
	AllowLevel(level Level) bool

	// Sets the maximum loggable level for this logger.
	SetMaxLevel(level Level)

	// Gets the maximum loggable level for this logger.
	GetMaxLevel() Level

	// Creates a new message for this logger. Will return nil if
	// the level is not loggable.
	//
	// The depth argument is used for skipping frames in callstack
	// logging, if supported.
	NewMessage(level Level, depth int, keysAndValues ...any) Message

	// Convenience method to create and send a message with at least
	// the "_message" key. Additional keys can be set by providing
	// a sequence of key-value pairs.
	Log(level Level, depth int, message string, keysAndValues ...any)

	// Convenience method to create and send a message with just
	// the "_message" key, where the message is created via the format
	// and args similarly to fmt.Printf.
	Logf(level Level, depth int, format string, args ...any)

	Critical(message string, keysAndValues ...any)
	Criticalf(format string, args ...any)
	Error(message string, keysAndValues ...any)
	Errorf(format string, args ...any)
	Warning(message string, keysAndValues ...any)
	Warningf(format string, args ...any)
	Notice(message string, keysAndValues ...any)
	Noticef(format string, args ...any)
	Info(message string, keysAndValues ...any)
	Infof(format string, args ...any)
	Debug(message string, keysAndValues ...any)
	Debugf(format string, args ...any)
}
