package world

import "github.com/gitfyu/mable/world/block"

// TODO this only exists for testing, will be removed soon

// Default is a global default world
var Default = createDefaultWorld()

func createDefaultWorld() *World {
	c := NewChunk(0, 0)
	for x := uint8(0); x < 16; x++ {
		for z := uint8(0); z < 16; z++ {
			c.SetBlock(x, 1, z, BlockData{block.Stone, 0})
		}
	}

	return NewWorld(map[ChunkPos]*Chunk{
		ChunkPos{0, 0}: c,
	})
}
