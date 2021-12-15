package mable

import (
	"bufio"
	"github.com/gitfyu/mable/network"
	"github.com/gitfyu/mable/network/protocol"
	"net"
)

// TODO the size param might not be needed
type packetHandler func(size int, h *connHandler) error

// idToPacketHandler acts as a map with a packet ID as key and a packetHandler as value
type idToPacketHandler []packetHandler

// stateToPacketHandlers acts as a map with a protocol.State as key to a idToPacketHandler value
var stateToPacketHandlers = [][]packetHandler{
	handshakeHandlers,
	statusHandlers,
}

type connHandler struct {
	conn  net.Conn
	dec   *network.PacketDecoder
	enc   *network.PacketEncoder
	state protocol.State
}

func newConnHandler(c net.Conn) *connHandler {
	return &connHandler{
		conn:  c,
		dec:   network.NewPacketDecoder(bufio.NewReader(c)),
		enc:   network.NewPacketEncoder(bufio.NewWriter(c)),
		state: protocol.StateHandshake,
	}
}

func (h *connHandler) handle() error {
	var size, id network.VarInt
	var idSize int

	for {
		if ok := h.dec.ReadVarInt(&size) && h.dec.ReadVarIntAndSize(&id, &idSize); !ok {
			return h.dec.LastError()
		}

		bodySize := int(size) - idSize

		if !h.validId(id) {
			// Since a lot of packets will probably never be implemented, unknown packets are simply ignored
			if !h.dec.Skip(bodySize) {
				return h.dec.LastError()
			}
			continue
		}

		if err := h.handlePacket(id, bodySize); err != nil {
			return err
		}
	}
}

func (h *connHandler) validId(id network.VarInt) bool {
	return id >= 0 && int(id) < len(stateToPacketHandlers[h.state])
}

func (h *connHandler) handlePacket(id network.VarInt, bodySize int) error {
	return stateToPacketHandlers[h.state][id](bodySize, h)
}

func (h *connHandler) Close() error {
	return h.conn.Close()
}
