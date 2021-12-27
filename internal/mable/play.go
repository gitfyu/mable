package mable

import (
	"github.com/gitfyu/mable/entity"
	"github.com/gitfyu/mable/protocol/packet/play"
	"github.com/gitfyu/mable/world"
	"github.com/google/uuid"
)

// handlePlay creates the player and handles all packets until the connection is closed
func handlePlay(c *conn, username string, id uuid.UUID) error {
	p := entity.NewPlayer(username, id, c)
	defer p.Destroy()

	if err := writeJoinGame(c, p.GetEntityID()); err != nil {
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

	for c.IsOpen() {
		pk, err := c.readPacket()
		if err != nil {
			return err
		}

		p.HandlePacket(pk)
	}

	return nil
}

func writeJoinGame(c *conn, id entity.ID) error {
	pk := play.OutJoinGame{
		EntityID:      int(id),
		Gamemode:      1,
		Dimension:     0,
		Difficulty:    1,
		MaxPlayers:    0,
		LevelType:     "flat",
		ReduceDbgInfo: false,
	}
	return c.WritePacket(&pk)
}
