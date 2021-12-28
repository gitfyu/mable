package entity

import (
	"github.com/gitfyu/mable/chat"
	"github.com/gitfyu/mable/protocol"
	"github.com/gitfyu/mable/protocol/packet"
	"github.com/gitfyu/mable/protocol/packet/play"
	"github.com/gitfyu/mable/world"
	"github.com/gitfyu/mable/world/biome"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"math"
	"sync"
	"time"
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
	Entity
	name string
	uid  uuid.UUID
	conn PlayerConn

	// world- and position related values
	pos          world.Pos
	posLock      sync.RWMutex
	chunkUpdates chan interface{}

	pings     chan int32
	destroyed chan struct{}
}

// NewPlayer constructs a new player
func NewPlayer(name string, uid uuid.UUID, conn PlayerConn) *Player {
	p := &Player{
		Entity:       NewEntity(),
		name:         name,
		uid:          uid,
		conn:         conn,
		destroyed:    make(chan struct{}),
		chunkUpdates: make(chan interface{}, 100),
	}

	go p.keepAlive()
	go p.updateChunks()

	return p
}

// Destroy should be called to clean up resources when this Player is no longer needed. The Player should not be used
// again afterwards.
func (p *Player) Destroy() {
	close(p.destroyed)
}

// setCoords updates the player's coordinates, without sending an update to the client
func (p *Player) setCoords(x, y, z float64) {
	p.posLock.Lock()
	defer p.posLock.Unlock()

	p.pos.X, p.pos.Y, p.pos.Z = x, y, z
}

// SetPos changes the player' position
func (p *Player) SetPos(pos world.Pos) {
	p.posLock.Lock()
	defer p.posLock.Unlock()

	p.pos = pos
	p.conn.WritePacket(&play.OutPosition{
		X:     pos.X,
		Y:     pos.Y + PlayerEyeHeight,
		Z:     pos.Z,
		Yaw:   pos.Yaw,
		Pitch: pos.Pitch,
	})
}

func (p *Player) GetChunkPos() world.ChunkPos {
	p.posLock.RLock()
	defer p.posLock.RUnlock()

	return world.ChunkPos{
		X: int32(math.Floor(p.pos.X / 16)),
		Z: int32(math.Floor(p.pos.Z / 16)),
	}
}

func (p *Player) updateChunks() {
	// TODO configurable chunk update rate
	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()

	// TODO
	const viewDist = 2

	subId := uint32(p.GetEntityID())
	chunks := make(map[world.ChunkPos]*world.Chunk)

	for {
		select {
		case <-ticker.C:
			center := p.GetChunkPos()

			// Unload chunks that have gone out of range
			for pos, c := range chunks {
				if pos.Dist(center) > viewDist {
					c.Unsubscribe(subId)
					delete(chunks, pos)
				}
			}

			// Load new chunks
			for x := center.X - viewDist; x <= center.X+viewDist; x++ {
				for z := center.Z - viewDist; z <= center.Z+viewDist; z++ {
					pos := world.ChunkPos{X: x, Z: z}
					if _, loaded := chunks[pos]; !loaded {
						c := p.pos.World.GetChunk(pos)
						if c != nil {
							c.Subscribe(subId, p.chunkUpdates)
							chunks[pos] = c

							p.SendChunkData(x, z, c)
						}
					}
				}
			}
		case msg := <-p.chunkUpdates:
			log.Debug().Interface("msg", msg).Msg("Chunk update")
		case <-p.destroyed:
			for _, c := range chunks {
				c.Unsubscribe(subId)
			}
			return
		}
	}
}

// keepAlive will frequently ping the player to prevent them from disconnecting
func (p *Player) keepAlive() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.Ping()
		case <-p.destroyed:
			return
		}
	}
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

func (p *Player) Ping() {
	p.conn.WritePacket(&play.OutKeepAlive{
		ID: 0,
	})
}
