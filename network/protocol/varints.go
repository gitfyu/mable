package protocol

import (
	"errors"
)

type (
	VarInt  int32
	VarLong int64
)

const (
	varIntMaxBytes  = 5
	varLongMaxBytes = 10
)

var (
	errVarIntTooBig = errors.New("VarInt too big")
)

// Implementation is based on https://wiki.vg/Protocol#VarInt_and_VarLong and
// https://github.com/Tnze/go-mc/blob/master/net/packet/types.go

type ByteReader interface {
	ReadByte() (byte, error)
}

type ByteWriter interface {
	WriteByte(byte) error
}

// ReadVarInt reads a single VarInt. The size parameter receives the number of bytes read, unless it is nil.
func ReadVarInt(r ByteReader, v *VarInt, size *int) error {
	var tmp VarLong
	err := readVarIntOrLong(r, &tmp, varIntMaxBytes, size)

	*v = VarInt(tmp)
	return err
}

// ReadVarLong reads a single VarLong. The size parameter receives the number of bytes read, unless it is nil.
func ReadVarLong(r ByteReader, v *VarLong, size *int) error {
	return readVarIntOrLong(r, v, varLongMaxBytes, size)
}

// WriteVarInt writes a single VarInt
func WriteVarInt(w ByteWriter, v VarInt) error {
	return writeVarIntOrLong(w, uint64(uint32(v)))
}

// WriteVarLong writes a single VarLong
func WriteVarLong(w ByteWriter, v VarLong) error {
	return writeVarIntOrLong(w, uint64(v))
}

func readVarIntOrLong(r ByteReader, v *VarLong, maxSize int, size *int) error {
	var n int

	for {
		b, err := r.ReadByte()
		if err != nil {
			return err
		}

		*v |= VarLong(b&0x7F) << (n * 7)
		n++
		if n > maxSize {
			return errVarIntTooBig
		}

		if (b & 0x80) == 0 {
			break
		}
	}

	if size != nil {
		*size = n
	}

	return nil
}

func writeVarIntOrLong(w ByteWriter, v uint64) error {
	for {
		b := v & 0x7F
		v >>= 7
		if v != 0 {
			b |= 0x80
		}

		if err := w.WriteByte(byte(b)); err != nil {
			return err
		}

		if v == 0 {
			return nil
		}
	}
}
