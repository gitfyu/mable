package log

import "strings"

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

func (l Level) str() string {
	return levelStrings[l+2]
}
