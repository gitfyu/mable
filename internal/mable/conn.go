package mable

import (
	"encoding/json"
	"errors"
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

type conn struct {
	serv      *Server
	conn      net.Conn
	state     protocol.State
	version   protocol.Version
	readBuf   *packet.Buffer
	reader    *packet.Reader
	writer    *packet.Writer
	writeLock sync.Mutex
	closed    int32
}

func newConn(s *Server, c net.Conn) *conn {
	return &conn{
		serv:    s,
		conn:    c,
		state:   protocol.StateHandshake,
		readBuf: packet.AcquireBuffer(),
		reader: packet.NewReader(c, packet.ReaderConfig{
			MaxSize: s.cfg.MaxPacketSize,
		}),
		writer: packet.NewWriter(c),
	}
}

func (c *conn) handle() error {
	s, err := handleHandshake(c)
	if err != nil {
		return err
	}

	switch s {
	case protocol.StateStatus:
		return handleStatus(c)
	case protocol.StateLogin:
		c.state = s
		if err := handleLogin(c); err != nil {
			return err
		}

		c.state = protocol.StatePlay
		return handlePlay(c)
	default:
		return errors.New("unknown state")
	}
}

// Close closes the connection, causing the client to be disconnected. This function may be called concurrently. Only
// the first call to it will actually close the connection, any further calls will simply be ignored.
func (c *conn) Close() error {
	// Documentation for net.Conn.Close doesn't seem to indicate whether it can safely be called multiple times, so
	// this will prevent duplicate calls just in case
	if !atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		return nil
	}

	packet.ReleaseBuffer(c.readBuf)
	return c.conn.Close()
}

// IsOpen returns whether the connection is still open
func (c *conn) IsOpen() bool {
	return atomic.LoadInt32(&c.closed) == 0
}

func (c *conn) readPacket() (packet.ID, *packet.Buffer, error) {
	id, err := c.reader.ReadPacket(c.readBuf)
	return id, c.readBuf, err
}

// WritePacket writes a single packet to the client. This function may be called concurrently.
func (c *conn) WritePacket(id packet.ID, buf *packet.Buffer) error {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()

	return c.writer.WritePacket(id, buf)
}

// Disconnect kicks the player with a specified reason
func (c *conn) Disconnect(reason *chat.Msg) error {
	str, err := json.Marshal(reason)
	if err != nil {
		return err
	}

	var id packet.ID

	switch c.state {
	case protocol.StatePlay:
		id = packet.PlayDisconnect
	case protocol.StateLogin:
		id = packet.LoginDisconnect
	default:
		return errActionUnsupportedState
	}

	buf := packet.AcquireBuffer()
	defer packet.ReleaseBuffer(buf)

	buf.WriteStringFromBytes(str)

	if err := c.WritePacket(id, buf); err != nil {
		return err
	}

	return c.Close()
}
