package play

import (
	"github.com/gitfyu/mable/internal/protocol"
)

type Position struct {
	X, Y, Z    float64
	Yaw, Pitch float32
	Flags      uint8
}

func (_ Position) PacketID() uint {
	return 0x08
}

func (p *Position) MarshalPacket(w *protocol.WriteBuffer) {
	w.WriteFloat64(p.X)
	w.WriteFloat64(p.Y)
	w.WriteFloat64(p.Z)
	w.WriteFloat32(p.Yaw)
	w.WriteFloat32(p.Pitch)
	w.WriteUint8(p.Flags)
}
