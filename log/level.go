package log

import "strings"

// Level represents a logging level, with TraceLevel being the lowest and ErrorLevel the highest. The default Level (its
// zero value) is InfoLevel.
type Level int8

const (
	TraceLevel Level = iota - 2
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
)

var levelStrings = []string{
	"TRACE",
	"DEBUG",
	"INFO ",
	"WARN ",
	"ERROR",
}

// LevelFromString returns the Level corresponding to the given input string. If the input is not recognised, InfoLevel
// will be returned.
func LevelFromString(str string) Level {
	switch strings.ToLower(str) {
	case "trace":
		return TraceLevel
	case "debug":
		return DebugLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	default:
		fallthrough
	case "info":
		return InfoLevel
	}
}

// String returns the string representation of a Level. The returned string will always be 5 characters long, shorter
// level names are padded with a space at the end.
func (l Level) String() string {
	return levelStrings[l+2]
}
