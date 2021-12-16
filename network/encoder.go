package network

import (
	"bufio"
	"bytes"
	"github.com/gitfyu/mable/network/protocol"
)

// PacketEncoder is a utility for writing packets. The write functions all write to an internal buffer, which can be
// dispatched using WritePacket, which will automatically prepend the size of the packet. The write functions return
// a boolean indicating success to allow easy chaining of calls, for example:
//
// ok := e.WriteVarInt(v1) && e.WriteVarInt(v2)
//
// In case !ok, LastError can be used to retrieve the error. Use NewPacketEncoder to create a new instance of
// PacketEncoder.
type PacketEncoder struct {
	out        *bufio.Writer
	packetData bytes.Buffer
	err        error
}

// NewPacketEncoder creates a new PacketEncoder
func NewPacketEncoder(w *bufio.Writer) *PacketEncoder {
	return &PacketEncoder{
		out: w,
	}
}

// LastError returns the error that occurred during a previous call to a write function. If the previous operation was
// successful, the return value is undefined.
func (e *PacketEncoder) LastError() error {
	return e.err
}

// WriteVarInt writes a protocol.VarInt to the buffer
func (e *PacketEncoder) WriteVarInt(v protocol.VarInt) bool {
	if err := protocol.WriteVarInt(&e.packetData, v); err != nil {
		e.err = err
		return false
	}
	return true
}

// WriteVarLong writes a protocol.VarLong to the buffer
func (e *PacketEncoder) WriteVarLong(v protocol.VarLong) bool {
	if err := protocol.WriteVarLong(&e.packetData, v); err != nil {
		e.err = err
		return false
	}
	return true
}

// WriteString writes a string to the buffer and prepends the length
func (e *PacketEncoder) WriteString(s string) bool {
	if !e.WriteVarInt(protocol.VarInt(len(s))) {
		return false
	}

	if _, err := e.packetData.Write([]byte(s)); err != nil {
		e.err = err
		return false
	}

	return true
}

// WritePacket sends the current size and contents of the buffer to the output writer, after which the buffer will be
// reset. The flush parameter determines if the output writer should be flushed.
func (e *PacketEncoder) WritePacket(flush bool) bool {
	if err := protocol.WriteVarLong(e.out, protocol.VarLong(e.packetData.Len())); err != nil {
		e.err = err
		return false
	}

	if _, err := e.out.Write(e.packetData.Bytes()); err != nil {
		e.err = err
		return false
	}

	e.packetData.Reset()
	if flush {
		if err := e.out.Flush(); err != nil {
			e.err = err
			return false
		}
	}

	return true
}
