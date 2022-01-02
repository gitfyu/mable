package play

import (
	"github.com/gitfyu/mable/internal/protocol"
)

type JoinGame struct {
	EntityID      int
	Gamemode      uint8
	Dimension     int8
	Difficulty    uint8
	MaxPlayers    uint8
	LevelType     string
	ReduceDbgInfo bool
}

func (_ JoinGame) PacketID() uint {
	return 0x01
}

func (c *JoinGame) MarshalPacket(w *protocol.WriteBuffer) {
	w.WriteUint32(uint32(c.EntityID))
	w.WriteUint8(c.Gamemode)
	w.WriteUint8(uint8(c.Dimension))
	w.WriteUint8(c.Difficulty)
	w.WriteUint8(c.MaxPlayers)
	w.WriteString(c.LevelType)
	w.WriteBool(c.ReduceDbgInfo)
}
