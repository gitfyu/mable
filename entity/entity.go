package entity

import "sync/atomic"

type ID int32

var entityIdCounter int32

// GenId generates a new entity ID
func GenId() ID {
	return ID(atomic.AddInt32(&entityIdCounter, 1))
}

// Entity is the base type for all entities
type Entity struct {
	id ID
}

func NewEntity() Entity {
	return Entity{id: GenId()}
}

func (e *Entity) GetEntityID() ID {
	return e.id
}
