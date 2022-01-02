package packet

import (
	"github.com/gitfyu/mable/internal/protocol"
)

// Inbound represents a packet sent from the client
type Inbound interface {
	UnmarshalPacket(r *protocol.ReadBuffer)
}

// Outbound represents a packet sent from the server
type Outbound interface {
	PacketID() uint
	MarshalPacket(w *protocol.WriteBuffer)
}

var idToPacket = make(map[uint]func() Inbound)

// RegisterInbound registers a supplier function that creates an Inbound packet for a specific protocol.State and id
func RegisterInbound(state protocol.State, id uint, supplier func() Inbound) {
	idToPacket[id<<2|uint(state)] = supplier
}

// createInbound creates an Inbound packet for the specified protocol.State and ID, or nil in case no instance could be
// created
func createInbound(state protocol.State, id uint) Inbound {
	f, ok := idToPacket[id<<2|uint(state)]
	if !ok {
		return nil
	}

	return f()
}
