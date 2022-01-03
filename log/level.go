package log

import "strings"

// Level represents a logging level, with Trace being the lowest and Error the highest. The default Level (its zero
// value) is Info.
type Level int8

const (
	Trace Level = iota - 2
	Debug
	Info
	Warn
	Error
)

var levelStrings = []string{
	"TRACE",
	"DEBUG",
	"INFO ",
	"WARN ",
	"ERROR",
}

// LevelFromString returns the Level corresponding to the given input string. If the input is not recognised, Info
// will be returned.
func LevelFromString(str string) Level {
	switch strings.ToLower(str) {
	case "trace":
		return Trace
	case "debug":
		return Debug
	case "warn":
		return Warn
	case "error":
		return Error
	default:
		fallthrough
	case "info":
		return Info
	}
}

// String returns the string representation of a Level. The returned string will always be 5 characters long, shorter
// level names are padded with a space at the end.
func (l Level) String() string {
	return levelStrings[l+2]
}
