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

func (p *Position) MarshalPacket(w protocol.Writer) error {
	if err := protocol.WriteFloat64(w, p.X); err != nil {
		return err
	}
	if err := protocol.WriteFloat64(w, p.Y); err != nil {
		return err
	}
	if err := protocol.WriteFloat64(w, p.Z); err != nil {
		return err
	}
	if err := protocol.WriteFloat32(w, p.Yaw); err != nil {
		return err
	}
	if err := protocol.WriteFloat32(w, p.Pitch); err != nil {
		return err
	}
	return w.WriteByte(p.Flags)
}
