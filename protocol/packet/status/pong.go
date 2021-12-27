package status

import (
	"github.com/gitfyu/mable/protocol"
)

type Pong struct {
	Time int64
}

func (p *Pong) MarshalPacket(w *protocol.WriteBuffer) {
	w.WriteVarInt(0x01)
	w.WriteUint64(uint64(p.Time))
}
