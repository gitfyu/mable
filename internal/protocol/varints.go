package protocol

import (
	"errors"
	"io"
	"math/bits"
)

// https://wiki.vg/VarInt_And_VarLong

const (
	// VarIntMaxBytes is the maximum number of bytes required to represent a single VarInt.
	VarIntMaxBytes = 5

	segmentBits = 0x7F
	continueBit = 0x80
)

var (
	errVarIntTooBig = errors.New("VarInt too big")
)

// VarIntSize returns the number of bytes required to write the given value as a VarInt.
func VarIntSize(v int32) int {
	return (31-bits.LeadingZeros32(uint32(v)))/7 + 1
}

// ReadVarInt reads a single VarInt from an io.ByteReader.
func ReadVarInt(r io.ByteReader) (int32, error) {
	var v, pos int32

	for {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}

		v |= int32(b&segmentBits) << pos
		if b&continueBit == 0 {
			break
		}

		pos += 7
		if pos >= 32 {
			return 0, errVarIntTooBig
		}
	}

	return v, nil
}

// WriteVarInt writes a single VarInt to io.ByteWriter.
func WriteVarInt(w io.ByteWriter, v int32) error {
	for {
		if v&(^segmentBits) == 0 {
			return w.WriteByte(byte(v))
		}

		err := w.WriteByte(byte(v&segmentBits) | continueBit)
		if err != nil {
			return err
		}

		v = int32(uint32(v) >> 7)
	}
}
