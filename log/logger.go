// Package log is designed as a replacement for https://github.com/rs/zerolog to allow efficiently logging in a
// human-friendly format with a similar, simple API.
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

	// MinLevel is the minimum New that will be logged, messages with a lower level will be dropped
	MinLevel Level
}

// Trace is a shortcut for New(Trace, msg)
func (l *Logger) Trace(msg string) *Msg {
	return l.New(Trace, msg)
}

// Debug is a shortcut for New(Debug, msg)
func (l *Logger) Debug(msg string) *Msg {
	return l.New(Debug, msg)
}

// Info is a shortcut for New(Info, msg)
func (l *Logger) Info(msg string) *Msg {
	return l.New(Info, msg)
}

// Warn is a shortcut for New(Warn, msg)
func (l *Logger) Warn(msg string) *Msg {
	return l.New(Warn, msg)
}

// Error is a shortcut for New(Error, msg)
func (l *Logger) Error(msg string) *Msg {
	return l.New(Error, msg)
}

// New creates a new Msg with the specified Level. If lvl is below MinLevel, this function does nothing and returns nil.
// The message will not be written until Msg.Log is called.
func (l *Logger) New(lvl Level, msg string) *Msg {
	if lvl < l.MinLevel {
		return nil
	}
	return createMsg(lvl, l.Name, msg)
}
