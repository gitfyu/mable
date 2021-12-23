package entity

import (
	"context"
	"github.com/gitfyu/mable/chat"
	"github.com/gitfyu/mable/protocol"
	"github.com/gitfyu/mable/protocol/chunk"
	"github.com/gitfyu/mable/protocol/packet"
	"github.com/gitfyu/mable/world"
	"github.com/gitfyu/mable/world/biome"
	"github.com/gitfyu/mable/world/block"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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
	name      string
	uid       uuid.UUID
	conn      PlayerConn
	pos       world.Pos
	posLock   sync.RWMutex
	worldLeft *sync.Cond
	pings     chan int32
}

// NewPlayer constructs a new player
func NewPlayer(name string, uid uuid.UUID, conn PlayerConn) *Player {
	p := &Player{
		Entity: NewEntity(),
		name:   name,
		uid:    uid,
		conn:   conn,
	}
	p.worldLeft = sync.NewCond(&p.posLock)
	return p
}

// Destroy should be called to clean up resources when this Player is no longer needed. The Player should not be used
// again afterwards.
func (p *Player) Destroy() {
	p.posLock.Lock()
	defer p.posLock.Unlock()

	p.leaveWorld()
}

// SetPos changes the player' position
func (p *Player) SetPos(pos world.Pos) error {
	p.posLock.Lock()
	defer p.posLock.Unlock()

	if p.pos.World != pos.World {
		p.leaveWorld()
		p.enterWorld(pos.World)
	}

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

// leaveWorld removes the player from their current world, if they are in one. The calling goroutine MUST hold the
// Player.posLock!
func (p *Player) leaveWorld() {
	if p.pos.World != nil {
		p.pos.World.Unsubscribe(world.SubscriberID(p.GetEntityID()))
		// wait for all updates from previous world to be processed
		p.worldLeft.Wait()
	}
}

// enterWorld adds a player to a new world. The calling goroutine MUST hold the Player.posLock!
func (p *Player) enterWorld(w *world.World) {
	if w != nil {
		ch := w.Subscribe(world.SubscriberID(p.GetEntityID()))
		go func() {
			for msg := range ch {
				log.Debug().
					Str("receiver", p.name).
					Interface("msg", msg).
					Msg("World update received")
			}
			// signal that all updates have been processed
			p.worldLeft.Signal()
		}()
	}
}

// KeepAlive will keep pinging the player until the context is cancelled
func (p *Player) KeepAlive(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_ = p.Ping()
		case <-ctx.Done():
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
	buf.WriteVarInt(protocol.VarInt(chunk.TotalDataSize(1)))

	// blocks
	for y := 0; y < 16; y++ {
		for z := 0; z < 16; z++ {
			for x := 0; x < 16; x++ {
				buf.WriteUnsignedShortLittleEndian(chunk.EncodeBlockData(block.Stone, 0))
			}
		}
	}

	// block light
	for i := 0; i < chunk.LightDataSize; i++ {
		buf.WriteUnsignedByte(chunk.FullBright<<4 | chunk.FullBright)
	}

	// skylight
	for i := 0; i < chunk.LightDataSize; i++ {
		buf.WriteUnsignedByte(chunk.FullBright<<4 | chunk.FullBright)
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
