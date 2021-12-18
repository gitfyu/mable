package packet

import (
	"bytes"
	"encoding/binary"
	"github.com/gitfyu/mable/protocol"
	"sync"
)

// Builder is a utility for building packets. It is recommended to acquire an instance using AcquireBuilder,
// which you must release when you are done with it using ReleaseBuilder.
type Builder struct {
	buf bytes.Buffer
}

var builderPool = sync.Pool{
	New: func() interface{} {
		return &Builder{}
	},
}

// AcquireBuilder returns a Builder, which must be released afterwards using ReleaseBuilder
func AcquireBuilder() *Builder {
	return builderPool.Get().(*Builder)
}

// ReleaseBuilder releases a Builder, after which you must no longer use it
func ReleaseBuilder(b *Builder) {
	builderPool.Put(b)
}

// Init is the first function that you must call when writing a packet using the corresponding ID
func (p *Builder) Init(id ID) *Builder {
	// reserve space at the start for packet size
	if p.buf.Len() >= protocol.VarIntMaxBytes {
		p.buf.Truncate(protocol.VarIntMaxBytes)
	} else {
		p.buf.Write(make([]byte, protocol.VarIntMaxBytes))
	}

	return p.PutVarInt(protocol.VarInt(id))
}

func (p *Builder) PutVarInt(v protocol.VarInt) *Builder {
	b := make([]byte, protocol.VarIntSize(v))
	protocol.WriteVarInt(b, v)
	p.buf.Write(b)
	return p
}

func (p *Builder) PutBytes(b []byte) *Builder {
	p.buf.Write(b)
	return p
}

func (p *Builder) PutString(s string) *Builder {
	b := []byte(s)
	p.PutVarInt(protocol.VarInt(len(b)))
	return p.PutBytes(b)
}

func (p *Builder) PutStringFromBytes(b []byte) *Builder {
	p.PutVarInt(protocol.VarInt(len(b)))
	return p.PutBytes(b)
}

func (p *Builder) PutLong(v int64) *Builder {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return p.PutBytes(b)
}

// ToBytes converts all the packet data that has been written into a single byte slice. After calling this function,
// you should not call any functions from this Builder again.
func (p *Builder) ToBytes() []byte {
	b := p.buf.Bytes()
	size := protocol.VarInt(len(b) - protocol.VarIntMaxBytes)
	sizeBytes := protocol.VarIntSize(size)

	b = b[protocol.VarIntMaxBytes-sizeBytes:]
	protocol.WriteVarInt(b, size)
	return b
}
