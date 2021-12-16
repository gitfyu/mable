package network

import (
	"github.com/gitfyu/mable/network/protocol"
	"github.com/pkg/errors"
	"io"
)

// PacketData is a utility for decoding packets. To use it, load the raw data using Load and then use the getter
// functions to read the values. The getter functions will panic if you try to read more data than available in the
// buffer.
type PacketData struct {
	data []byte
}

// Load initializes the PacketData with new data
func (r *PacketData) Load(data []byte) {
	r.data = data
}

func (r *PacketData) GetVarInt() protocol.VarInt {
	var v protocol.VarInt
	if err := protocol.ReadVarInt(r, &v); err != nil {
		panic(errors.WithStack(io.EOF))
	}

	return v
}

func (r *PacketData) GetString() string {
	n := int(r.GetVarInt())
	if n < 0 || n > len(r.data) {
		panic(errors.WithStack(io.EOF))
	} else if n == 0 {
		return ""
	} else {
		s := string(r.data[:n])
		r.data = r.data[n:]
		return s
	}
}

func (r *PacketData) GetUnsignedShort() uint16 {
	if len(r.data) < 2 {
		panic(errors.WithStack(io.EOF))
	}

	v := uint16(r.data[0])<<8 | uint16(r.data[1])
	r.data = r.data[2:]
	return v
}

// ReadByte should not be used directly
func (r *PacketData) ReadByte() (byte, error) {
	if len(r.data) == 0 {
		return 0, io.EOF
	}
	b := r.data[0]
	r.data = r.data[1:]
	return b, nil
}
