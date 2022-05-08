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

func (BulkChunkData) PacketID() uint {
	return 0x26
}

func (c *BulkChunkData) MarshalPacket(w protocol.Writer) error {
	if err := protocol.WriteBool(w, c.SkyLightIncluded); err != nil {
		return err
	}
	if err := protocol.WriteVarInt(w, c.ChunkCount); err != nil {
		return err
	}

	for i := range c.Meta {
		if err := protocol.WriteUint32(w, uint32(c.Meta[i].X)); err != nil {
			return err
		}
		if err := protocol.WriteUint32(w, uint32(c.Meta[i].Z)); err != nil {
			return err
		}
		if err := protocol.WriteUint16(w, c.Meta[i].SectionMask); err != nil {
			return err
		}
	}

	_, err := w.Write(c.Data)
	return err
}
