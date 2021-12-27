package play

import "github.com/gitfyu/mable/protocol"

type OutChunkData struct {
	X, Z      int32
	FullChunk bool
	Mask      uint16
	Data      []byte
}

func (c *OutChunkData) MarshalPacket(w *protocol.WriteBuffer) {
	w.WriteVarInt(0x21)
	w.WriteUint32(uint32(c.X))
	w.WriteUint32(uint32(c.Z))
	w.WriteBool(c.FullChunk)
	w.WriteUint16(c.Mask)
	w.WriteVarInt(len(c.Data))
	w.WriteBytes(c.Data)
}
