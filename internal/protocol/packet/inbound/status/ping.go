package status

import (
	"github.com/gitfyu/mable/internal/protocol"
	"github.com/gitfyu/mable/internal/protocol/packet"
)

type Ping struct {
	Time int64
}

func init() {
	packet.RegisterInbound(protocol.StateStatus, 0x01, func() packet.Inbound {
		return &Ping{}
	})
}

func (p *Ping) UnmarshalPacket(r protocol.Reader) error {
	t, err := protocol.ReadUint64(r)
	p.Time = int64(t)
	return err
}
