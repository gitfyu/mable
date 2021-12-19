package mable

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gitfyu/mable/chat"
	"github.com/gitfyu/mable/protocol"
	"github.com/gitfyu/mable/protocol/packet"
	"net"
	"sync"
	"sync/atomic"
)

var (
	errActionUnsupportedState = errors.New("action not supported in current state")
)

// stateToPacketHandlers acts as a map with a protocol.State as key to a packetHandlerLookup value
var stateToPacketHandlers = []packetHandlerLookup{
	handshakeHandlers,
	statusHandlers,
	loginHandlers,
}

type connHandler struct {
	serv      *Server
	conn      net.Conn
	state     protocol.State
	version   protocol.Version
	writer    *packet.Writer
	writeLock sync.Mutex
	// closed acts as an atomic 'boolean' for Close and IsOpen
	closed int32
}

func newConnHandler(s *Server, c net.Conn) *connHandler {
	return &connHandler{
		serv:   s,
		conn:   c,
		state:  protocol.StateHandshake,
		writer: packet.NewWriter(c),
	}
}

func (h *connHandler) handle() error {
	r := packet.NewReader(h.conn, packet.ReaderConfig{
		MaxSize: h.serv.cfg.MaxPacketSize,
	})

	for h.IsOpen() {
		id, buf, err := r.ReadPacket()
		if err != nil {
			return err
		}

		if err := h.handlePacket(id, buf); err != nil {
			return fmt.Errorf("failed to handle packet %d: %w", id, err)
		}
	}

	return nil
}

func (h *connHandler) validId(id packet.ID) bool {
	return int(id) < len(stateToPacketHandlers[h.state])
}

// handlePacket processes a packet. Packets that are not implemented will simply be ignored. This function will invoke
// packet.ReleaseBuffer on the provided Buffer before returning.
func (h *connHandler) handlePacket(id packet.ID, data *packet.Buffer) (err error) {
	defer packet.ReleaseBuffer(data)

	if !h.validId(id) {
		// Ignore unknown packets
		return
	}

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
func (h *connHandler) WritePacket(id packet.ID, buf *packet.Buffer) error {
	h.writeLock.Lock()
	defer h.writeLock.Unlock()

	return h.writer.WritePacket(id, buf)
}

// Disconnect kicks the player with a specified reason
func (h *connHandler) Disconnect(reason *chat.Msg) error {
	str, err := json.Marshal(reason)
	if err != nil {
		return err
	}

	switch h.state {
	// TODO impl this for 'play' state in the future as well
	case protocol.StateLogin:
		buf := packet.AcquireBuffer()
		defer packet.ReleaseBuffer(buf)

		buf.WriteStringFromBytes(str)

		if err := h.WritePacket(packet.LoginDisconnect, buf); err != nil {
			return err
		}

		return h.Close()
	default:
		return errActionUnsupportedState
	}
}
