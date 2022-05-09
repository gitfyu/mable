package game

import (
	"github.com/gitfyu/mable/chat"
	"github.com/gitfyu/mable/internal/protocol/packet"
	outbound "github.com/gitfyu/mable/internal/protocol/packet/outbound/play"
	"github.com/google/uuid"
)

const PlayerEyeHeight = 1.62

// PlayerConn represents a player's network connection.
type PlayerConn interface {
	// WritePacket sends a packet to the player.
	WritePacket(pk packet.Outbound)
	// Disconnect kicks the player from the server.
	Disconnect(reason *chat.Msg)
}

// Player represents a player entity.
type Player struct {
	id     ID
	name   string
	uid    uuid.UUID
	conn   PlayerConn
	world  *World
	pos    Pos
	chunks map[ChunkPos]*Chunk
}

// NewPlayer constructs a new Player.
// The created Player will not be associated with any World yet.
func NewPlayer(name string, uid uuid.UUID, conn PlayerConn) *Player {
	return &Player{
		id:     newEntityID(),
		name:   name,
		uid:    uid,
		conn:   conn,
		chunks: make(map[ChunkPos]*Chunk),
	}
}

// EntityID implements Entity.EntityID.
func (p *Player) EntityID() ID {
	return p.id
}

// Close releases resources associated with the Player.
func (p *Player) Close() error {
	p.SetWorld(nil)
	return nil
}

// SetWorld moves the player to a different World.
func (p *Player) SetWorld(w *World) {
	if p.world != nil {
		p.world.RemoveEntity(p.id)
	}

	p.world = w
	if w != nil {
		w.AddEntity(p)
	}
}

// Teleport moves the player to a new Pos.
func (p *Player) Teleport(pos Pos) {
	p.pos = pos
	p.updateChunks()
	p.conn.WritePacket(&outbound.Position{
		X:     pos.X,
		Y:     pos.Y + PlayerEyeHeight,
		Z:     pos.Z,
		Yaw:   pos.Yaw,
		Pitch: pos.Pitch,
	})
}

func (p *Player) tick() {
	p.conn.WritePacket(&outbound.KeepAlive{
		ID: 0,
	})
}

// updateChunks updates the chunks map for the player based on their current position.
func (p *Player) updateChunks() {
	// TODO properly calculate view distance
	const viewDist = 2
	center := ChunkPosFromWorldCoords(p.pos.X, p.pos.Z)

	// unload old chunks
	for pos := range p.chunks {
		if center.Dist(pos) > viewDist {
			p.conn.WritePacket(&outbound.ChunkData{
				X:         pos.X,
				Z:         pos.Z,
				FullChunk: true,
				Mask:      0,
			})
			delete(p.chunks, pos)
		}
	}

	pk := &outbound.BulkChunkData{
		SkyLightIncluded: true,
	}

	// load new chunks
	for x := center.X - viewDist; x <= center.X+viewDist; x++ {
		for z := center.Z - viewDist; z <= center.Z+viewDist; z++ {
			pos := ChunkPos{x, z}
			if _, loaded := p.chunks[pos]; loaded {
				continue
			}

			c := p.world.GetChunk(pos)
			if c != nil {
				pk.ChunkCount++
				pk.Meta = append(pk.Meta, outbound.BulkChunkDataMeta{
					X:           x,
					Z:           z,
					SectionMask: c.sectionMask,
				})
				pk.Data = c.appendData(pk.Data)
				p.chunks[pos] = c

				if pk.ChunkCount == 10 {
					p.conn.WritePacket(pk)
					pk = &outbound.BulkChunkData{
						SkyLightIncluded: true,
					}
				}
			}
		}
	}

	if pk.ChunkCount > 0 {
		p.conn.WritePacket(pk)
	}
}
