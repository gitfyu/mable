package game

import "sync/atomic"

// entityIdCounter holds the latest ID that was used.
var entityIdCounter int32

// ID is a unique identifier for an Entity.
type ID int32

// Entity represents a Minecraft entity, such as a player or a dropped item.
type Entity interface {
	// EntityID returns the ID for this entity. This function may be called concurrently.
	EntityID() ID

	tick()
}

// newEntityID generates a new ID. This function may be called concurrently.
func newEntityID() ID {
	return ID(atomic.AddInt32(&entityIdCounter, 1))
}
