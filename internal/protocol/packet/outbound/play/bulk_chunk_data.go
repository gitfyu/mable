package play

import (
	"github.com/gitfyu/mable/internal/protocol"
)

type BulkChunkDataMeta struct {
	X, Z        int32
	SectionMask uint16
}

type BulkChunkData struct {
	SkyLightIncluded bool
	ChunkCount       int32
	Meta             []BulkChunkDataMeta
	Data             []byte
}

func (_ BulkChunkData) PacketID() uint {
	return 0x26
}

func (c *BulkChunkData) MarshalPacket(w *protocol.WriteBuffer) {
	w.WriteBool(c.SkyLightIncluded)
	w.WriteVarInt(c.ChunkCount)

	for i := range c.Meta {
		w.WriteUint32(uint32(c.Meta[i].X))
		w.WriteUint32(uint32(c.Meta[i].Z))
		w.WriteUint16(c.Meta[i].SectionMask)
	}

	w.WriteBytes(c.Data)
}
