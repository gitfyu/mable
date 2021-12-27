package play

import "github.com/gitfyu/mable/protocol"

type OutKeepAlive struct {
	ID int
}

type InKeepAlive struct {
	ID int
}

func (k *OutKeepAlive) MarshalPacket(w *protocol.WriteBuffer) {
	w.WriteVarInt(0x00)
	w.WriteVarInt(k.ID)
}

func (k *InKeepAlive) UnmarshalPacket(r *protocol.ReadBuffer) {
	k.ID = r.ReadVarInt()
}
