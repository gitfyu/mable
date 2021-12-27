package protocol

import (
	"math"
)

type WriteBuffer struct {
	data []byte
	off  int
}

func (w *WriteBuffer) Reset() {
	w.off = 0
}

func (w *WriteBuffer) ensureSpace(n int) {
	if w.off+n > len(w.data) {
		newData := make([]byte, len(w.data)*2+n)
		copy(newData, w.data)
		w.data = newData
	}
}

func (w *WriteBuffer) WriteBool(v bool) {
	if v {
		w.WriteUint8(1)
	} else {
		w.WriteUint8(0)
	}
}

func (w *WriteBuffer) WriteUint8(v uint8) {
	w.ensureSpace(1)
	w.data[w.off] = v
	w.off++
}

func (w *WriteBuffer) WriteUint16(v uint16) {
	w.ensureSpace(2)
	w.data[w.off+0] = byte(v >> 8)
	w.data[w.off+1] = byte(v)
	w.off += 2
}

func (w *WriteBuffer) WriteUint32(v uint32) {
	w.ensureSpace(4)
	w.data[w.off+0] = byte(v >> 24)
	w.data[w.off+1] = byte(v >> 16)
	w.data[w.off+2] = byte(v >> 8)
	w.data[w.off+3] = byte(v)
	w.off += 4
}

func (w *WriteBuffer) WriteUint64(v uint64) {
	w.ensureSpace(8)
	w.data[w.off+0] = byte(v >> 56)
	w.data[w.off+1] = byte(v >> 48)
	w.data[w.off+2] = byte(v >> 40)
	w.data[w.off+3] = byte(v >> 32)
	w.data[w.off+4] = byte(v >> 24)
	w.data[w.off+5] = byte(v >> 16)
	w.data[w.off+6] = byte(v >> 8)
	w.data[w.off+7] = byte(v)
	w.off += 8
}

func (w *WriteBuffer) WriteVarInt(v int32) {
	size := VarIntSize(v)
	w.ensureSpace(size)
	WriteVarInt(w.data[w.off:], v)
	w.off += size
}

func (w *WriteBuffer) WriteFloat32(v float32) {
	w.WriteUint32(math.Float32bits(v))
}

func (w *WriteBuffer) WriteFloat64(v float64) {
	w.WriteUint64(math.Float64bits(v))
}

func (w *WriteBuffer) WriteBytes(b []byte) {
	w.ensureSpace(len(b))
	copy(w.data[w.off:], b)
	w.off += len(b)
}

func (w *WriteBuffer) WriteString(str string) {
	b := []byte(str)
	w.WriteVarInt(int32(len(b)))
	w.WriteBytes(b)
}

func (w *WriteBuffer) Len() int {
	return w.off
}

func (w *WriteBuffer) Bytes() []byte {
	return w.data[:w.off]
}
