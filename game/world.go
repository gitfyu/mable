package game

// World represents a world within the server.
type World struct {
	chunks   map[ChunkPos]*Chunk
	entities map[ID]Entity
}

// NewWorld constructs a new World containing predefined chunks.
func NewWorld(chunks map[ChunkPos]*Chunk) *World {
	return &World{
		chunks:   chunks,
		entities: make(map[ID]Entity),
	}
}

// AddEntity adds an Entity to the world.
func (w *World) AddEntity(e Entity) {
	w.entities[e.EntityID()] = e
}

// RemoveEntity removes the entity associated with the specified id.
// If no such entity exists, this function does nothing.
func (w *World) RemoveEntity(id ID) {
	delete(w.entities, id)
}

// GetChunk gets the Chunk at the specified position, or nil if it does not exist.
func (w *World) GetChunk(pos ChunkPos) *Chunk {
	return w.chunks[pos]
}

// tick updates the World.
func (w *World) tick() {
	for _, e := range w.entities {
		e.tick()
	}
}
