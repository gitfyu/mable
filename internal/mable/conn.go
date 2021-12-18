package mable

import (
	"errors"
	"github.com/gitfyu/mable/network"
	"github.com/gitfyu/mable/network/protocol"
	"github.com/gitfyu/mable/network/protocol/packet"
	"github.com/rs/zerolog/log"
	"net"
	"sync/atomic"
)

var errPacketHandlerPanic = errors.New("panic while handling packet")

// stateToPacketHandlers acts as a map with a protocol.State as key to a packetHandlerLookup value
var stateToPacketHandlers = []packetHandlerLookup{
	handshakeHandlers,
	statusHandlers,
	loginHandlers,
}

type connHandler struct {
	serv  *Server
	conn  net.Conn
	state protocol.State
	// closed acts as an atomic 'boolean' for Close and IsOpen
	closed int32
}

func newConnHandler(s *Server, c net.Conn) *connHandler {
	return &connHandler{
		serv:  s,
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
		MaxPacketSize: h.serv.cfg.MaxPacketSize,
	})

	for h.IsOpen() {
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

	return nil
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

// Close closes the connection, causing the client to be disconnected. This function may be called concurrently. Only
// the first call to it will actually close the connection, any further calls will simply be ignored.
func (h *connHandler) Close() error {
	// Documentation for net.Conn.Close doesn't seem to indicate whether it can safely be called multiple times, so
	// this will prevent duplicate calls just in case
	if !atomic.CompareAndSwapInt32(&h.closed, 0, 1) {
		return nil
	}

	return h.conn.Close()
}

// IsOpen returns whether the connection is still open
func (h *connHandler) IsOpen() bool {
	return atomic.LoadInt32(&h.closed) == 0
}

// WritePacket writes a single packet to the client. This function may be called concurrently.
func (h *connHandler) WritePacket(buf *network.PacketBuilder) error {
	_, err := h.conn.Write(buf.ToBytes())
	return err
}
