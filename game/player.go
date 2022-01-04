package game

import (
	"github.com/gitfyu/mable/chat"
	"github.com/gitfyu/mable/internal/protocol/packet"
	outbound "github.com/gitfyu/mable/internal/protocol/packet/outbound/play"
	"github.com/google/uuid"
	"sync"
)

const PlayerEyeHeight = 1.62

// PlayerConn represents a player's network connection
type PlayerConn interface {
	// WritePacket sends a packet to the player
	WritePacket(pk packet.Outbound)
	// Disconnect kicks the player from the server
	Disconnect(reason *chat.Msg)
}

// Player represents a player entity
type Player struct {
	id        ID
	name      string
	uid       uuid.UUID
	conn      PlayerConn
	world     *World
	worldLock sync.RWMutex
	pos       Pos
	chunks    map[ChunkPos]*Chunk
}

// NewPlayer constructs a new player
func NewPlayer(name string, uid uuid.UUID, conn PlayerConn, w *World) *Player {
	p := &Player{
		id:     newEntityID(),
		name:   name,
		uid:    uid,
		conn:   conn,
		world:  w,
		chunks: make(map[ChunkPos]*Chunk),
	}
	w.AddEntity(p)
	return p
}

func (p *Player) EntityID() ID {
	return p.id
}

// Close releases resources associated with the Player. This function should only be called once and will always return
// nil.
func (p *Player) Close() error {
	p.SetWorld(nil)
	return nil
}

func (p *Player) SetWorld(w *World) {
	p.worldLock.Lock()
	defer p.worldLock.Unlock()

	if p.world != nil {
		p.world.RemoveEntity(p.id)
	}

	p.world = w
	if w != nil {
		w.AddEntity(p)
	}
}

// Teleport changes the player' position
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

// sendChunkData sends a chunk data packet to the player. To unload a chunk, set the chunk parameter to nil.
func (p *Player) sendChunkData(chunkX, chunkZ int32, chunk *Chunk) {
	var mask uint16
	var data []byte

	if chunk != nil {
		mask = chunk.sectionMask
		data = make([]byte, chunk.dataSize)
		chunk.writeData(data)
	} else {
		mask = 0
		data = []byte{}
	}

	p.conn.WritePacket(&outbound.ChunkData{
		X:         chunkX,
		Z:         chunkZ,
		FullChunk: true,
		Mask:      mask,
		Data:      data,
	})
}

func (p *Player) tick() {
	p.conn.WritePacket(&outbound.KeepAlive{
		ID: 0,
	})
}

// updateChunks updates the chunks map for the player based on their current position
func (p *Player) updateChunks() {
	p.worldLock.Lock()
	defer p.worldLock.Unlock()

	// TODO properly calculate view distance
	const viewDist = 2
	center := ChunkPosFromWorldCoords(p.pos.X, p.pos.Z)

	// unload old chunks
	for pos := range p.chunks {
		if center.Dist(pos) > viewDist {
			p.sendChunkData(pos.X, pos.Z, nil)
			delete(p.chunks, pos)
		}
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
				p.sendChunkData(pos.X, pos.Z, c)
				p.chunks[pos] = c
			}
		}
	}
}
