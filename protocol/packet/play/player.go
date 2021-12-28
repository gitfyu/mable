package play

import (
	"github.com/gitfyu/mable/protocol"
	"github.com/gitfyu/mable/protocol/packet"
)

type InPlayer struct {
	OnGround bool
}

type InPlayerPos struct {
	X, Y, Z  float64
	OnGround bool
}

type InPlayerLook struct {
	Yaw, Pitch float32
	OnGround   bool
}

type InPlayerPosLook struct {
	X, Y, Z    float64
	Yaw, Pitch float32
	OnGround   bool
}

func init() {
	packet.RegisterInbound(protocol.StatePlay, 0x03, func() packet.Inbound {
		return &InPlayer{}
	})
	packet.RegisterInbound(protocol.StatePlay, 0x04, func() packet.Inbound {
		return &InPlayerPos{}
	})
	packet.RegisterInbound(protocol.StatePlay, 0x05, func() packet.Inbound {
		return &InPlayerLook{}
	})
	packet.RegisterInbound(protocol.StatePlay, 0x06, func() packet.Inbound {
		return &InPlayerPosLook{}
	})
}

func (p *InPlayer) UnmarshalPacket(r *protocol.ReadBuffer) {
	p.OnGround = r.ReadBool()
}

func (p *InPlayerPos) UnmarshalPacket(r *protocol.ReadBuffer) {
	p.X, p.Y, p.Z = r.ReadFloat64(), r.ReadFloat64(), r.ReadFloat64()
	p.OnGround = r.ReadBool()
}

func (p *InPlayerLook) UnmarshalPacket(r *protocol.ReadBuffer) {
	p.Yaw, p.Pitch = r.ReadFloat32(), r.ReadFloat32()
	p.OnGround = r.ReadBool()
}

func (p *InPlayerPosLook) UnmarshalPacket(r *protocol.ReadBuffer) {
	p.X, p.Y, p.Z = r.ReadFloat64(), r.ReadFloat64(), r.ReadFloat64()
	p.Yaw, p.Pitch = r.ReadFloat32(), r.ReadFloat32()
	p.OnGround = r.ReadBool()
}
