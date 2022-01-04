package protocol

import (
	"io"
	"math"
)

// ReadBuffer is a utility for reading data types that are common in packets from a byte slice. Most of its functions
// do not return an error, instead they will panic.
type ReadBuffer struct {
	data []byte
}

// ReadAll reads exactly n bytes from an io.Reader.
func (r *ReadBuffer) ReadAll(src io.Reader, n int) error {
	if cap(r.data) >= n {
		r.data = r.data[:n]
	} else {
		r.data = make([]byte, n)
	}
	_, err := io.ReadFull(src, r.data)
	return err
}

func (r *ReadBuffer) ReadBool() bool {
	return r.ReadUint8() == 1
}

func (r *ReadBuffer) ReadUint8() uint8 {
	v := r.data[0]
	r.data = r.data[1:]
	return v
}

func (r *ReadBuffer) ReadUint16() uint16 {
	v := uint16(r.data[1]) | uint16(r.data[0])<<8
	r.data = r.data[2:]
	return v
}

func (r *ReadBuffer) ReadUint32() uint32 {
	v := uint32(r.data[3]) | uint32(r.data[2])<<8 | uint32(r.data[1])<<16 | uint32(r.data[0])<<24
	r.data = r.data[4:]
	return v
}

func (r *ReadBuffer) ReadUint64() uint64 {
	v := uint64(r.data[7]) |
		uint64(r.data[6])<<8 |
		uint64(r.data[5])<<16 |
		uint64(r.data[4])<<24 |
		uint64(r.data[3])<<32 |
		uint64(r.data[2])<<40 |
		uint64(r.data[1])<<48 |
		uint64(r.data[0])<<56
	r.data = r.data[8:]
	return v
}

func (r *ReadBuffer) ReadVarInt() int {
	v, _ := ReadVarInt(r)
	return int(v)
}

func (r *ReadBuffer) ReadFloat32() float32 {
	return math.Float32frombits(r.ReadUint32())
}

func (r *ReadBuffer) ReadFloat64() float64 {
	return math.Float64frombits(r.ReadUint64())
}

func (r *ReadBuffer) ReadString() string {
	n := r.ReadVarInt()
	s := string(r.data[:n])
	r.data = r.data[n:]

	return s
}

func (r *ReadBuffer) ReadByte() (byte, error) {
	b := r.data[0]
	r.data = r.data[1:]
	return b, nil
}
