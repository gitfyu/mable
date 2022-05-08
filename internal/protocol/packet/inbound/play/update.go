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

func (p *Update) UnmarshalPacket(r protocol.Reader) error {
	var err error
	if p.HasPos {
		if p.X, err = protocol.ReadFloat64(r); err != nil {
			return err
		}
		if p.Y, err = protocol.ReadFloat64(r); err != nil {
			return err
		}
		if p.Z, err = protocol.ReadFloat64(r); err != nil {
			return err
		}
	}
	if p.HasLook {
		if p.Yaw, err = protocol.ReadFloat32(r); err != nil {
			return err
		}
		if p.Pitch, err = protocol.ReadFloat32(r); err != nil {
			return err
		}
	}
	p.OnGround, err = protocol.ReadBool(r)
	return err
}
