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
	// Version returns the protocol.Version of the player's connection
	Version() protocol.Version
	// WritePacket sends a packet to the player
	WritePacket(id packet.ID, buf *packet.Buffer) error
	// Disconnect kicks the player from the server
	Disconnect(reason *chat.Msg) error
}

// Player represents a player entity
type Player struct {
	Entity
	name         string
	uid          uuid.UUID
	conn         PlayerConn
	world        *world.World
	pos          world.Pos
	worldPosLock sync.RWMutex
	pings        chan int32
}

// NewPlayer constructs a new player
func NewPlayer(name string, uid uuid.UUID, conn PlayerConn, w *world.World) *Player {
	return &Player{
		Entity: NewEntity(),
		name:   name,
		uid:    uid,
		conn:   conn,
		world:  w,
	}
}

// Update will keep pinging the player until the context is cancelled
func (p *Player) Update(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	worldUpdates := make(chan interface{}, 10)
	sub := p.world.Subscribe(worldUpdates)
	defer p.world.Unsubscribe(sub)

	for {
		select {
		case msg := <-worldUpdates:
			log.Debug().
				Str("receiver", p.name).
				Interface("msg", msg).
				Msg("World update received")
		case <-ticker.C:
			_ = p.Ping()
		case <-ctx.Done():
			return
		}
	}
}

// Teleport moves the player to the given position
func (p *Player) Teleport(pos world.Pos) error {
	buf := packet.AcquireBuffer()
	defer packet.ReleaseBuffer(buf)

	buf.WriteDouble(pos.X)
	buf.WriteDouble(pos.Y + PlayerEyeHeight)
	buf.WriteDouble(pos.Z)
	buf.WriteFloat(pos.Yaw)
	buf.WriteFloat(pos.Pitch)

	if p.conn.Version() == protocol.Version_1_8 {
		// flags indicating all values are absolute
		buf.WriteSignedByte(0)
	} else {
		// on ground, useless
		buf.WriteBool(false)
	}

	p.worldPosLock.Lock()
	defer p.worldPosLock.Unlock()

	p.pos = pos

	return p.conn.WritePacket(packet.PlayServerPosAndLook, buf)
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
	if p.conn.Version() == protocol.Version_1_8 {
		buf.WriteVarInt(0)
	} else {
		buf.WriteInt(0)
	}

	return p.conn.WritePacket(packet.PlayServerKeepAlive, buf)
}
