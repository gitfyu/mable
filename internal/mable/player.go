package mable

import (
	"github.com/gitfyu/mable/entity"
	"github.com/gitfyu/mable/protocol"
	"github.com/gitfyu/mable/protocol/packet"
	"github.com/gitfyu/mable/world"
	"github.com/google/uuid"
)

type player struct {
	conn *conn
	name string
	uid  uuid.UUID
	id   entity.ID
}

func (p *player) GetEntityID() entity.ID {
	return p.id
}

func (p *player) Teleport(pos world.Position) error {
	buf := packet.AcquireBuffer()
	defer packet.ReleaseBuffer(buf)

	buf.WriteDouble(pos.X)
	buf.WriteDouble(pos.Y + entity.PlayerEyeHeight)
	buf.WriteDouble(pos.Z)
	buf.WriteFloat(pos.Yaw)
	buf.WriteFloat(pos.Pitch)

	if p.conn.version == protocol.Version_1_8 {
		// flags indicating all values are absolute
		buf.WriteSignedByte(0)
	} else {
		// on ground, useless
		buf.WriteBool(false)
	}

	return p.conn.WritePacket(packet.PlayPosAndLook, buf)
}
