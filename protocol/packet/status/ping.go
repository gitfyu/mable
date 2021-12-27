package status

import (
	"github.com/gitfyu/mable/protocol"
)

type Ping struct {
	Time int64
}

func (p *Ping) UnmarshalPacket(r *protocol.ReadBuffer) {
	p.Time = int64(r.ReadUint64())
}
