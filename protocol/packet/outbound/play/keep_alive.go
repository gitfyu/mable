package play

import (
	"github.com/gitfyu/mable/protocol"
)

type KeepAlive struct {
	ID int
}

func (_ KeepAlive) PacketID() uint {
	return 0x00
}

func (k *KeepAlive) MarshalPacket(w *protocol.WriteBuffer) {
	w.WriteVarInt(int32(k.ID))
}
