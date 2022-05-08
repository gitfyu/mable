package protocol

import (
	"encoding/json"
	"io"
	"math"

	"github.com/gitfyu/mable/chat"
)

// Reader combines all interfaces needed to be able to
// read any datatype in the Minecraft protocol.
type Reader interface {
	io.Reader
	io.ByteReader
}

// Writer combines all interfaces needed to be able to
// write any datatype in the Minecraft protocol.
type Writer interface {
	io.Writer
	io.ByteWriter
	io.StringWriter
}

func ReadBool(r io.ByteReader) (bool, error) {
	b, err := r.ReadByte()
	return b != 0, err
}

func ReadUint16(r io.Reader) (uint16, error) {
	var b [2]byte
	_, err := r.Read(b[:])
	return uint16(b[1]) | uint16(b[0])<<8, err
}

func ReadUint32(r io.Reader) (uint32, error) {
	var b [4]byte
	_, err := r.Read(b[:])
	return uint32(b[3]) |
		uint32(b[2])<<8 |
		uint32(b[1])<<16 |
		uint32(b[0])<<24, err
}

func ReadUint64(r io.Reader) (uint64, error) {
	var b [8]byte
	_, err := r.Read(b[:])
	return uint64(b[7]) |
		uint64(b[6])<<8 |
		uint64(b[5])<<16 |
		uint64(b[4])<<24 |
		uint64(b[3])<<32 |
		uint64(b[2])<<40 |
		uint64(b[1])<<48 |
		uint64(b[0])<<56, err
}

func ReadFloat32(r io.Reader) (float32, error) {
	v, err := ReadUint32(r)
	return math.Float32frombits(v), err
}

func ReadFloat64(r io.Reader) (float64, error) {
	v, err := ReadUint64(r)
	return math.Float64frombits(v), err
}

func ReadString(r Reader) (string, error) {
	// TODO cap maximum length
	len, err := ReadVarInt(r)
	if err != nil {
		return "", err
	}

	b := make([]byte, len)
	if _, err = r.Read(b); err != nil {
		return "", err
	}

	return string(b), nil
}

func WriteBool(w io.ByteWriter, v bool) error {
	var err error
	if v {
		err = w.WriteByte(1)
	} else {
		err = w.WriteByte(0)
	}
	return err
}

func WriteUint16(w io.Writer, v uint16) error {
	b := [2]byte{
		byte(v >> 8),
		byte(v),
	}
	_, err := w.Write(b[:])
	return err
}

func WriteUint32(w io.Writer, v uint32) error {
	b := [4]byte{
		byte(v >> 24),
		byte(v >> 16),
		byte(v >> 8),
		byte(v),
	}
	_, err := w.Write(b[:])
	return err
}

func WriteUint64(w io.Writer, v uint64) error {
	b := [8]byte{
		byte(v >> 56),
		byte(v >> 48),
		byte(v >> 40),
		byte(v >> 32),
		byte(v >> 24),
		byte(v >> 16),
		byte(v >> 8),
		byte(v),
	}
	_, err := w.Write(b[:])
	return err
}

func WriteFloat32(w io.Writer, v float32) error {
	return WriteUint32(w, math.Float32bits(v))
}

func WriteFloat64(w io.Writer, v float64) error {
	return WriteUint64(w, math.Float64bits(v))
}

func WriteString(w Writer, s string) error {
	if err := WriteVarInt(w, int32(len(s))); err != nil {
		return err
	}
	_, err := w.WriteString(s)
	return err
}

func WriteByteArray(w Writer, b []byte) error {
	if err := WriteVarInt(w, int32(len(b))); err != nil {
		return err
	}
	_, err := w.Write(b)
	return err
}

func WriteChat(w Writer, v *chat.Msg) error {
	str, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return WriteByteArray(w, str)
}
