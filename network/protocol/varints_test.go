package protocol

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

type varIntTest struct {
	// the value
	val VarInt
	// the encoded value
	bytes []byte
}

// Obtained from https://wiki.vg/Protocol#VarInt_and_VarLong
var varIntTests = []varIntTest{
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
}

var varIntInvalidTest = varIntTest{
	val: 0, bytes: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f},
}

func Test_ReadVarInt(t *testing.T) {
	for _, test := range varIntTests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			r := bufio.NewReader(bytes.NewReader(test.bytes))
			var val VarInt

			if err := ReadVarInt(r, &val); err != nil {
				t.Error(err)
			}
			if test.val != val {
				t.Errorf("Expected val %d, got %d", test.val, val)
			}
		})
	}
}

func Test_ReadVarInt_Invalid(t *testing.T) {
	r := bufio.NewReader(bytes.NewReader(varIntInvalidTest.bytes))
	var val VarInt
	if err := ReadVarInt(r, &val); err == nil {
		t.Error("Expected error")
	}
}

func Test_WriteVarInt(t *testing.T) {
	for _, test := range varIntTests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			var buf bytes.Buffer
			w := bufio.NewWriter(&buf)
			if err := WriteVarInt(w, test.val); err != nil {
				t.Error(err)
			}

			if err := w.Flush(); err != nil {
				t.Error(err)
			}

			got := buf.Bytes()
			if !bytes.Equal(test.bytes, got) {
				t.Errorf("Expected %d, got %d", test.bytes, got)
			}
		})
	}
}

func Test_VarIntSize(t *testing.T) {
	for _, test := range varIntTests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			if VarIntSize(test.val) != len(test.bytes) {
				t.Errorf("Expected %d, got %d", len(test.bytes), VarIntSize(test.val))
			}
		})
	}
}
