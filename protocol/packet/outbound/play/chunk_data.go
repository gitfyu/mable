package play

import "github.com/gitfyu/mable/protocol"

type ChunkData struct {
	X, Z      int32
	FullChunk bool
	Mask      uint16
	Data      []byte
}

func (_ ChunkData) PacketID() uint {
	return 0x21
}

func (c *ChunkData) MarshalPacket(w *protocol.WriteBuffer) {
	w.WriteUint32(uint32(c.X))
	w.WriteUint32(uint32(c.Z))
	w.WriteBool(c.FullChunk)
	w.WriteUint16(c.Mask)
	w.WriteByteArrayWithLength(c.Data)
}
