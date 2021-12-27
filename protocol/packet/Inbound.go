package packet

import (
	"github.com/gitfyu/mable/protocol"
)

var idToPacket = make(map[uint]func() Inbound)

// RegisterInbound registers a supplier function that creates a packet.Inbound for a specific protocol.State and id
func RegisterInbound(state protocol.State, id uint, supplier func() Inbound) {
	idToPacket[id<<2|uint(state)] = supplier
}

// createInbound creates an Inbound instance for the specified protocol.State and ID, or nil in case no instance could
// be created
func createInbound(state protocol.State, id uint) Inbound {
	f, ok := idToPacket[id<<2|uint(state)]
	if !ok {
		return nil
	}

	return f()
}

type Inbound interface {
	UnmarshalPacket(r *protocol.ReadBuffer)
}
