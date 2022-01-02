package play

import (
	"github.com/gitfyu/mable/internal/protocol"
	"github.com/gitfyu/mable/internal/protocol/packet"
)

type Update struct {
	HasPos, HasLook, OnGround bool
	X, Y, Z                   float64
	Yaw, Pitch                float32
}

func init() {
	packet.RegisterInbound(protocol.StatePlay, 0x03, func() packet.Inbound {
		return &Update{}
	})
	packet.RegisterInbound(protocol.StatePlay, 0x04, func() packet.Inbound {
		return &Update{HasPos: true}
	})
	packet.RegisterInbound(protocol.StatePlay, 0x05, func() packet.Inbound {
		return &Update{HasLook: true}
	})
	packet.RegisterInbound(protocol.StatePlay, 0x06, func() packet.Inbound {
		return &Update{HasPos: true, HasLook: true}
	})
}

func (p *Update) UnmarshalPacket(r *protocol.ReadBuffer) {
	if p.HasPos {
		p.X, p.Y, p.Z = r.ReadFloat64(), r.ReadFloat64(), r.ReadFloat64()
	}
	if p.HasLook {
		p.Yaw, p.Pitch = r.ReadFloat32(), r.ReadFloat32()
	}
	p.OnGround = r.ReadBool()
}
