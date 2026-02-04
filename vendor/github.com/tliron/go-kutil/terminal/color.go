package terminal

import (
	"fmt"
	"os"
	"strconv"
)

var ColorizeStdout bool
var ColorizeStderr bool

type CleanupFunc func() error

const (
	escapePrefix = "\x1b["
	escapeSuffix = "m"

	ResetCode   = escapePrefix + "0" + escapeSuffix
	RedCode     = escapePrefix + "31" + escapeSuffix
	GreenCode   = escapePrefix + "32" + escapeSuffix
	YellowCode  = escapePrefix + "33" + escapeSuffix
	BlueCode    = escapePrefix + "34" + escapeSuffix
	MagentaCode = escapePrefix + "35" + escapeSuffix
	CyanCode    = escapePrefix + "36" + escapeSuffix
	GrayCode    = escapePrefix + "90" + escapeSuffix
)

// Attempts to enable colorization for [os.Stdout] and [os.Stderr] according to
// the colorize argument. If it succeeds it may return two [CleanupFunc]s that
// should be called to restore [os.Stdout] and [os.Stderr] to their original mode.
// The [CleanupFunc]s may also be nil.
//
// The colorize argument can be:
//
//   - "true", "TRUE", "True", "t", "T", "1": Attempts to enable
//     colorization if [os.Stdout] and [os.Stdout] support it.
//     If it succeeds will set [ColorizeStdout], [StdoutStylist],
//     [ColorizeStderr], and [StderrStylist] accordingly.
//   - "false", "FALSE", "False", "f", "F", "0": Does nothing.
//   - "force": Sets [ColorizeStdout], [StdoutStylist],
//     [ColorizeStderr], and [StderrStylist] as if coloriziation
//     were enabled.
//
// Other colorize values will return an error.
func InitializeColorization(colorize string) (CleanupFunc, CleanupFunc, error) {
	if colorize == "force" {
		ColorizeStdout = true
		StdoutStylist = NewStylist(true)
		ColorizeStderr = true
		StderrStylist = StdoutStylist
	} else if colorize_, err := strconv.ParseBool(colorize); err == nil {
		if colorize_ {
			var cleanupStdout CleanupFunc
			var cleanupStderr CleanupFunc
			var ok bool
			var err error

			if ok, cleanupStdout, err = EnableColor(os.Stdout); err == nil {
				if ok {
					ColorizeStdout = true
					StdoutStylist = NewStylist(true)
				}
			} else {
				return nil, nil, err
			}

			if ok, cleanupStderr, err = EnableColor(os.Stderr); err == nil {
				if ok {
					ColorizeStderr = true
					StderrStylist = NewStylist(true)
				}
			} else {
				if cleanupStdout != nil {
					cleanupStdout()
				}
				return nil, nil, err
			}

			return cleanupStdout, cleanupStderr, nil
		}
	} else {
		return nil, nil, fmt.Errorf("\"--colorize\" must be \"true\", \"false\", or \"force\": %s", colorize)
	}

	return nil, nil, nil
}

//
// Colorizer
//

type Colorizer func(name string) string

// ([Colorizer] signature)
func ColorRed(s string) string {
	return RedCode + s + ResetCode
}

// ([Colorizer] signature)
func ColorGreen(s string) string {
	return GreenCode + s + ResetCode
}

// ([Colorizer] signature)
func ColorYellow(s string) string {
	return YellowCode + s + ResetCode
}

// ([Colorizer] signature)
func ColorBlue(s string) string {
	return BlueCode + s + ResetCode
}

// ([Colorizer] signature)
func ColorMagenta(s string) string {
	return MagentaCode + s + ResetCode
}

// ([Colorizer] signature)
func ColorCyan(s string) string {
	return CyanCode + s + ResetCode
}

// ([Colorizer] signature)
func ColorGray(s string) string {
	return GrayCode + s + ResetCode
}
