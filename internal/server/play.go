package server

import (
	"github.com/gitfyu/mable/game"
	"github.com/gitfyu/mable/internal/protocol/packet/outbound/play"
	"github.com/google/uuid"
)

// handlePlay creates the player and handles all packets until the connection is closed.
func handlePlay(c *conn, username string, id uuid.UUID) error {
	g := c.serv.game
	p := game.NewPlayer(username, id, c)

	defer g.Schedule(func() {
		p.Close()
	})

	g.Schedule(func() {
		p.SetWorld(g.DefaultWorld())
		c.WritePacket(&play.JoinGame{
			EntityID:      int(p.EntityID()),
			Gamemode:      1,
			Dimension:     0,
			Difficulty:    1,
			MaxPlayers:    0,
			LevelType:     "flat",
			ReduceDbgInfo: false,
		})
		p.Teleport(game.Pos{
			X: 8,
			Y: 16,
			Z: 8,
		})
	})

	for c.IsOpen() {
		pk, err := c.readPacket()
		if err != nil {
			return err
		}

		g.Schedule(func() {
			p.HandlePacket(pk)
		})
	}

	return nil
}
