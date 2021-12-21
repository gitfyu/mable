package chunk

import "github.com/gitfyu/mable/world/biome"

const BiomeDataSize = 256

// biomeDataAllPlains is a pre-computed array containing only plains biomes, sent in the chunk data packet
var biomeDataAllPlains []byte

func init() {
	biomeDataAllPlains = make([]byte, BiomeDataSize)
	for i := 0; i < BiomeDataSize; i++ {
		biomeDataAllPlains[i] = byte(biome.Plains)
	}
}
