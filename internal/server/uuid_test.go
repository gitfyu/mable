package server

import (
	"bytes"
	"github.com/google/uuid"
	"testing"
)

func Test_generateOfflineUUID(t *testing.T) {
	got := generateOfflineUUID("test123")
	expect, err := uuid.Parse("be4c4b88-c56b-3b93-aec4-4bc0d038a924")
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(got[:], expect[:]) {
		t.Errorf("Expected %d, got %d", got, expect)
	}
}
