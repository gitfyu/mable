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

func (JoinGame) PacketID() uint {
	return 0x01
}

func (c *JoinGame) MarshalPacket(w protocol.Writer) error {
	if err := protocol.WriteUint32(w, uint32(c.EntityID)); err != nil {
		return err
	}
	if err := w.WriteByte(c.Gamemode); err != nil {
		return err
	}
	if err := w.WriteByte(uint8(c.Dimension)); err != nil {
		return err
	}
	if err := w.WriteByte(c.Difficulty); err != nil {
		return err
	}
	if err := w.WriteByte(c.MaxPlayers); err != nil {
		return err
	}
	if err := protocol.WriteString(w, c.LevelType); err != nil {
		return err
	}
	return protocol.WriteBool(w, c.ReduceDbgInfo)
}
