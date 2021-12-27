package status

import (
	"github.com/gitfyu/mable/protocol"
	"github.com/gitfyu/mable/protocol/packet"
)

type Ping struct {
	Time int64
}

func init() {
	packet.RegisterInbound(protocol.StateStatus, 0x01, func() packet.Inbound {
		return &Ping{}
	})
}

func (p *Ping) UnmarshalPacket(r *protocol.ReadBuffer) {
	p.Time = int64(r.ReadUint64())
}
