package login

import (
	"github.com/gitfyu/mable/internal/protocol"
	"github.com/gitfyu/mable/internal/protocol/packet"
)

type Start struct {
	Username string
}

func init() {
	packet.RegisterInbound(protocol.StateLogin, 0x00, func() packet.Inbound {
		return &Start{}
	})
}

func (s *Start) UnmarshalPacket(r *protocol.ReadBuffer) {
	s.Username = r.ReadString()
}
