package mable

import (
	"errors"
	"github.com/gitfyu/mable/protocol"
	"github.com/gitfyu/mable/protocol/packet"
)

func handleHandshake(c *connHandler) (protocol.State, error) {
	id, buf, err := c.readPacket()
	if err != nil {
		return 0, err
	}
	if id != packet.Handshake {
		return 0, errors.New("expected handshake")
	}

	ver, err := buf.ReadVarInt()
	if err != nil {
		return 0, err
	}
	c.version = protocol.Version(ver)

	// address
	if _, err := buf.ReadString(); err != nil {
		return 0, err
	}

	// port
	if _, err := buf.ReadUnsignedShort(); err != nil {
		return 0, err
	}

	nextState, err := buf.ReadVarInt()
	if err != nil {
		return 0, err
	}

	return protocol.State(nextState), nil
}
