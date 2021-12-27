package play

import (
	"github.com/gitfyu/mable/protocol"
)

type OutPosition struct {
	X, Y, Z    float64
	Yaw, Pitch float32
	Flags      uint8
}

func (p *OutPosition) MarshalPacket(w *protocol.WriteBuffer) {
	w.WriteVarInt(0x08)
	w.WriteFloat64(p.X)
	w.WriteFloat64(p.Y)
	w.WriteFloat64(p.Z)
	w.WriteFloat32(p.Yaw)
	w.WriteFloat32(p.Pitch)
	w.WriteUint8(p.Flags)
}
