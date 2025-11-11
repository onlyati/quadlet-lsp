package commonlog

import (
	"fmt"
	"io"
	"strings"

	"github.com/tliron/kutil/terminal"
)

var backend Backend

// Sets the current backend.
//
// A nil backend will disable all logging
// (but the APIs would still not fail).
func SetBackend(backend_ Backend) {
	backend = backend_
}

// Configures the current backend. Verbosity is mapped to maximum
// loggable level as follows:
//
//   - -4 and below: [None]
//   - -3: [Critical]
//   - -2: [Error]
//   - -1: [Warning]
//   - 0: [Notice]
//   - 1: [Info]
//   - 2 and above: [Debug]
//
// Note that -4 ([None]) is a special case that is often optimized to turn
// off as much processing as possible.
//
// No-op if no backend was set.
func Configure(verbosity int, path *string) {
	if backend != nil {
		backend.Configure(verbosity, path)
	}
}

// Convenience method to call [Configure] while automatically overriding
// the verbosity with -4 ([None]) if [terminal.Quiet] is set to false
// and the path is empty (meaning we want to log to stdout).
func Initialize(verbosity int, path string) {
	if path == "" {
		if terminal.Quiet {
			verbosity = -4
		}
		Configure(verbosity, nil)
	} else {
		Configure(verbosity, &path)
	}
}

// Gets the current backend's [io.Writer]. Guaranteed to always return
// a valid non-nil value.
//
// Will be [io.Discard] if writing is unsupported by the backend or if
// no backend was set.
func GetWriter() io.Writer {
	if backend != nil {
		if writer := backend.GetWriter(); writer != nil {
			return writer
		}
	}

	return io.Discard
}

// Returns true if a level is loggable for the given name on the
// current backend.
//
// Returns false if no backend was set.
func AllowLevel(level Level, name ...string) bool {
	if backend != nil {
		return backend.AllowLevel(level, name...)
	} else {
		return false
	}
}

// Sets the maximum loggable level for the given name on the
// current backend. Will become the default maximum level for
// names deeper in the hierarchy unless explicitly set for
// them.
//
// No-op if no backend was set.
func SetMaxLevel(level Level, name ...string) {
	if backend != nil {
		backend.SetMaxLevel(level, name...)
	}
}

// Gets the maximum loggable level for the given name on the
// current backend.
//
// Returns [None] if no backend was set.
func GetMaxLevel(name ...string) Level {
	if backend != nil {
		return backend.GetMaxLevel(name...)
	} else {
		return None
	}
}

// Creates a new message for the given name on the current backend.
// Will return nil if the level is not loggable for the name, is
// [None], or if no backend was set.
//
// The depth argument is used for skipping frames in callstack
// logging, if supported.
func NewMessage(level Level, depth int, name ...string) Message {
	if (backend != nil) && (level != None) {
		return backend.NewMessage(level, depth+1, name...)
	} else {
		return nil
	}
}

// Calls [NewMessage] with [Critical] level.
func NewCriticalMessage(depth int, name ...string) Message {
	return NewMessage(Critical, depth+1, name...)
}

// Calls [NewMessage] with [Error] level.
func NewErrorMessage(depth int, name ...string) Message {
	return NewMessage(Error, depth+1, name...)
}

// Calls [NewMessage] with [Warning] level.
func NewWarningMessage(depth int, name ...string) Message {
	return NewMessage(Warning, depth+1, name...)
}

// Calls [NewMessage] with [Notice] level.
func NewNoticeMessage(depth int, name ...string) Message {
	return NewMessage(Notice, depth+1, name...)
}

// Calls [NewMessage] with [Info] level.
func NewInfoMessage(depth int, name ...string) Message {
	return NewMessage(Info, depth+1, name...)
}

// Calls [NewMessage] with [Debug] level.
func NewDebugMessage(depth int, name ...string) Message {
	return NewMessage(Debug, depth+1, name...)
}

// Provides a [BackendLogger] instance for the given path. The path
// is converted to a name using [PathToName].
//
// Note that this function will always return a new instance and that
// instances for the same path are functionally equivalent.
func GetLogger(path string) Logger {
	name := PathToName(path)
	return NewBackendLogger(name...)
}

// Calls [GetLogger] with [fmt.Sprintf] for the path.
func GetLoggerf(format string, values ...any) Logger {
	return GetLogger(fmt.Sprintf(format, values...))
}

// Convenience method to call a function and log the error, if
// returned, using [Logger.Error]. If task is not empty it will
// be prefixed to the error message.
func CallAndLogError(f func() error, task string, log Logger) {
	if err := f(); err != nil {
		if task != "" {
			log.Errorf("%s: %s", task, err.Error())
		} else {
			log.Error(err.Error())
		}
	}
}

// Convenience method to call a function and log the error, if
// returned, using [Logger.Warning]. If task is not empty it will
// be prefixed to the error message.
func CallAndLogWarning(f func() error, task string, log Logger) {
	if err := f(); err != nil {
		if task != "" {
			log.Warningf("%s: %s", task, err.Error())
		} else {
			log.Warning(err.Error())
		}
	}
}

func PathToName(path string) []string {
	name := strings.Split(path, ".")
	if len(name) == 0 {
		name = nil
	}
	return name
}
