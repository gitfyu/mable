package mable

import (
	"github.com/gitfyu/mable/entity"
	"github.com/gitfyu/mable/protocol"
	"github.com/gitfyu/mable/protocol/packet"
)

func handlePlay(c *conn) error {
	if err := writeJoinGame(c); err != nil {
		return err
	}
	if err := writeSpawnPosition(c); err != nil {
		return err
	}
	if err := writePlayerPosAndLook(c); err != nil {
		return err
	}

	// TODO
	select {}
}

func writeJoinGame(c *conn) error {
	buf := packet.AcquireBuffer()
	defer packet.ReleaseBuffer(buf)

	buf.WriteInt(int32(c.player.id))
	// survival gamemode
	buf.WriteUnsignedByte(uint8(0))
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

func writeSpawnPosition(c *conn) error {
	buf := packet.AcquireBuffer()
	defer packet.ReleaseBuffer(buf)

	var x, y, z int32

	if c.version == protocol.Version_1_8 {
		buf.WritePosition(x, y, z)
	} else {
		buf.WriteInt(x)
		buf.WriteInt(y)
		buf.WriteInt(z)
	}

	return c.WritePacket(packet.PlaySpawnPosition, buf)
}

func writePlayerPosAndLook(c *conn) error {
	buf := packet.AcquireBuffer()
	defer packet.ReleaseBuffer(buf)

	var x, y, z float64
	var yaw, pitch float32

	buf.WriteDouble(x)
	buf.WriteDouble(y + entity.PlayerEyeHeight)
	buf.WriteDouble(z)
	buf.WriteFloat(yaw)
	buf.WriteFloat(pitch)

	if c.version == protocol.Version_1_8 {
		// bitfield for absolute/relative values
		buf.WriteSignedByte(0)
	} else {
		// on ground
		buf.WriteBool(false)
	}

	return c.WritePacket(packet.PlayPosAndLook, buf)
}
