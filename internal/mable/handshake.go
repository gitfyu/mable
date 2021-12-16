package mable

import (
	"errors"
	"github.com/gitfyu/mable/network/protocol"
)

var errInvalidState = errors.New("invalid state")

var handshakeHandlers = idToPacketHandler{
	handleHandshake,
}

func handleHandshake(_ int, h *connHandler) error {
	var protoVer protocol.VarInt
	var addr string
	var port uint16
	var nextState protocol.VarInt

	ok := h.dec.ReadVarInt(&protoVer) &&
		h.dec.ReadString(&addr) &&
		h.dec.ReadUnsignedShort(&port) &&
		h.dec.ReadVarInt(&nextState)

	if !ok {
		return h.dec.LastError()
	}

	switch s := protocol.State(nextState); s {
	// TODO support login
	case protocol.StateStatus:
		h.state = s
	default:
		return errInvalidState
	}

	return nil
}
