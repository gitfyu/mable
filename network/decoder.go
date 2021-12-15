package network

import (
	"bufio"
	"io"
)

// Decoder is a wrapper around bufio.Reader used to decode common data types in the Minecraft protocol. Its functions
// all return a bool instead of an error to allow easy chaining of several calls, for example:
//
// ok := d.ReadVarIntAndSize(&v1) && d.ReadVarIntAndSize(&v2) && d.ReadVarIntAndSize(v3)
//
// If !ok, you can obtain the error using LastError.
type Decoder struct {
	Reader *bufio.Reader
	err    error
}

// LastError returns the error that occurred during a previous call to a Decoder function. If the previous operation was
// successful, the return value is undefined.
func (d *Decoder) LastError() error {
	return d.err
}

// ReadVarInt is the same as ReadVarIntAndSize, except it does not return the size.
func (d *Decoder) ReadVarInt(v *VarInt) bool {
	return d.ReadVarIntAndSize(v, nil)
}

// ReadVarIntAndSize reads a single VarInt. n will be set to the number of bytes read, unless it is set to nil.
// If the result is too big, LastError will be ErrVarIntTooBig.
func (d *Decoder) ReadVarIntAndSize(v *VarInt, n *int) bool {
	var size int
	var b byte

	for {
		b, d.err = d.Reader.ReadByte()
		if d.err != nil {
			return false
		}

		*v |= VarInt(b&0x7F) << (size * 7)
		size++
		if size > varIntMaxBytes {
			d.err = ErrVarIntTooBig
			return false
		}

		if (b & 0x80) == 0 {
			break
		}
	}

	if n != nil {
		*n = size
	}

	return true
}

// ReadVarLong reads a single VarLong. If the result is too big, LastError will be ErrVarLongTooBig.
func (d *Decoder) ReadVarLong(v *VarLong) bool {
	var size int32
	var b byte

	for {
		b, d.err = d.Reader.ReadByte()
		if d.err != nil {
			return false
		}

		*v |= VarLong(b&0x7F) << (size * 7)
		size++
		if size > varLongMaxBytes {
			d.err = ErrVarLongTooBig
			return false
		}

		if (b & 0x80) == 0 {
			break
		}
	}

	return true
}

// ReadString reads a single string. LastError will report ErrStringNegLen or ErrStringTooBig for illegal inputs.
func (d *Decoder) ReadString(s *string) bool {
	var size VarInt
	if !d.ReadVarInt(&size) {
		return false
	}

	if size < 0 {
		d.err = ErrStringNegLen
		return false
	} else if size == 0 {
		*s = ""
		return true
	} else if size > stringMaxBytes {
		d.err = ErrStringTooBig
		return false
	}

	buf := make([]byte, size)
	if _, err := io.ReadFull(d.Reader, buf); err != nil {
		d.err = err
		return false
	}

	*s = string(buf)
	return true
}

func (d *Decoder) ReadByte(v *byte) bool {
	b, err := d.Reader.ReadByte()
	if err != nil {
		d.err = err
		return false
	}

	*v = b
	return true
}

func (d *Decoder) ReadUnsignedShort(v *uint16) bool {
	var b1, b2 byte
	if !d.ReadByte(&b1) || !d.ReadByte(&b2) {
		return false
	}

	*v = uint16(b1)<<8 | uint16(b2)
	return true
}
