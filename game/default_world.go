package game

import (
	"github.com/gitfyu/mable/block"
)

// TODO this only exists for testing, will be removed soon

// DefaultWorld is a global default world
var DefaultWorld = createDefaultWorld()

func createDefaultWorld() *World {
	c := NewChunk()
	for x := uint8(0); x < 16; x++ {
		for z := uint8(0); z < 16; z++ {
			for y := uint8(1); y < 100; y += 5 {
				c.SetBlock(x, y, z, BlockData{block.Stone, 0})
			}
		}
	}

	return NewWorld(map[ChunkPos]*Chunk{
		ChunkPos{0, 0}: c,
	})
}
