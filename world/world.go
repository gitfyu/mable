package world

type SubscriberID int32

type World struct {
	chunks map[ChunkPos]*Chunk
}

// NewWorld constructs a new World
func NewWorld(chunks map[ChunkPos]*Chunk) *World {
	return &World{
		chunks: chunks,
	}
}

// GetChunk gets the Chunk at the specified position, or nil if it does not exist
func (w *World) GetChunk(pos ChunkPos) *Chunk {
	return w.chunks[pos]
}
