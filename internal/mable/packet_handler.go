package mable

import (
	"github.com/gitfyu/mable/protocol/packet"
)

type packetHandler func(h *conn, data *packet.Buffer) error

// packetHandlerLookup acts as a map with a packet ID as key and a packetHandler as value
type packetHandlerLookup []packetHandler

type packetHandlers map[packet.ID]packetHandler

// newPacketHandlerLookup constructs a packetHandlerLookup from a map, where the keys are packet.ID values. This
// function will panic for invalid packet.ID values.
func newPacketHandlerLookup(h packetHandlers) packetHandlerLookup {
	maxId := packet.ID(0)
	for k, _ := range h {
		if k > maxId {
			maxId = k
		} else if k < 0 {
			panic("illegal packet ID")
		}
	}

	lookup := make(packetHandlerLookup, maxId+1)
	for k, v := range h {
		lookup[k] = v
	}

	return lookup
}
