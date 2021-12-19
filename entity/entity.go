package entity

import "sync/atomic"

type ID int32

var entityIdCounter int32

// GenId generates a new entity ID
func GenId() ID {
	return ID(atomic.AddInt32(&entityIdCounter, 1))
}
