package mable

import (
	"github.com/gitfyu/mable/protocol"
	"github.com/gitfyu/mable/protocol/packet"
	"github.com/gitfyu/mable/world"
	"time"
)

func handlePlay(c *conn) error {
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

	go func() {
		ticker := time.NewTicker(time.Second * 5)
		for {
			<-ticker.C
			_ = c.player.Ping()
		}
	}()

	// TODO
	select {}
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

	return c.WritePacket(packet.PlayJoinGame, buf)
}
