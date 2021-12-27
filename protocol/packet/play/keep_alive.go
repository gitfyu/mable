package play

import (
	"github.com/gitfyu/mable/protocol"
	"github.com/gitfyu/mable/protocol/packet"
)

type InKeepAlive struct {
	ID int
}

func init() {
	packet.RegisterInbound(protocol.StatePlay, 0x00, func() packet.Inbound {
		return &InKeepAlive{}
	})
}

func (k *InKeepAlive) UnmarshalPacket(r *protocol.ReadBuffer) {
	k.ID = r.ReadVarInt()
}

type OutKeepAlive struct {
	ID int
}

func (_ OutKeepAlive) PacketID() uint {
	return 0x00
}

func (k *OutKeepAlive) MarshalPacket(w *protocol.WriteBuffer) {
	w.WriteVarInt(int32(k.ID))
}
