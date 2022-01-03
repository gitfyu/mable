package log

import (
	"io"
	"os"
)

var (
	// Writer will receive the log messages, defaults to os.Stdout
	Writer io.Writer = os.Stdout

	// ErrorHandler will be invoked if an error occurs while logging a message. The default function will simply panic.
	ErrorHandler = func(err error) {
		panic(err)
	}
)

// Logger is used to write messages
type Logger struct {
	// Name is prepended to every message, convention is fully uppercase names
	Name string

	// MinLevel is the minimum new that will be logged, messages with a lower level will be dropped
	MinLevel Level
}

func (l *Logger) Trace(msg string) *Msg {
	return l.new(Trace, msg)
}

func (l *Logger) Debug(msg string) *Msg {
	return l.new(Debug, msg)
}

func (l *Logger) Info(msg string) *Msg {
	return l.new(Info, msg)
}

func (l *Logger) Warn(msg string) *Msg {
	return l.new(Warn, msg)
}

func (l *Logger) Error(msg string) *Msg {
	return l.new(Error, msg)
}

func (l *Logger) new(lvl Level, msg string) *Msg {
	if lvl < l.MinLevel {
		return nil
	}
	if len(l.Name) == 0 {
		l.Name = "DEFAULT"
	}
	return createMsg(lvl, l.Name, msg)
}
