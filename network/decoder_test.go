package network

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

func decoderFromBytes(b []byte) *PacketDecoder {
	return NewPacketDecoder(bufio.NewReader(bytes.NewReader(b)))
}

type varLongTestCase struct {
	// whether the test is supposed to fail
	invalid bool
	// whether this is a VarInt or VarLong
	long bool
	// the value
	val VarLong
	// the encoded value
	bytes []byte
}

// Obtained from https://wiki.vg/Protocol#VarInt_and_VarLong
var varLongTestCases = []varLongTestCase{
	// VarInts
	{val: 0, bytes: []byte{0x00}},
	{val: 1, bytes: []byte{0x01}},
	{val: 2, bytes: []byte{0x02}},
	{val: 127, bytes: []byte{0x7f}},
	{val: 128, bytes: []byte{0x80, 0x01}},
	{val: 255, bytes: []byte{0xff, 0x01}},
	{val: 25565, bytes: []byte{0xdd, 0xc7, 0x01}},
	{val: 2097151, bytes: []byte{0xff, 0xff, 0x7f}},
	{val: 2147483647, bytes: []byte{0xff, 0xff, 0xff, 0xff, 0x07}},
	{val: -1, bytes: []byte{0xff, 0xff, 0xff, 0xff, 0x0f}},
	{val: -2147483648, bytes: []byte{0x80, 0x80, 0x80, 0x80, 0x08}},
	// val is unused for this case
	{invalid: true, val: 0, bytes: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f}},

	// VarLongs
	{long: true, val: 9223372036854775807, bytes: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f}},
	{long: true, val: -1, bytes: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01}},
	{long: true, val: -2147483648, bytes: []byte{0x80, 0x80, 0x80, 0x80, 0xf8, 0xff, 0xff, 0xff, 0xff, 0x01}},
	{long: true, val: -9223372036854775808, bytes: []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x01}},
	// val is unused for this case
	{invalid: true, long: true, val: 0, bytes: []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x01}},
}

func TestDecoder_readVarLong(t *testing.T) {
	for _, c := range varLongTestCases {
		t.Run(fmt.Sprintf("%v", c), func(t *testing.T) {
			dec := decoderFromBytes(c.bytes)
			maxSize := varIntMaxBytes

			if c.long {
				maxSize = varLongMaxBytes
			}

			var v VarLong
			var n int

			if !dec.readVarLong(&v, maxSize, &n) {
				if c.invalid {
					return
				}

				t.Error(dec.LastError())
			}

			if c.invalid {
				t.Error("Expected error")
			}

			if c.long {
				if c.val != v {
					t.Errorf("Expected value %d, got %d", c.val, v)
				}
			} else {
				if VarInt(c.val) != VarInt(v) {
					t.Errorf("Expected value %d, got %d", c.val, v)
				}
			}

			if len(c.bytes) != n {
				t.Errorf("Expected size %d, got %d", len(c.bytes), n)
			}
		})
	}
}

func TestDecoder_ReadString_Valid(t *testing.T) {
	// TODO obtain more test cases
	decoder := decoderFromBytes([]byte{0x05, 0x6D, 0x61, 0x62, 0x6C, 0x65})
	var s string

	if !decoder.ReadString(&s) {
		t.Error(decoder.LastError())
	}

	if s != "mable" {
		t.Errorf("Expected 'mable', got '%s'", s)
	}
}

func TestDecoder_ReadString_TooBig(t *testing.T) {
	decoder := decoderFromBytes([]byte{0xff, 0xff, 0xff, 0xff, 0x07, 0x6D, 0x61, 0x62, 0x6C, 0x65})
	var s string

	if decoder.ReadString(&s) || decoder.LastError() != ErrStringTooBig {
		t.Errorf("Expected ErrStringTooBig, got %s", decoder.LastError())
	}
}

func TestDecoder_ReadString_NegLen(t *testing.T) {
	decoder := decoderFromBytes([]byte{0xff, 0xff, 0xff, 0xff, 0x0f, 0x6D, 0x61, 0x62, 0x6C, 0x65})
	var s string

	if decoder.ReadString(&s) || decoder.LastError() != ErrStringNegLen {
		t.Errorf("Expected ErrStringNegLen, got %s", decoder.LastError())
	}
}
