package mable

import (
	"errors"
	"github.com/gitfyu/mable/network"
	"github.com/gitfyu/mable/network/protocol"
	"github.com/gitfyu/mable/network/protocol/packet"
	"github.com/rs/zerolog/log"
	"net"
)

var errPacketHandlerPanic = errors.New("panic while handling packet")

type packetHandler func(h *connHandler, data *network.PacketData) error

// idToPacketHandler acts as a map with a packet ID as key and a packetHandler as value
type idToPacketHandler []packetHandler

// stateToPacketHandlers acts as a map with a protocol.State as key to a idToPacketHandler value
var stateToPacketHandlers = [][]packetHandler{
	handshakeHandlers,
	statusHandlers,
}

type connHandler struct {
	conn  net.Conn
	state protocol.State
}

func newConnHandler(c net.Conn) *connHandler {
	return &connHandler{
		conn:  c,
		state: protocol.StateHandshake,
	}
}

func (h *connHandler) handle() error {
	var id packet.ID
	var buf []byte
	var err error
	var data network.PacketData

	r := network.NewReader(h.conn, network.ReaderConfig{
		// TODO currently this is just an arbitrarily chosen limit
		MaxPacketSize: 1 << 16,
	})

	for {
		id, buf, err = r.ReadPacket(buf)
		if err != nil {
			return err
		}

		if !h.validId(id) {
			// Ignore unknown packets
			continue
		}

		data.Load(buf)
		if err := h.handlePacket(id, &data); err != nil {
			return err
		}
	}
}

func (h *connHandler) validId(id packet.ID) bool {
	return id >= 0 && int(id) < len(stateToPacketHandlers[h.state])
}

func (h *connHandler) handlePacket(id packet.ID, data *network.PacketData) (err error) {
	defer func() {
		if r := recover(); r != nil {
			e := log.Debug().
				Int("id", int(id)).
				Int("state", int(h.state))

			if err, ok := r.(error); ok {
				e.Err(err)
			}

			e.Msg("Error handling packet")
			err = errPacketHandlerPanic
		}
	}()
	return stateToPacketHandlers[h.state][id](h, data)
}

// Close closes the connection, causing the client to be disconnected
func (h *connHandler) Close() error {
	return h.conn.Close()
}

// WritePacket writes a single packet to the client. This function may be called concurrently.
func (h *connHandler) WritePacket(buf *network.PacketBuilder) error {
	_, err := h.conn.Write(buf.ToBytes())
	return err
}
