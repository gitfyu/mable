package mable

import (
	"github.com/gitfyu/mable/entity"
	"github.com/gitfyu/mable/protocol/packet"
	"github.com/gitfyu/mable/world"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// handlePlay creates the player and handles all packets until the connection is closed
func handlePlay(c *conn, username string, id uuid.UUID) error {
	p := entity.NewPlayer(username, id, c)
	defer p.Destroy()

	if err := writeJoinGame(c, p.GetEntityID()); err != nil {
		return err
	}
	if err := p.SendChunkData(0, 0); err != nil {
		return err
	}

	err := p.SetPos(world.Pos{
		World: world.Default,
		X:     8,
		Y:     16,
		Z:     8,
	})
	if err != nil {
		return err
	}

	handlers := map[packet.ID]func(buffer *packet.Buffer) error{
		packet.PlayClientKeepAlive: func(data *packet.Buffer) error {
			i, err := data.ReadVarInt()
			if err != nil {
				return err
			}

			log.Debug().Int("id", int(i)).Msg("KeepAlive")
			return nil
		},
	}

	for c.IsOpen() {
		id, data, err := c.readPacket()
		if err != nil {
			return err
		}
		
		if h, ok := handlers[id]; ok {
			if err := h(data); err != nil {
				return err
			}
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
	// disable reduced debug info
	buf.WriteBool(false)

	return c.WritePacket(packet.PlayServerJoinGame, buf)
}
