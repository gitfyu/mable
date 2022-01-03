package log

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

const (
	// ErrorKey is the key that will be used for the error variable in Msg.Err
	ErrorKey = "err"
)

// Msg represents a log message. You should not construct a Msg yourself, instead use a Logger to create them. All
// exported functions can safely be called with a nil receiver, in which case they are no-ops.
type Msg struct {
	buf []byte
	len int
}

var msgPool = sync.Pool{
	New: func() interface{} {
		return new(Msg)
	},
}

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

func (m *Msg) Int(key string, val int64) *Msg {
	if m != nil {
		m.appendVar(key, strconv.FormatInt(val, 10))
	}
	return m
}

func (m *Msg) Uint(key string, val uint64) *Msg {
	if m != nil {
		m.appendVar(key, strconv.FormatUint(val, 10))
	}
	return m
}

func (m *Msg) F32(key string, val float32) *Msg {
	if m != nil {
		m.appendVar(key, strconv.FormatFloat(float64(val), 'f', -1, 32))
	}
	return m
}

func (m *Msg) F64(key string, val float64) *Msg {
	if m != nil {
		m.appendVar(key, strconv.FormatFloat(val, 'f', -1, 64))
	}
	return m
}

func (m *Msg) Str(key string, val string) *Msg {
	if m != nil {
		m.appendVar(key, val)
	}
	return m
}

func (m *Msg) Stringer(key string, val fmt.Stringer) *Msg {
	if m != nil {
		m.appendVar(key, val.String())
	}
	return m
}

func (m *Msg) Interface(key string, val interface{}) *Msg {
	if m != nil {
		m.appendVar(key, fmt.Sprintf("%#v", val))
	}
	return m
}

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

	m.reserve(1)
	m.appendByte('\n')

	if _, err := Writer.Write(m.buf[:m.len]); err != nil {
		ErrorHandler(err)
	}

	msgPool.Put(m)
}

func (m *Msg) appendByte(b byte) {
	m.buf[m.len] = b
	m.len++
}

func (m *Msg) append(b string) {
	copy(m.buf[m.len:], b)
	m.len += len(b)
}

func (m *Msg) appendVar(key string, val string) {
	// len(" [" + key + "=" + val + "]")
	m.reserve(2 + len(key) + 1 + len(val) + 1)

	m.append(" [")
	m.append(key)
	m.appendByte('=')
	m.append(val)
	m.appendByte(']')
}

func createMsg(lvl Level, name string, msg string) *Msg {
	const timeFormat = "15:04"

	m := msgPool.Get().(*Msg)
	m.len = 0

	// len("12:34:56 ERR [" + name + "] " + msg)
	m.reserve(len(timeFormat) + 1 + len(lvl.str()) + 2 + len(name) + 2 + len(msg))

	m.append(time.Now().Format(timeFormat))
	m.appendByte(' ')
	m.append(lvl.str())
	m.append(" [")
	m.append(name)
	m.append("] ")
	m.append(msg)

	return m
}

func (m *Msg) reserve(n int) {
	if m.len+n > len(m.buf) {
		newBuf := make([]byte, len(m.buf)*2+n)
		copy(newBuf, m.buf)
		m.buf = newBuf
	}
}
