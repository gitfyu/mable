package chunk

const LightDataSize = 2048

// lightDataFullBright is a pre-computed array containing maximum brightness levels, sent in the chunk data packet
var lightDataFullBright []byte

func init() {
	lightDataFullBright = make([]byte, LightDataSize)
	for i := 0; i < LightDataSize; i++ {
		lightDataFullBright[i] = 15<<4 | 15
	}
}
