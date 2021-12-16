package protocol

import (
	"errors"
	"io"
	"math/bits"
)

type (
	VarInt int32
)

const (
	VarIntMaxBytes = 5
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

func ReadVarInt(r io.ByteReader, v *VarInt) error {
	var n int

	for {
		b, err := r.ReadByte()
		if err != nil {
			return err
		}

		*v |= VarInt(b&0x7F) << (n * 7)
		n++
		if n > VarIntMaxBytes {
			return errVarIntTooBig
		}

		if (b & 0x80) == 0 {
			break
		}
	}

	return nil
}

func WriteVarInt(w io.ByteWriter, v VarInt) error {
	val := uint32(v)

	for {
		b := val & 0x7F
		val >>= 7
		if val != 0 {
			b |= 0x80
		}

		if err := w.WriteByte(byte(b)); err != nil {
			return err
		}

		if val == 0 {
			return nil
		}
	}
}
