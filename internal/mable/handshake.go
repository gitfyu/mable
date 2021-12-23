package mable

import (
	"errors"
	"github.com/gitfyu/mable/protocol"
	"github.com/gitfyu/mable/protocol/packet"
)

// handleHandshake processes the handshake packet
func handleHandshake(c *conn) (protocol.State, uint, error) {
	id, buf, err := c.readPacket()
	if err != nil {
		return 0, 0, err
	}
	if id != packet.Handshake {
		return 0, 0, errors.New("expected handshake")
	}

	ver, err := buf.ReadVarInt()
	if err != nil {
		return 0, 0, err
	}

	// address
	if _, err := buf.ReadString(); err != nil {
		return 0, 0, err
	}

	// port
	if _, err := buf.ReadUnsignedShort(); err != nil {
		return 0, 0, err
	}

	nextState, err := buf.ReadVarInt()
	if err != nil {
		return 0, 0, err
	}

	return protocol.State(nextState), uint(ver), nil
}
