package protocol

import (
	"errors"
	"io"
)

type ReadBuffer struct {
	data []byte
}

func (r *ReadBuffer) ReadAll(src io.Reader, n int) error {
	// TODO reuse existing data if big enough
	r.data = make([]byte, n)
	_, err := io.ReadFull(src, r.data)
	return err
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
	// TODO currently both this file and varints.go implement reading varints, merge them in the future
	// https://github.com/Tnze/go-mc/blob/master/net/packet/types.go#L265
	var v uint32
	i := 0

	for b := byte(0x80); b&0x80 != 0; i++ {
		if i > 5 {
			panic(errors.New("VarInt too big"))
		}

		b = r.data[i]
		v |= uint32(b&0x7F) << uint32(7*i)
	}

	r.data = r.data[i:]
	return int(v)
}

func (r *ReadBuffer) ReadString() string {
	n := r.ReadVarInt()
	s := string(r.data[:n])
	r.data = r.data[n:]

	return s
}
