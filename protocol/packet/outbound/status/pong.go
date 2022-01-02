package status

import (
	"github.com/gitfyu/mable/protocol"
)

type Pong struct {
	Time int64
}

func (_ Pong) PacketID() uint {
	return 0x01
}

func (p *Pong) MarshalPacket(w *protocol.WriteBuffer) {
	w.WriteUint64(uint64(p.Time))
}
