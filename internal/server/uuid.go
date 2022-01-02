package server

import (
	"crypto/md5"
	"github.com/google/uuid"
)

// generateOfflineUUID generates a UUID from a username in the same way as the vanilla server
func generateOfflineUUID(username string) uuid.UUID {
	b := md5.Sum([]byte("OfflinePlayer:" + username))
	b[6] &= 0x0f
	b[6] |= 0x30
	b[8] &= 0x3f
	b[8] |= 0x80
	u, _ := uuid.FromBytes(b[:])
	return u
}
