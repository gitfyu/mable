package game

import (
	"github.com/gitfyu/mable/block"
)

// TODO this only exists for testing, will be removed soon

// DefaultWorld is a global default world.
var DefaultWorld = createDefaultWorld()

func createDefaultWorld() *World {
	chunks := make(map[ChunkPos]*Chunk)
	for x := int32(-1); x <= 1; x++ {
		for z := int32(-1); z <= 1; z++ {
			c := NewChunk()
			for dx := uint8(0); dx < 16; dx++ {
				for dz := uint8(0); dz < 16; dz++ {
					for dy := uint8(1); dy < 100; dy += 5 {
						c.SetBlock(dx, dy, dz, block.Stone.ToData())
					}
				}
			}

			chunks[ChunkPos{x, z}] = c
		}
	}
	return NewWorld(chunks)
}
