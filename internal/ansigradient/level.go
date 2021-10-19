package ansigradient

import (
	"github.com/jwalton/go-supportscolor"
)

// ColorLevel represents the ANSI color level supported by the terminal.
type ColorLevel = supportscolor.ColorLevel

const (
	// LevelNone represents a terminal that does not support color at all.
	LevelNone ColorLevel = supportscolor.None
	// LevelBasic represents a terminal with basic 16 color support.
	LevelBasic ColorLevel = supportscolor.Basic
	// LevelAnsi256 represents a terminal with 256 color support.
	LevelAnsi256 ColorLevel = supportscolor.Ansi256
	// LevelAnsi16m represents a terminal with full true color support.
	LevelAnsi16m ColorLevel = supportscolor.Ansi16m
)

var defaultLevel *ColorLevel = nil

// SetLevel is used to override the auto-detected color level.
func SetLevel(level ColorLevel) {
	l := level
	defaultLevel = &l
}

// GetLevel returns the currently configured color level.
func GetLevel() ColorLevel {
	if defaultLevel == nil {
		l := supportscolor.Stdout().Level
		defaultLevel = &l
	}
	return *defaultLevel
}
