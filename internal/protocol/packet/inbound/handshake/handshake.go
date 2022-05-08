package handshake

import (
	"github.com/gitfyu/mable/internal/protocol"
	"github.com/gitfyu/mable/internal/protocol/packet"
)

type Handshake struct {
	ProtoVer  int32
	Addr      string
	Port      uint16
	NextState int32
}

func init() {
	packet.RegisterInbound(protocol.StateHandshake, 0x00, func() packet.Inbound {
		return &Handshake{}
	})
}

func (h *Handshake) UnmarshalPacket(r protocol.Reader) error {
	var err error

	if h.ProtoVer, err = protocol.ReadVarInt(r); err != nil {
		return err
	}
	if h.Addr, err = protocol.ReadString(r); err != nil {
		return err
	}
	if h.Port, err = protocol.ReadUint16(r); err != nil {
		return err
	}
	h.NextState, err = protocol.ReadVarInt(r)
	return err
}
