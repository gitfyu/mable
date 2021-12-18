package mable

import (
	"errors"
	"github.com/gitfyu/mable/network"
	"github.com/gitfyu/mable/network/protocol"
	"github.com/gitfyu/mable/network/protocol/packet"
)

var errInvalidState = errors.New("invalid state")

var handshakeHandlers = newPacketHandlerLookup(
	packetHandlers{
		packet.Handshake: handleHandshake,
	},
)

func handleHandshake(h *connHandler, data *network.PacketData) error {
	// protocol version
	h.version = protocol.Version(data.GetVarInt())
	// address
	_ = data.GetString()
	// port
	_ = data.GetUnsignedShort()
	nextState := data.GetVarInt()

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
