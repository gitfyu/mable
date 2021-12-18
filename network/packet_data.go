package network

import (
	"encoding/binary"
	"github.com/gitfyu/mable/network/protocol"
	"io"
)

// PacketData is a utility for decoding packets. To use it, load the raw data using Load and then use the getter
// functions to read the values. The getter functions will panic if you try to read more data than available in the
// buffer or if the data is invalid.
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
		panic(io.EOF)
	}

	return v
}

func (r *PacketData) GetString() string {
	n := int(r.GetVarInt())
	if n < 0 || n > len(r.data) {
		panic(io.EOF)
	} else if n == 0 {
		return ""
	} else {
		s := string(r.data[:n])
		r.data = r.data[n:]
		return s
	}
}

func (r *PacketData) GetBytes(n int) []byte {
	if len(r.data) < n {
		panic(io.EOF)
	}

	b := r.data[:n]
	r.data = r.data[n:]
	return b
}

func (r *PacketData) GetUnsignedShort() uint16 {
	return binary.BigEndian.Uint16(r.GetBytes(2))
}

func (r *PacketData) GetLong() int64 {
	return int64(binary.BigEndian.Uint64(r.GetBytes(8)))
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
