package game

import (
	"github.com/gitfyu/mable/biome"
	"github.com/gitfyu/mable/block"
	"math"
)

const (
	// chunkSectionBlocksSize is the number of bytes used for block data per chunkSection
	chunkSectionBlocksSize = 16 * 16 * 16 * 2

	// chunkSectionsPerChunk is the maximum number of chunkSection instances within a single Chunk
	chunkSectionsPerChunk = 16

	// lightDataSize is the number of bytes used for both block- and skylight data per chunkSection
	lightDataSize = 16 * 16 * 16 / 2 * 2

	// biomeDataSize is the number of bytes used for biome data in a single Chunk
	biomeDataSize = 256
)

var (
	// cachedLightData contains pre-generated block- and skylight data for a single full-sized chunk
	cachedLightData [lightDataSize * chunkSectionsPerChunk]byte

	// cachedBiomeData contains pre-generated biome data for a single chunk
	cachedBiomeData [biomeDataSize]byte
)

func init() {
	// I (currently) don't care about light or biome data for chunks, since it is not really relevant for Mable's
	// intended use case, which means that I can just pre-generate all this data and re-use it for every chunk, instead
	// of having to recompute it every time.

	const fullBright = 15
	for i := range cachedLightData {
		cachedLightData[i] = fullBright<<4 | fullBright
	}

	for i := range cachedBiomeData {
		cachedBiomeData[i] = uint8(biome.Plains)
	}
}

// ChunkPos contains a pair of chunk coordinates
type ChunkPos struct {
	X, Z int32
}

// ChunkPosFromWorldCoords returns the ChunkPos for the given world coordinates
func ChunkPosFromWorldCoords(x, z float64) ChunkPos {
	return ChunkPos{
		X: int32(math.Floor(x / 16)),
		Z: int32(math.Floor(z / 16)),
	}
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

// chunkSection represents a 16-block tall section within a chunk
type chunkSection [chunkSectionBlocksSize]byte

// Chunk represents a 16x16x256 area in a World
type Chunk struct {
	listeners map[uint32]chan<- interface{}

	// sectionMask is a bitmask where the nth bit indicates if sections[n] is set
	sectionMask uint16

	// sectionCount is the number of chunkSection instances stored in sections
	sectionCount int

	// sections contains all chunkSection instances for this Chunk. It is possible that not all indices contain a
	// chunkSection, in which case they will be nil.
	sections [chunkSectionsPerChunk]*chunkSection

	// dataSize is the total size required for the buffer that should be passed to writeData
	dataSize int
}

// NewChunk constructs a new Chunk
func NewChunk() *Chunk {
	return &Chunk{
		listeners: make(map[uint32]chan<- interface{}),
		dataSize:  biomeDataSize,
	}
}

// SetBlock changes a block in the chunk. Note that the coordinates are relative to the chunk, not world coordinates.
// Coordinates must all be within the range [0,15] or the function will panic.
func (c *Chunk) SetBlock(x, y, z uint8, data block.Data) {
	sectionIdx := y >> 4
	c.createSectionIfNotExists(sectionIdx)

	section := c.sections[sectionIdx]
	idx := int(y&15)<<9 | int(z)<<5 | int(x)<<1
	v := data.ToUint16()

	section[idx] = uint8(v)
	section[idx+1] = uint8(v >> 8)
}

// createSectionIfNotExists creates and stores a new chunkSection at the specified index if it does not exist yet
func (c *Chunk) createSectionIfNotExists(index uint8) {
	if c.sectionMask&(1<<index) != 0 {
		return
	}

	c.sectionCount++
	c.sectionMask |= 1 << index
	c.sections[index] = new(chunkSection)
	c.dataSize += chunkSectionBlocksSize + lightDataSize
}

// writeData will write the data for this chunk to the buffer, to be sent in a packet. The required size of the buffer
// is specified in Chunk.dataSize, a smaller buffer will cause the function to panic.
func (c *Chunk) writeData(buf []byte) {
	off := 0

	// blocks
	for i := 0; i < chunkSectionsPerChunk; i++ {
		if c.sectionMask&(1<<i) != 0 {
			copy(buf[off:], c.sections[i][:])
			off += chunkSectionBlocksSize
		}
	}

	// light
	lightDataSize := lightDataSize * c.sectionCount
	copy(buf[off:], cachedLightData[:lightDataSize])
	off += lightDataSize

	// biomes
	copy(buf[off:], cachedBiomeData[:])
}

// Subscribe registers the specified channel to receive updates for this Chunk. The specified ID must be unique to the
// subscriber.
func (c *Chunk) Subscribe(id uint32, ch chan<- interface{}) {
	ch <- "Subbed"
	c.listeners[id] = ch
}

// Unsubscribe cancels the subscription associated with the specified ID
func (c *Chunk) Unsubscribe(id uint32) {
	delete(c.listeners, id)
}

// Broadcast broadcasts a message to all subscribers of this Chunk
func (c *Chunk) Broadcast(msg interface{}) {
	for _, ch := range c.listeners {
		ch <- msg
	}
}
