package mable

import (
	"errors"
	"github.com/gitfyu/mable/protocol"
	"github.com/gitfyu/mable/protocol/packet"
)

var errInvalidState = errors.New("invalid state")

var handshakeHandlers = newPacketHandlerLookup(
	packetHandlers{
		packet.Handshake: handleHandshake,
	},
)

func handleHandshake(h *connHandler, p *packet.Packet) error {
	// protocol version
	h.version = protocol.Version(p.GetVarInt())
	// address
	_ = p.GetString()
	// port
	_ = p.GetUnsignedShort()
	nextState := p.GetVarInt()

	switch s := protocol.State(nextState); s {
	case protocol.StateStatus:
		fallthrough
	case protocol.StateLogin:
		// TODO ensure the protocol version is supported
		h.state = s
	default:
		return errInvalidState
	}

	return nil
}
