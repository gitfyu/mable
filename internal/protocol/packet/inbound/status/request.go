package status

import (
	"github.com/gitfyu/mable/internal/protocol"
	"github.com/gitfyu/mable/internal/protocol/packet"
)

type Request struct{}

func init() {
	r := Request{}
	packet.RegisterInbound(protocol.StateStatus, 0x00, func() packet.Inbound {
		return r
	})
}

func (Request) UnmarshalPacket(protocol.Reader) error {
	return nil
}
