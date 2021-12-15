package network

import (
	"bufio"
	"io"
)

// PacketDecoder is a wrapper around bufio.Reader used to decode common data types in the Minecraft protocol. Its
// functions all return a bool instead of an error to allow easy chaining of several calls, for example:
//
// ok := d.ReadVarIntAndSize(&v1) && d.ReadVarIntAndSize(&v2) && d.ReadVarIntAndSize(&v3)
//
// If !ok, you can obtain the error using LastError.
type PacketDecoder struct {
	reader *bufio.Reader
	err    error
}

func NewPacketDecoder(r *bufio.Reader) *PacketDecoder {
	return &PacketDecoder{
		reader: r,
	}
}

// LastError returns the error that occurred during a previous call to a PacketDecoder function. If the previous
// operation was successful, the return value is undefined.
func (d *PacketDecoder) LastError() error {
	return d.err
}

func (d *PacketDecoder) Skip(n int) bool {
	if _, err := d.reader.Discard(n); err != nil {
		d.err = err
		return false
	}

	return true
}

// ReadVarInt is the same as ReadVarIntAndSize, except it does not return the size.
func (d *PacketDecoder) ReadVarInt(v *VarInt) bool {
	return d.ReadVarIntAndSize(v, nil)
}

// ReadVarIntAndSize reads a single VarInt. n will be set to the number of bytes read, unless it is set to nil.
// If the result is too big, LastError will be ErrVarIntTooBig.
func (d *PacketDecoder) ReadVarIntAndSize(v *VarInt, n *int) bool {
	var tmp VarLong
	ok := d.readVarLong(&tmp, varIntMaxBytes, n)

	*v = VarInt(tmp)
	return ok
}

// ReadVarLong reads a single VarLong. If the result is too big, LastError will be ErrVarIntTooBig.
func (d *PacketDecoder) ReadVarLong(v *VarLong) bool {
	return d.readVarLong(v, varLongMaxBytes, nil)
}

func (d *PacketDecoder) readVarLong(v *VarLong, maxSize int, size *int) bool {
	var n int
	var b byte

	for {
		b, d.err = d.reader.ReadByte()
		if d.err != nil {
			return false
		}

		*v |= VarLong(b&0x7F) << (n * 7)
		n++
		if n > maxSize {
			d.err = ErrVarIntTooBig
			return false
		}

		if (b & 0x80) == 0 {
			break
		}
	}

	if size != nil {
		*size = n
	}

	return true
}

// ReadString reads a single string. LastError will report ErrStringNegLen or ErrStringTooBig for illegal inputs.
func (d *PacketDecoder) ReadString(s *string) bool {
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
	if _, err := io.ReadFull(d.reader, buf); err != nil {
		d.err = err
		return false
	}

	*s = string(buf)
	return true
}

func (d *PacketDecoder) ReadByte(v *byte) bool {
	b, err := d.reader.ReadByte()
	if err != nil {
		d.err = err
		return false
	}

	*v = b
	return true
}

func (d *PacketDecoder) ReadUnsignedShort(v *uint16) bool {
	var b1, b2 byte
	if !d.ReadByte(&b1) || !d.ReadByte(&b2) {
		return false
	}

	*v = uint16(b1)<<8 | uint16(b2)
	return true
}
