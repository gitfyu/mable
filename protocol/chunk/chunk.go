package chunk

import (
	"github.com/gitfyu/mable/protocol"
	"github.com/gitfyu/mable/protocol/packet"
	"github.com/gitfyu/mable/world/block"
	"math/bits"
)

// SectionMask is a 16-bit value, where every bit refers to a 16x16x16 section in a chunk. The least-significant bit
// refers to the lowest section (y=0), the most-significant bit refers to the highest (y=240).
type SectionMask uint16

// count returns the number of sections in the mask
func (s SectionMask) count() int {
	return bits.OnesCount16(uint16(s))
}

// totalDataSize returns the number of bytes that all the chunk data will consist of
func (s SectionMask) totalDataSize() int {
	i := s.count()
	j := i * 2 * 16 * 16 * 16
	k := i * 16 * 16 * 16 / 2
	l := i * 16 * 16 * 16 / 2
	i1 := 256
	return j + k + l + i1
}

func encodeBlockData(id block.ID, data int) uint16 {
	return uint16(id)<<4 | uint16(data)&16
}

func WriteChunkData(buf *packet.Buffer, mask SectionMask) {
	buf.WriteVarInt(protocol.VarInt(mask.totalDataSize()))

	// number of sections
	n := mask.count()

	// blocks
	for i := 0; i < n; i++ {
		for y := 0; y < 16; y++ {
			for z := 0; z < 16; z++ {
				for x := 0; x < 16; x++ {
					buf.WriteUnsignedShortLittleEndian(encodeBlockData(block.Stone, 0))
				}
			}
		}
	}

	// block light
	for i := 0; i < n; i++ {
		buf.Write(lightDataFullBright)
	}

	// skylight
	for i := 0; i < n; i++ {
		buf.Write(lightDataFullBright)
	}

	// biomes
	buf.Write(biomeDataAllPlains)
}
