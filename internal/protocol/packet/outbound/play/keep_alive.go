package play

import (
	"github.com/gitfyu/mable/internal/protocol"
)

type KeepAlive struct {
	ID int
}

func (_ KeepAlive) PacketID() uint {
	return 0x00
}

func (k *KeepAlive) MarshalPacket(w protocol.Writer) error {
	return protocol.WriteVarInt(w, int32(k.ID))
}
