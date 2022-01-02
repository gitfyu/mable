package play

import (
	"github.com/gitfyu/mable/internal/protocol"
	"github.com/gitfyu/mable/internal/protocol/packet"
)

type KeepAlive struct {
	ID int
}

func init() {
	packet.RegisterInbound(protocol.StatePlay, 0x00, func() packet.Inbound {
		return &KeepAlive{}
	})
}

func (k *KeepAlive) UnmarshalPacket(r *protocol.ReadBuffer) {
	k.ID = r.ReadVarInt()
}
