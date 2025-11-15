package util

import (
	"github.com/tliron/kutil/terminal"
)

// Attempts to enable colorization for [os.Stdout] and [os.Stderr] according to
// the colorize argument. If cleanup is required, will set register it with
// [OnExit]. Errors will [Fail].
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
// Other colorize values will [Fail].
//
// See [terminal.InitializeColorization].
func InitializeColorization(colorize string) {
	cleanupStdout, cleanupStderr, err := terminal.InitializeColorization(colorize)
	FailOnError(err)
	if cleanupStdout != nil {
		OnExitError(cleanupStdout)
	}
	if cleanupStderr != nil {
		OnExitError(cleanupStderr)
	}
}
