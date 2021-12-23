package mable

import (
	"context"
	"github.com/gitfyu/mable/protocol"
	"github.com/gitfyu/mable/protocol/packet"
	"github.com/gitfyu/mable/world"
)

func handlePlay(c *conn) error {
	defer c.player.Close()

	if err := writeJoinGame(c); err != nil {
		return err
	}
	if err := c.player.SendChunkData(0, 0); err != nil {
		return err
	}
	if err := c.player.SetSpawnPos(0, 0, 0); err != nil {
		return err
	}
	if err := c.player.Teleport(world.NewPos(8, 16, 8, 0, 0)); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go c.player.Update(ctx)

	for c.IsOpen() {
		_, _, err := c.readPacket()
		if err != nil {
			return err
		}
	}

	return nil
}

func writeJoinGame(c *conn) error {
	buf := packet.AcquireBuffer()
	defer packet.ReleaseBuffer(buf)

	buf.WriteInt(int32(c.player.GetEntityID()))
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
