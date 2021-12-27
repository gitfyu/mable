package world

import "testing"

func testSetBlock(t *testing.T, c *Chunk, x, y, z uint8, expect bool) {
	v := c.SetBlock(x, y, z, BlockData{0, 0})

	if v != expect {
		t.Errorf("Expected %t, got %t (X:%d, Y:%d, Z:%d)", v, expect, x, y, z)
	}
}

func TestChunk_SetBlock(t *testing.T) {
	c := NewChunk(1, 2)

	// TODO currently only really tests whether the section is computed properly
	testSetBlock(t, c, 0, 0, 0, false)
	testSetBlock(t, c, 0, 15, 0, false)
	testSetBlock(t, c, 0, 16, 0, true)
	testSetBlock(t, c, 0, 31, 0, true)
	testSetBlock(t, c, 0, 32, 0, true)
	testSetBlock(t, c, 0, 47, 0, true)
	testSetBlock(t, c, 0, 48, 0, false)
	testSetBlock(t, c, 0, 255, 0, false)
}
