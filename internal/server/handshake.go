package server

import (
	"errors"

	"github.com/gitfyu/mable/internal/protocol"
	"github.com/gitfyu/mable/internal/protocol/packet/inbound/handshake"
)

// handleHandshake processes the handshake sequence and returns the next protocol.State and the client's protocol
// version.
func handleHandshake(c *conn) (protocol.State, int32, error) {
	pk, err := c.readPacket()
	if err != nil {
		return 0, 0, err
	}
	h, ok := pk.(*handshake.Handshake)
	if !ok {
		return 0, 0, errors.New("expected handshake")
	}

	return protocol.State(h.NextState), h.ProtoVer, nil
}
