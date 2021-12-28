package mable

import (
	"github.com/gitfyu/mable/entity/player"
	"github.com/gitfyu/mable/protocol/packet/play"
	"github.com/gitfyu/mable/world"
	"github.com/google/uuid"
)

// handlePlay creates the player and handles all packets until the connection is closed
func handlePlay(c *conn, username string, id uuid.UUID) error {
	p := player.NewPlayer(username, id, c, world.Default)
	defer p.Destroy()

	c.WritePacket(&play.OutJoinGame{
		EntityID:      int(p.GetEntityID()),
		Gamemode:      1,
		Dimension:     0,
		Difficulty:    1,
		MaxPlayers:    0,
		LevelType:     "flat",
		ReduceDbgInfo: false,
	})
	p.Teleport(world.Pos{
		X: 8,
		Y: 16,
		Z: 8,
	})

	for c.IsOpen() {
		pk, err := c.readPacket()
		if err != nil {
			return err
		}

		p.HandlePacket(pk)
	}

	return nil
}
