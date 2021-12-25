package entity

import (
	"github.com/gitfyu/mable/chat"
	"github.com/gitfyu/mable/protocol"
	"github.com/gitfyu/mable/protocol/packet"
	"github.com/gitfyu/mable/world"
	"github.com/gitfyu/mable/world/biome"
	"github.com/gitfyu/mable/world/block"
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
	WritePacket(id packet.ID, buf *packet.Buffer) error
	// Disconnect kicks the player from the server
	Disconnect(reason *chat.Msg) error
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
func (p *Player) SetPos(pos world.Pos) error {
	p.posLock.Lock()
	defer p.posLock.Unlock()

	p.pos = pos

	buf := packet.AcquireBuffer()
	defer packet.ReleaseBuffer(buf)

	buf.WriteDouble(p.pos.X)
	buf.WriteDouble(p.pos.Y + PlayerEyeHeight)
	buf.WriteDouble(p.pos.Z)
	buf.WriteFloat(p.pos.Yaw)
	buf.WriteFloat(p.pos.Pitch)

	// flags indicating all values are absolute
	buf.WriteSignedByte(0)

	return p.conn.WritePacket(packet.PlayServerPosAndLook, buf)
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
			_ = p.Ping()
		case <-p.destroyed:
			return
		}
	}
}

// TODO currently the actual data being sent is hardcoded, in the future it should be passed as a parameter

func (p *Player) SendChunkData(chunkX, chunkZ int32) error {
	buf := packet.AcquireBuffer()
	defer packet.ReleaseBuffer(buf)

	buf.WriteInt(chunkX)
	buf.WriteInt(chunkZ)
	// true means full chunk
	buf.WriteBool(true)

	// mask, first bit set means only the lowest section is sent
	buf.WriteUnsignedShort(1)
	buf.WriteVarInt(protocol.VarInt(protocol.ChunkDataSize(1)))

	// blocks
	for y := 0; y < 16; y++ {
		for z := 0; z < 16; z++ {
			for x := 0; x < 16; x++ {
				buf.WriteUnsignedShortLittleEndian(protocol.EncodeBlockData(block.Stone, 0))
			}
		}
	}

	// block light
	for i := 0; i < protocol.LightDataSize; i++ {
		buf.WriteUnsignedByte(protocol.FullBright<<4 | protocol.FullBright)
	}

	// skylight
	for i := 0; i < protocol.LightDataSize; i++ {
		buf.WriteUnsignedByte(protocol.FullBright<<4 | protocol.FullBright)
	}

	// biomes
	for i := 0; i < 256; i++ {
		buf.WriteUnsignedByte(uint8(biome.Plains))
	}

	return p.conn.WritePacket(packet.PlayServerChunkData, buf)
}

func (p *Player) Ping() error {
	buf := packet.AcquireBuffer()
	defer packet.ReleaseBuffer(buf)

	// TODO currently the same arbitrary ID is sent every time, since the server has no use for the response (yet)
	buf.WriteVarInt(0)

	return p.conn.WritePacket(packet.PlayServerKeepAlive, buf)
}
