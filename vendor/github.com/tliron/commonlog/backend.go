package commonlog

import (
	"io"
)

//
// Backend
//

type Backend interface {
	// Configures the backend. Verbosity is mapped to maximum
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
	Configure(verbosity int, path *string)

	// Gets the backend's [io.Writer]. Can be nil if unsupported.
	GetWriter() io.Writer

	NewMessage(level Level, depth int, name ...string) Message

	// Returns true if a level is loggable for the given name.
	AllowLevel(level Level, name ...string) bool

	// Sets the maximum loggable level for the given name. Will become the
	// default maximum level for names deeper in the hierarchy unless
	// explicitly set for them.
	SetMaxLevel(level Level, name ...string)

	// Gets the maximum loggable level for the given name.
	GetMaxLevel(name ...string) Level
}
