package handshake

import (
	"github.com/gitfyu/mable/protocol"
	"github.com/gitfyu/mable/protocol/packet"
)

type Handshake struct {
	ProtoVer  uint
	Addr      string
	Port      uint16
	NextState protocol.State
}

func init() {
	packet.RegisterInbound(protocol.StateHandshake, 0x00, func() packet.Inbound {
		return &Handshake{}
	})
}

func (h *Handshake) UnmarshalPacket(r *protocol.ReadBuffer) {
	h.ProtoVer = uint(r.ReadVarInt())
	h.Addr = r.ReadString()
	h.Port = r.ReadUint16()
	h.NextState = protocol.State(r.ReadVarInt())
}
