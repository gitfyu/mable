package handshake

import (
	"github.com/gitfyu/mable/protocol"
)

type Handshake struct {
	ProtoVer  uint
	Addr      string
	Port      uint16
	NextState protocol.State
}

func (h *Handshake) UnmarshalPacket(r *protocol.ReadBuffer) {
	h.ProtoVer = uint(r.ReadVarInt())
	h.Addr = r.ReadString()
	h.Port = r.ReadUint16()
	h.NextState = protocol.State(r.ReadVarInt())
}
