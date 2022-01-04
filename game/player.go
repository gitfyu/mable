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
}

// NewPlayer constructs a new player
func NewPlayer(name string, uid uuid.UUID, conn PlayerConn, w *World) *Player {
	p := &Player{
		id:    newEntityID(),
		name:  name,
		uid:   uid,
		conn:  conn,
		world: w,
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

func (p *Player) tick() {
	p.conn.WritePacket(&outbound.KeepAlive{
		ID: 0,
	})
}

// Teleport changes the player' position
func (p *Player) Teleport(pos Pos) {
	p.pos = pos
	p.conn.WritePacket(&outbound.Position{
		X:     pos.X,
		Y:     pos.Y + PlayerEyeHeight,
		Z:     pos.Z,
		Yaw:   pos.Yaw,
		Pitch: pos.Pitch,
	})
}

// TODO currently the actual data being sent is hardcoded, in the future it should be passed as a parameter

func (p *Player) SendChunkData(chunkX, chunkZ int32, c *Chunk) {
	pk := outbound.ChunkData{
		X:         chunkX,
		Z:         chunkZ,
		FullChunk: true,
		Mask:      c.sectionMask,
		Data:      make([]byte, c.dataSize),
	}

	c.writeData(pk.Data)
	p.conn.WritePacket(&pk)
}
