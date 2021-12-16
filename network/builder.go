package network

import (
	"bytes"
	"github.com/gitfyu/mable/network/protocol"
	"github.com/gitfyu/mable/network/protocol/packet"
	"sync"
)

// PacketBuilder is a utility for building packets. It is recommended to acquire an instance using AcquirePacketBuilder,
// which you must release when you are done with it using ReleasePacketBuilder.
type PacketBuilder struct {
	buf bytes.Buffer
}

var packetBuilderPool = sync.Pool{
	New: func() interface{} {
		return &PacketBuilder{}
	},
}

// AcquirePacketBuilder returns a PacketBuilder, which must be released afterwards using ReleasePacketBuilder
func AcquirePacketBuilder() *PacketBuilder {
	return packetBuilderPool.Get().(*PacketBuilder)
}

// ReleasePacketBuilder releases a PacketBuilder, after which you must no longer use it
func ReleasePacketBuilder(b *PacketBuilder) {
	packetBuilderPool.Put(b)
}

// Init is the first function that you must call when writing a packet using the corresponding packet.ID
func (p *PacketBuilder) Init(id packet.ID) *PacketBuilder {
	// reserve space at the start for packet size
	if p.buf.Len() >= protocol.VarIntMaxBytes {
		p.buf.Truncate(protocol.VarIntMaxBytes)
	} else {
		p.buf.Write(make([]byte, protocol.VarIntMaxBytes))
	}

	return p.PutVarInt(protocol.VarInt(id))
}

func (p *PacketBuilder) PutVarInt(v protocol.VarInt) *PacketBuilder {
	b := make([]byte, protocol.VarIntSize(v))
	protocol.WriteVarInt(b, v)
	p.buf.Write(b)
	return p
}

func (p *PacketBuilder) PutBytes(b []byte) *PacketBuilder {
	p.buf.Write(b)
	return p
}

func (p *PacketBuilder) PutString(s string) *PacketBuilder {
	b := []byte(s)
	p.PutVarInt(protocol.VarInt(len(b)))
	return p.PutBytes(b)
}

// ToBytes converts all the packet data that has been written into a single byte slice. After calling this function,
// you should not call any functions from this PacketBuilder again.
func (p *PacketBuilder) ToBytes() []byte {
	b := p.buf.Bytes()
	size := protocol.VarInt(len(b) - protocol.VarIntMaxBytes)
	sizeBytes := protocol.VarIntSize(size)

	b = b[protocol.VarIntMaxBytes-sizeBytes:]
	protocol.WriteVarInt(b, size)
	return b
}
