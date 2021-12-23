package mable

import (
	"context"
	"github.com/gitfyu/mable/entity"
	"github.com/gitfyu/mable/protocol"
	"github.com/gitfyu/mable/protocol/packet"
	"github.com/gitfyu/mable/world"
	"github.com/google/uuid"
)

func handlePlay(c *conn, username string, id uuid.UUID) error {
	p := entity.NewPlayer(username, id, c, world.Default)

	if err := writeJoinGame(c, p.GetEntityID()); err != nil {
		return err
	}
	if err := p.SendChunkData(0, 0); err != nil {
		return err
	}
	if err := p.SetSpawnPos(0, 0, 0); err != nil {
		return err
	}
	if err := p.Teleport(world.NewPos(8, 16, 8, 0, 0)); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go p.Update(ctx)

	for c.IsOpen() {
		_, _, err := c.readPacket()
		if err != nil {
			return err
		}
	}

	return nil
}

func writeJoinGame(c *conn, id entity.ID) error {
	buf := packet.AcquireBuffer()
	defer packet.ReleaseBuffer(buf)

	buf.WriteInt(int32(id))
	// creative gamemode
	buf.WriteUnsignedByte(uint8(1))
	// overworld dimension
	buf.WriteSignedByte(0)
	// easy difficulty
	buf.WriteUnsignedByte(1)
	// max players, unused
	buf.WriteUnsignedByte(0)
	// level type
	buf.WriteString("flat")

	if c.version == protocol.Version_1_8 {
		// disable reduced debug info
		buf.WriteBool(false)
	}

	return c.WritePacket(packet.PlayServerJoinGame, buf)
}
