package entity

import (
	"github.com/gitfyu/mable/chat"
	"github.com/gitfyu/mable/protocol"
	"github.com/gitfyu/mable/protocol/packet"
	"github.com/gitfyu/mable/world"
	"github.com/google/uuid"
	"sync"
)

const PlayerEyeHeight = 1.62

type PlayerConn interface {
	Version() protocol.Version
	WritePacket(id packet.ID, buf *packet.Buffer) error
	Disconnect(reason *chat.Msg) error
}

type Player struct {
	Entity
	name    string
	uid     uuid.UUID
	conn    PlayerConn
	pos     world.Position
	posLock sync.RWMutex
}

func NewPlayer(name string, uid uuid.UUID, conn PlayerConn) *Player {
	return &Player{
		Entity: NewEntity(),
		name:   name,
		uid:    uid,
		conn:   conn,
	}
}

func (p *Player) SetSpawnPos(x, y, z int32) error {
	if p.conn == nil {
		return nil
	}

	buf := packet.AcquireBuffer()
	defer packet.ReleaseBuffer(buf)

	if p.conn.Version() == protocol.Version_1_8 {
		buf.WritePosition(x, y, z)
	} else {
		buf.WriteInt(x)
		buf.WriteInt(y)
		buf.WriteInt(z)
	}

	return p.conn.WritePacket(packet.PlaySpawnPosition, buf)
}

func (p *Player) Teleport(pos world.Position) error {
	p.posLock.Lock()
	defer p.posLock.Unlock()

	p.pos = pos

	if p.conn == nil {
		return nil
	}

	buf := packet.AcquireBuffer()
	defer packet.ReleaseBuffer(buf)

	buf.WriteDouble(pos.X)
	buf.WriteDouble(pos.Y + PlayerEyeHeight)
	buf.WriteDouble(pos.Z)
	buf.WriteFloat(pos.Yaw)
	buf.WriteFloat(pos.Pitch)

	if p.conn.Version() == protocol.Version_1_8 {
		// flags indicating all values are absolute
		buf.WriteSignedByte(0)
	} else {
		// on ground, useless
		buf.WriteBool(false)
	}

	return p.conn.WritePacket(packet.PlayPosAndLook, buf)
}
