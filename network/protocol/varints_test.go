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

type varLongTest struct {
	// the value
	val VarLong
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

// Obtained from https://wiki.vg/Protocol#VarInt_and_VarLong
var varLongTests = []varLongTest{
	{val: 0, bytes: []byte{0x00}},
	{val: 1, bytes: []byte{0x01}},
	{val: 2, bytes: []byte{0x02}},
	{val: 127, bytes: []byte{0x7f}},
	{val: 128, bytes: []byte{0x80, 0x01}},
	{val: 255, bytes: []byte{0xff, 0x01}},
	{val: 2147483647, bytes: []byte{0xff, 0xff, 0xff, 0xff, 0x07}},
	{val: 9223372036854775807, bytes: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f}},
	{val: -1, bytes: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01}},
	{val: -2147483648, bytes: []byte{0x80, 0x80, 0x80, 0x80, 0xf8, 0xff, 0xff, 0xff, 0xff, 0x01}},
	{val: -9223372036854775808, bytes: []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x01}},
}

var varLongInvalidTest = varIntTest{
	val: 0, bytes: []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x01},
}

func Test_ReadVarInt(t *testing.T) {
	for _, test := range varIntTests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			testRead(t, func(r ByteReader, v *VarLong) error {
				var tmp VarInt
				err := ReadVarInt(r, &tmp)
				*v = VarLong(tmp)
				return err
			}, test.bytes, VarLong(test.val))
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

func Test_ReadVarLong(t *testing.T) {
	for _, test := range varLongTests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			testRead(t, func(r ByteReader, v *VarLong) error {
				return ReadVarLong(r, v)
			}, test.bytes, test.val)
		})
	}
}

func Test_ReadVarLong_Invalid(t *testing.T) {
	r := bufio.NewReader(bytes.NewReader(varLongInvalidTest.bytes))
	var val VarLong
	if err := ReadVarLong(r, &val); err == nil {
		t.Error("Expected error")
	}
}

func Test_WriteVarInt(t *testing.T) {
	for _, test := range varIntTests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			testWrite(t, func(w ByteWriter) error {
				return WriteVarInt(w, test.val)
			}, test.bytes)
		})
	}
}

func Test_WriteVarLong(t *testing.T) {
	for _, test := range varLongTests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			testWrite(t, func(w ByteWriter) error {
				return WriteVarLong(w, test.val)
			}, test.bytes)
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

func Test_VarLongSize(t *testing.T) {
	for _, test := range varLongTests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			if VarLongSize(test.val) != len(test.bytes) {
				t.Errorf("Expected %d, got %d", len(test.bytes), VarLongSize(test.val))
			}
		})
	}
}

func testRead(t *testing.T, readFunc func(ByteReader, *VarLong) error, b []byte, expectVal VarLong) {
	r := bufio.NewReader(bytes.NewReader(b))
	var val VarLong

	if err := readFunc(r, &val); err != nil {
		t.Error(err)
	}
	if expectVal != val {
		t.Errorf("Expected val %d, got %d", expectVal, val)
	}
}

func testWrite(t *testing.T, writeFunc func(ByteWriter) error, expect []byte) {
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	if err := writeFunc(w); err != nil {
		t.Error(err)
	}

	if err := w.Flush(); err != nil {
		t.Error(err)
	}

	got := buf.Bytes()
	if !bytes.Equal(expect, got) {
		t.Errorf("Expected %d, got %d", expect, got)
	}
}
