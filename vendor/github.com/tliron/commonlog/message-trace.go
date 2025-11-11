package commonlog

import (
	"runtime"
)

var Trace bool

// Adds "_file" and "_line" keys to a message if [Trace] is true.
// These are taken from the top of the callstack. Provide depth > 0
// to skip frames in the callstack.
func TraceMessage(message Message, depth int) Message {
	if Trace && (message != nil) {
		if _, file, line, ok := runtime.Caller(depth + 2); ok {
			message.Set(FILE, file)
			message.Set(LINE, line)
		}
	}

	return message
}
