package status

import (
	"github.com/gitfyu/mable/protocol"
	"github.com/gitfyu/mable/protocol/packet"
)

type Request struct{}

func init() {
	r := Request{}
	packet.RegisterInbound(protocol.StateStatus, 0x00, func() packet.Inbound {
		return r
	})
}

func (_ Request) UnmarshalPacket(_ *protocol.ReadBuffer) {
}
