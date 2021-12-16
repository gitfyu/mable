package protocol

import (
	"errors"
	"io"
	"math/bits"
)

type (
	VarInt  int32
	VarLong int64
)

const (
	VarIntMaxBytes  = 5
	VarLongMaxBytes = 10
)

var (
	errVarIntTooBig = errors.New("VarInt too big")
)

// Implementation is based on https://wiki.vg/Protocol#VarInt_and_VarLong and
// https://github.com/Tnze/go-mc/blob/master/net/packet/types.go

// VarIntSize returns the number of bytes required to write for the given value
func VarIntSize(v VarInt) int {
	return (31-bits.LeadingZeros32(uint32(v)))/7 + 1
}

// VarLongSize returns the number of bytes required to write for the given value
func VarLongSize(v VarLong) int {
	return (63-bits.LeadingZeros64(uint64(v)))/7 + 1
}

// ReadVarInt reads a single VarInt
func ReadVarInt(r io.ByteReader, v *VarInt) error {
	var tmp VarLong
	err := readVarIntOrLong(r, &tmp, VarIntMaxBytes)

	*v = VarInt(tmp)
	return err
}

// ReadVarLong reads a single VarLong
func ReadVarLong(r io.ByteReader, v *VarLong) error {
	return readVarIntOrLong(r, v, VarLongMaxBytes)
}

// WriteVarInt writes a single VarInt
func WriteVarInt(w io.ByteWriter, v VarInt) error {
	return writeVarIntOrLong(w, uint64(uint32(v)))
}

// WriteVarLong writes a single VarLong
func WriteVarLong(w io.ByteWriter, v VarLong) error {
	return writeVarIntOrLong(w, uint64(v))
}

func readVarIntOrLong(r io.ByteReader, v *VarLong, maxSize int) error {
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

	return nil
}

func writeVarIntOrLong(w io.ByteWriter, v uint64) error {
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
