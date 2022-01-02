package protocol

import (
	"github.com/gitfyu/mable/block"
)

// ChunkDataSize returns the number of bytes that all the chunk data will consist of for the given number of sections
func ChunkDataSize(sectionCount int) int {
	blockDataSize := sectionCount * 2 * 16 * 16 * 16
	blockLightSize := sectionCount * 16 * 16 * 16 / 2
	skyLightSize := sectionCount * 16 * 16 * 16 / 2
	biomesSize := 256
	return blockDataSize + blockLightSize + skyLightSize + biomesSize
}

// EncodeBlockData packs a block.ID and additional data into a single uint16
func EncodeBlockData(id block.ID, data int) uint16 {
	return uint16(id)<<4 | uint16(data)&16
}
