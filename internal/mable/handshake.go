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

func handleHandshake(h *connHandler, data *packet.Buffer) error {
	ver, err := data.ReadVarInt()
	if err != nil {
		return err
	}
	h.version = protocol.Version(ver)

	// address
	if _, err := data.ReadString(); err != nil {
		return err
	}

	// port
	if _, err := data.ReadUnsignedShort(); err != nil {
		return err
	}

	nextState, err := data.ReadVarInt()
	if err != nil {
		return err
	}

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
