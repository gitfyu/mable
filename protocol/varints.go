package protocol

import (
	"errors"
	"io"
	"math/bits"
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
func VarIntSize(v int32) int {
	return (31-bits.LeadingZeros32(uint32(v)))/7 + 1
}

func ReadVarInt(r io.ByteReader) (int32, error) {
	var v uint32
	var n int

	for {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}

		v |= uint32(b&0x7F) << (n * 7)
		n++
		if n > VarIntMaxBytes {
			return 0, errVarIntTooBig
		}

		if (b & 0x80) == 0 {
			break
		}
	}

	return int32(v), nil
}

// WriteVarInt writes a VarInt to a byte slice. If the slice is too small, this function will panic. You can check the
// required size using VarIntSize.
func WriteVarInt(buf []byte, v int32) {
	uv := uint32(v)
	for i := 0; uv != 0; i++ {
		b := uv & 0x7F
		uv >>= 7
		if uv != 0 {
			b |= 0x80
		}

		buf[i] = byte(b)
	}
}
