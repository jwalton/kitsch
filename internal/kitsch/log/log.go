// Package log provides logging for kitsch prompt.
package log

import (
	"fmt"

	"github.com/jwalton/gchalk"
)

// If verbose is false, we'll only print the first warning that comes along, and
// hide all "info" messages.
var verbose = false

var warningShowed = false

// SetVerbose sets verbose logging.
//
// In non-verbose mode, most logging is hidden (we're trying to show a prompt
// here, so we don't want to bombard the user with messages).  "Info"s are hidden,
// the first "Warn" will be displayed but the rest will be hidden.  "Error"s
// are shown but should be rare.
func SetVerbose(v bool) {
	verbose = v
}

// Info prints an info-level message to stderr.
func Info(message ...interface{}) {
	if verbose {
		print(gchalk.Stderr.BrightCyan("Info: "))
		println(fmt.Sprint(message...))
	}
}

// Warn prints a warn-level message to stderr.  If non-verbose mode, only the
// first warning will be displayed.  Once the user fixes that warning, or the
// user runs `check` or in verbose mode, we can show them more warnings.
func Warn(message ...interface{}) {
	if verbose || !warningShowed {
		warningShowed = true
		print(gchalk.Stderr.BrightYellow("Warn: "))
		println(fmt.Sprint(message...))
	}
}

// Error prints an error message to stderr.
func Error(message ...interface{}) {
	print(gchalk.Stderr.BrightRed("Err : "))
	println(fmt.Sprint(message...))
}
