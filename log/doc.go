/*
Package log is designed as a replacement for https:github.com/rs/zerolog to allow efficiently logging in a
human-friendly format with a similar, simple API.

Example usage:
	logger := Logger{Name: "MAIN"}
	logger.InfoLevel("Example message").Int("num", 123).Log()

You can also adjust several global settings that apply to every Logger:
	log.Writer = someWriter
	log.ErrorHandler = func(err error) {
		// do something
	}
	log.ErrorKey = "error"
Note: these global settings should be set before any Logger is used, they cannot be called concurrently!

Logger.New (and all shortcut functions that call it, such as Logger.Info) may safely be called concurrently, however
the returned Msg should not be used concurrently.
*/
package log
