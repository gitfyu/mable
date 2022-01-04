package game

import "sync/atomic"

var entityIdCounter int32

type ID int32

type Entity interface {
	EntityID() ID
	tick()
}

func newEntityID() ID {
	return ID(atomic.AddInt32(&entityIdCounter, 1))
}
