package block

import (
	"testing"
)

func TestData_Type(t *testing.T) {
	for id := ID(0); id < maxID; id++ {
		d := id.ToData()
		if d.Type() != id {
			t.Errorf("Expected %d, got %d", id, d.Type())
		}
	}
}

func TestData_Metadata(t *testing.T) {
	for meta := uint8(0); meta < maxMetadata; meta++ {
		d := Stone.ToDataWithMetadata(meta)
		if d.Metadata() != meta {
			t.Errorf("Expected %d, got %d", meta, d.Metadata())
		}
	}
}
