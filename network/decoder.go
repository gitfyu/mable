package network

import (
	"bufio"
	"github.com/gitfyu/mable/network/protocol"
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

// ReadVarInt reads a single protocol.VarInt
func (d *PacketDecoder) ReadVarInt(v *protocol.VarInt) bool {
	if err := protocol.ReadVarInt(d.reader, v); err != nil {
		d.err = err
		return false
	}

	return true
}

// ReadVarLong reads a single protocol.VarLong
func (d *PacketDecoder) ReadVarLong(v *protocol.VarLong) bool {
	if err := protocol.ReadVarLong(d.reader, v); err != nil {
		d.err = err
		return false
	}

	return true
}

// ReadString reads a single string. LastError will report ErrStringNegLen or ErrStringTooBig for illegal inputs.
func (d *PacketDecoder) ReadString(s *string) bool {
	var size protocol.VarInt
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
