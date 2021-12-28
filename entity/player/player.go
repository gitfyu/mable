package player

import (
	"github.com/gitfyu/mable/chat"
	"github.com/gitfyu/mable/entity"
	"github.com/gitfyu/mable/protocol"
	"github.com/gitfyu/mable/protocol/packet"
	"github.com/gitfyu/mable/protocol/packet/play"
	"github.com/gitfyu/mable/world"
	"github.com/gitfyu/mable/world/biome"
	"github.com/google/uuid"
	"sync"
)

const EyeHeight = 1.62

// Conn represents a player's network connection
type Conn interface {
	// WritePacket sends a packet to the player
	WritePacket(pk packet.Outbound)
	// Disconnect kicks the player from the server
	Disconnect(reason *chat.Msg)
}

// Player represents a player entity
type Player struct {
	id        entity.ID
	name      string
	uid       uuid.UUID
	conn      Conn
	world     *world.World
	worldLock sync.RWMutex
	pos       world.Pos
	destroyed chan struct{}
}

// NewPlayer constructs a new player
func NewPlayer(name string, uid uuid.UUID, conn Conn, w *world.World) *Player {
	p := &Player{
		id:        entity.NewID(),
		name:      name,
		uid:       uid,
		conn:      conn,
		world:     w,
		destroyed: make(chan struct{}),
	}
	w.AddEntity(p)
	return p
}

func (p *Player) GetEntityID() entity.ID {
	return p.id
}

// Destroy should be called to clean up resources when this Player is no longer needed. The Player should not be used
// again afterwards.
func (p *Player) Destroy() {
	p.SetWorld(nil)
	close(p.destroyed)
}

func (p *Player) SetWorld(w *world.World) {
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

func (p *Player) Tick() {
	p.conn.WritePacket(&play.OutKeepAlive{
		ID: 0,
	})
}

// Teleport changes the player' position
func (p *Player) Teleport(pos world.Pos) {
	p.pos = pos
	p.conn.WritePacket(&play.OutPosition{
		X:     pos.X,
		Y:     pos.Y + EyeHeight,
		Z:     pos.Z,
		Yaw:   pos.Yaw,
		Pitch: pos.Pitch,
	})
}

// TODO currently the actual data being sent is hardcoded, in the future it should be passed as a parameter

func (p *Player) SendChunkData(chunkX, chunkZ int32, c *world.Chunk) {
	pk := play.OutChunkData{
		X:         chunkX,
		Z:         chunkZ,
		FullChunk: true,
		Mask:      1,
		Data:      make([]byte, protocol.ChunkDataSize(1)),
	}

	// TODO this code currently assumes that the chunk will only write one section, which is not always the case
	c.WriteBlocks(pk.Data)
	off := 16 * 16 * 16 * 2

	// block light
	for i := 0; i < protocol.LightDataSize; i++ {
		pk.Data[off+i] = protocol.FullBright<<4 | protocol.FullBright
	}
	off += protocol.LightDataSize

	// skylight
	for i := 0; i < protocol.LightDataSize; i++ {
		pk.Data[off+i] = protocol.FullBright<<4 | protocol.FullBright
	}
	off += protocol.LightDataSize

	// biomes
	for i := 0; i < 256; i++ {
		pk.Data[off+i] = uint8(biome.Plains)
	}

	p.conn.WritePacket(&pk)
}
