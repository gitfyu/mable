package world

import (
	"github.com/gitfyu/mable/protocol/packet"
	"github.com/gitfyu/mable/world/block"
	"sync"
)

// ChunkPos contains a pair of chunk coordinates
type ChunkPos struct {
	X, Z int32
}

// Dist returns the distance between two ChunkPos values, in chunks
func (p ChunkPos) Dist(other ChunkPos) int32 {
	dx := p.X - other.X
	if dx < 0 {
		dx = -dx
	}
	dz := p.Z - other.Z
	if dz < 0 {
		dz = -dz
	}

	if dx > dz {
		return dx
	} else {
		return dz
	}
}

// BlockData contains a 12-bit block ID and a 4-bit data value
type BlockData struct {
	id   block.ID
	data uint8
}

// toUint16 encodes a BlockData value to an uint16, to be used in packets
func (b BlockData) toUint16() uint16 {
	return uint16(b.id)<<4 | uint16(b.data)&16
}

// Chunk represents a 16x16x256 area in a World
type Chunk struct {
	listeners  map[uint32]chan<- interface{}
	lock       sync.RWMutex
	minSection uint8
	maxSection uint8
	blocks     []uint8
}

// NewChunk constructs a new Chunk, where minSection is the index of the lowest section and maxSection the index of the
// highest section (index 0 refers to Y 0-15, index 1 refers to Y 16-31, etc.)
func NewChunk(minSection uint8, maxSection uint8) *Chunk {
	return &Chunk{
		listeners:  make(map[uint32]chan<- interface{}),
		minSection: minSection,
		maxSection: maxSection,
		blocks:     make([]uint8, 2*16*16*16*int(maxSection-minSection+1)),
	}
}

// SetBlock changes a block in the chunk. Note that the coordinates are relative to the chunk, not world coordinates.
// The function returns true if the block place was successful, false otherwise.
func (c *Chunk) SetBlock(x, y, z uint8, data BlockData) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	sectionIdx := y >> 4
	if sectionIdx < c.minSection || sectionIdx > c.maxSection {
		return false
	}

	idx := int(sectionIdx-c.minSection)<<13 | int(y&15)<<9 | int(z)<<5 | int(x)<<1
	v := data.toUint16()

	c.blocks[idx] = uint8(v)
	c.blocks[idx+1] = uint8(v >> 8)
	return true
}

func (c *Chunk) WriteBlocks(buf *packet.Buffer) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	buf.Write(c.blocks)
}

// Subscribe registers the specified channel to receive updates for this Chunk. The specified ID must be unique to the
// subscriber.
func (c *Chunk) Subscribe(id uint32, ch chan<- interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	ch <- "Subbed"
	c.listeners[id] = ch
}

// Unsubscribe cancels the subscription associated with the specified ID
func (c *Chunk) Unsubscribe(id uint32) {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.listeners, id)
}

// Broadcast broadcasts a message to all subscribers of this Chunk
func (c *Chunk) Broadcast(msg interface{}) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	for _, ch := range c.listeners {
		ch <- msg
	}
}
