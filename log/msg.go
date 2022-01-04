package log

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

const (
	// ErrorKey is the key that will be used for the error variable in Msg.Err.
	ErrorKey = "err"

	defaultMsgBufSize = 64
)

// Msg represents a log message. You should not construct a Msg yourself, instead use a Logger to create them. All
// exported functions can safely be called with a nil receiver, in which case they are no-ops.
type Msg struct {
	buf []byte
}

var msgPool = sync.Pool{
	New: func() interface{} {
		return &Msg{
			buf: make([]byte, 0, defaultMsgBufSize),
		}
	},
}

// Bool appends a variable of type bool.
func (m *Msg) Bool(key string, val bool) *Msg {
	if m != nil {
		if val {
			m.appendVar(key, "true")
		} else {
			m.appendVar(key, "false")
		}
	}

	return m
}

// Int appends a variable of type int64.
func (m *Msg) Int(key string, val int64) *Msg {
	if m != nil {
		m.appendVar(key, strconv.FormatInt(val, 10))
	}
	return m
}

// Uint appends a variable of type uint64.
func (m *Msg) Uint(key string, val uint64) *Msg {
	if m != nil {
		m.appendVar(key, strconv.FormatUint(val, 10))
	}
	return m
}

// F32 appends a variable of type float32.
func (m *Msg) F32(key string, val float32) *Msg {
	if m != nil {
		m.appendVar(key, strconv.FormatFloat(float64(val), 'f', -1, 32))
	}
	return m
}

// F64 appends a variable of type float64.
func (m *Msg) F64(key string, val float64) *Msg {
	if m != nil {
		m.appendVar(key, strconv.FormatFloat(val, 'f', -1, 64))
	}
	return m
}

// Str appends a variable of type string.
func (m *Msg) Str(key string, val string) *Msg {
	if m != nil {
		m.appendVar(key, val)
	}
	return m
}

// Stringer appends a variable whose value should be retrieved from a fmt.Stringer. While you could also manually
// convert the value to a string and then call Str, this function has the advantage that in the event that this message
// is dropped (if its level is below the minimum logging level), the string conversion is also skipped, which should
// result in better performance.
func (m *Msg) Stringer(key string, val fmt.Stringer) *Msg {
	if m != nil {
		m.appendVar(key, val.String())
	}
	return m
}

// Interface appends a variable whose value is an interface.
func (m *Msg) Interface(key string, val interface{}) *Msg {
	if m != nil {
		m.appendVar(key, fmt.Sprintf("%#v", val))
	}
	return m
}

// Err appends an error, using ErrorKey as the variable key.
func (m *Msg) Err(err error) *Msg {
	if m != nil {
		m.appendVar(ErrorKey, err.Error())
	}
	return m
}

// Log writes the Msg. You should not use this Msg again afterwards.
func (m *Msg) Log() {
	if m == nil {
		return
	}

	m.buf = append(m.buf, '\n')
	if _, err := Writer.Write(m.buf); err != nil {
		ErrorHandler(err)
	}

	msgPool.Put(m)
}

func (m *Msg) reset() {
	m.buf = m.buf[:0]
}

func (m *Msg) appendVar(key string, val string) {
	m.buf = append(m.buf, ' ', '[')
	m.buf = append(m.buf, key...)
	m.buf = append(m.buf, '=')
	m.buf = append(m.buf, val...)
	m.buf = append(m.buf, ']')
}

func createMsg(lvl Level, name string, msg string) *Msg {
	const timeFormat = "15:04:05"

	m := msgPool.Get().(*Msg)
	m.reset()

	m.buf = append(m.buf, time.Now().Format(timeFormat)...)
	m.buf = append(m.buf, ' ')
	m.buf = append(m.buf, lvl.String()...)
	m.buf = append(m.buf, ' ', '[')
	m.buf = append(m.buf, name...)
	m.buf = append(m.buf, ']', ' ')
	m.buf = append(m.buf, msg...)

	return m
}
