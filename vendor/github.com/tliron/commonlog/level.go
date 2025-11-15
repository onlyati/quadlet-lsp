package commonlog

import (
	"fmt"
)

//
// Level
//

type Level int

const (
	None     Level = 0
	Critical Level = 1
	Error    Level = 2
	Warning  Level = 3
	Notice   Level = 4
	Info     Level = 5
	Debug    Level = 6
)

// ([fmt.Stringify] interface)
func (self Level) String() string {
	switch self {
	case None:
		return "None"
	case Critical:
		return "Critical"
	case Error:
		return "Error"
	case Warning:
		return "Warning"
	case Notice:
		return "Notice"
	case Info:
		return "Info"
	case Debug:
		return "Debug"
	default:
		panic(fmt.Sprintf("unsupported log level: %d", self))
	}
}

// Translates a verbosity number to a maximum loggable level as
// follows:
//
//   - -4 and below: [None]
//   - -3: [Critical]
//   - -2: [Error]
//   - -1: [Warning]
//   - 0: [Notice]
//   - 1: [Info]
//   - 2 and above: [Debug]
func VerbosityToMaxLevel(verbosity int) Level {
	if verbosity < -4 {
		return None
	} else {
		switch verbosity {
		case -4:
			return None
		case -3:
			return Critical
		case -2:
			return Error
		case -1:
			return Warning
		case 0:
			return Notice
		case 1:
			return Info
		default:
			return Debug
		}
	}
}
