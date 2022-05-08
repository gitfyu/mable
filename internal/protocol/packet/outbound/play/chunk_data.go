package play

import (
	"github.com/gitfyu/mable/internal/protocol"
)

type ChunkData struct {
	X, Z      int32
	FullChunk bool
	Mask      uint16
	Data      []byte
}

func (ChunkData) PacketID() uint {
	return 0x21
}

func (c *ChunkData) MarshalPacket(w protocol.Writer) error {
	if err := protocol.WriteUint32(w, uint32(c.X)); err != nil {
		return err
	}
	if err := protocol.WriteUint32(w, uint32(c.Z)); err != nil {
		return err
	}
	if err := protocol.WriteBool(w, c.FullChunk); err != nil {
		return err
	}
	if err := protocol.WriteUint16(w, c.Mask); err != nil {
		return err
	}
	return protocol.WriteByteArray(w, c.Data)
}
