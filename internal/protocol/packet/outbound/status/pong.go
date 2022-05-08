package status

import (
	"github.com/gitfyu/mable/internal/protocol"
)

type Pong struct {
	Time int64
}

func (Pong) PacketID() uint {
	return 0x01
}

func (p *Pong) MarshalPacket(w protocol.Writer) error {
	return protocol.WriteUint64(w, uint64(p.Time))
}
