package world

import (
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

// Chunk represents a 16x16x256 area in a World
type Chunk struct {
	listeners     map[uint32]chan<- interface{}
	listenersLock sync.RWMutex
}

// NewChunk constructs a new Chunk
func NewChunk() *Chunk {
	return &Chunk{
		listeners: make(map[uint32]chan<- interface{}),
	}
}

// Subscribe registers the specified channel to receive updates for this Chunk. The specified ID must be unique to the
// subscriber.
func (c *Chunk) Subscribe(id uint32, ch chan<- interface{}) {
	c.listenersLock.Lock()
	defer c.listenersLock.Unlock()

	ch <- "Subbed"
	c.listeners[id] = ch
}

// Unsubscribe cancels the subscription associated with the specified ID
func (c *Chunk) Unsubscribe(id uint32) {
	c.listenersLock.Lock()
	defer c.listenersLock.Unlock()

	delete(c.listeners, id)
}

// Broadcast broadcasts a message to all subscribers of this Chunk
func (c *Chunk) Broadcast(msg interface{}) {
	c.listenersLock.RLock()
	defer c.listenersLock.RUnlock()

	for _, ch := range c.listeners {
		ch <- msg
	}
}
