package network

import (
	"bufio"
	"bytes"
	"testing"
)

func decoderFromBytes(b []byte) *Decoder {
	return &Decoder{
		Reader: bufio.NewReader(bytes.NewReader(b)),
	}
}

// Test cases are obtained from https://wiki.vg/Protocol#VarInt_and_VarLong

var varInts = []VarInt{
	0, 1, 2, 127, 128, 255, 25565, 2097151, 2147483647, -1, -2147483648,
}

var varIntBytes = []byte{
	0x00,
	0x01,
	0x02,
	0x7f,
	0x80, 0x01,
	0xff, 0x01,
	0xdd, 0xc7, 0x01,
	0xff, 0xff, 0x7f,
	0xff, 0xff, 0xff, 0xff, 0x07,
	0xff, 0xff, 0xff, 0xff, 0x0f,
	0x80, 0x80, 0x80, 0x80, 0x08,
}

var varIntBytesTooBig = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0x0f}

func TestDecoder_ReadVarInt_Valid(t *testing.T) {
	decoder := decoderFromBytes(varIntBytes)

	for _, v := range varInts {
		var result VarInt

		if !decoder.ReadVarInt(&result) {
			t.Error(decoder.LastError())
		}

		if result != v {
			t.Errorf("Expected %d, got %d", v, result)
		}
	}
}

func TestDecoder_ReadVarInt_TooBig(t *testing.T) {
	decoder := decoderFromBytes(varIntBytesTooBig)
	var v VarInt

	if decoder.ReadVarInt(&v) || decoder.LastError() != ErrVarIntTooBig {
		t.Error("Expected ErrVarIntTooBig")
	}
}

var varLongs = []VarLong{
	0, 1, 2, 127, 128, 255, 2147483647, 9223372036854775807, -1, -2147483648, -9223372036854775808,
}

var varLongBytes = []byte{
	0x00,
	0x01,
	0x02,
	0x7f,
	0x80, 0x01,
	0xff, 0x01,
	0xff, 0xff, 0xff, 0xff, 0x07,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01,
	0x80, 0x80, 0x80, 0x80, 0xf8, 0xff, 0xff, 0xff, 0xff, 0x01,
	0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x01,
}

var varLongBytesTooBig = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01}

func TestDecoder_ReadVarLong_Valid(t *testing.T) {
	decoder := decoderFromBytes(varLongBytes)

	for _, v := range varLongs {
		var result VarLong

		if !decoder.ReadVarLong(&result) {
			t.Error(decoder.LastError())
		}

		if result != v {
			t.Errorf("Expected %d, got %d", v, result)
		}
	}
}

func TestDecoder_ReadVarLong_TooBig(t *testing.T) {
	decoder := decoderFromBytes(varLongBytesTooBig)
	var v VarLong

	if decoder.ReadVarLong(&v) || decoder.LastError() != ErrVarLongTooBig {
		t.Error("Expected ErrVarLongTooBig")
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
