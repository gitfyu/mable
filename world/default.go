package world

// TODO this only exists for testing, will be removed soon

// Default is a global default world
var Default = NewWorld(map[ChunkPos]*Chunk{
	ChunkPos{0, 0}: NewChunk(),
	ChunkPos{0, 1}: NewChunk(),
	ChunkPos{1, 0}: NewChunk(),
	ChunkPos{1, 1}: NewChunk(),
})
