package protocol

import (
	"encoding/json"
	"github.com/gitfyu/mable/chat"
	"math"
	"sync"
)

const (
	defaultWriteBufSize = 32
)

// WriteBuffer is a utility for writing data types commonly used in packets. Afterwards, its contents can be converted
// to a byte slice using WriteBuffer.Bytes.
type WriteBuffer struct {
	data      []byte
	varIntBuf [VarIntMaxBytes]byte
}

var writeBufPool = sync.Pool{
	New: func() interface{} {
		return &WriteBuffer{
			data: make([]byte, 0, defaultWriteBufSize),
		}
	},
}

// AcquireWriteBuffer obtains a new WriteBuffer, which should later be released using ReleaseWriteBuffer.
func AcquireWriteBuffer() *WriteBuffer {
	return writeBufPool.Get().(*WriteBuffer)
}

// ReleaseWriteBuffer releases resources used for a WriteBuffer. You should call this after you are done using a buffer,
// after which you should not use the buffer again.
func ReleaseWriteBuffer(w *WriteBuffer) {
	writeBufPool.Put(w)
}

// Reset resets the buffers contents.
func (w *WriteBuffer) Reset() {
	w.data = w.data[:0]
}

func (w *WriteBuffer) WriteBool(v bool) {
	var b uint8
	if v {
		b = 1
	} else {
		b = 0
	}
	w.data = append(w.data, b)
}

func (w *WriteBuffer) WriteUint8(v uint8) {
	w.data = append(w.data, v)
}

func (w *WriteBuffer) WriteUint16(v uint16) {
	w.data = append(w.data, byte(v>>8), byte(v))
}

func (w *WriteBuffer) WriteUint32(v uint32) {
	w.data = append(w.data,
		byte(v>>24),
		byte(v>>16),
		byte(v>>8),
		byte(v),
	)
}

func (w *WriteBuffer) WriteUint64(v uint64) {
	w.data = append(w.data,
		byte(v>>56),
		byte(v>>48),
		byte(v>>40),
		byte(v>>32),
		byte(v>>24),
		byte(v>>16),
		byte(v>>8),
		byte(v),
	)
}

func (w *WriteBuffer) WriteVarInt(v int32) {
	WriteVarInt(w.varIntBuf[:], v)
	w.data = append(w.data, w.varIntBuf[:VarIntSize(v)]...)
}

func (w *WriteBuffer) WriteFloat32(v float32) {
	w.WriteUint32(math.Float32bits(v))
}

func (w *WriteBuffer) WriteFloat64(v float64) {
	w.WriteUint64(math.Float64bits(v))
}

func (w *WriteBuffer) WriteBytes(b []byte) {
	w.data = append(w.data, b...)
}

func (w *WriteBuffer) WriteString(str string) {
	w.WriteVarInt(int32(len(str)))
	w.data = append(w.data, str...)
}

func (w *WriteBuffer) WriteByteArrayWithLength(b []byte) {
	w.WriteVarInt(int32(len(b)))
	w.data = append(w.data, b...)
}

func (w *WriteBuffer) WriteChat(msg *chat.Msg) {
	str, err := json.Marshal(msg)
	if err != nil {
		// should never happen
		panic(err)
	}

	w.WriteByteArrayWithLength(str)
}

// Len returns the number of bytes written to the buffer so far.
func (w *WriteBuffer) Len() int {
	return len(w.data)
}

// Bytes returns a view of the contents that have been written so far. The returned slice should not be modified
// directly and is only valid until the next write to the buffer.
func (w *WriteBuffer) Bytes() []byte {
	return w.data
}
