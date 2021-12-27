package mable

import (
	"errors"
	"github.com/gitfyu/mable/chat"
	"github.com/gitfyu/mable/protocol"
	"github.com/gitfyu/mable/protocol/packet"
	"github.com/gitfyu/mable/protocol/packet/login"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	errActionUnsupportedState = errors.New("action not supported in current state")
)

type conn struct {
	serv      *Server
	conn      net.Conn
	state     protocol.State
	reader    *packet.Reader
	writer    *packet.Writer
	writeLock sync.Mutex
	closed    int32
}

func newConn(s *Server, c net.Conn) *conn {
	return &conn{
		serv:  s,
		conn:  c,
		state: protocol.StateHandshake,
		reader: packet.NewReader(c, packet.ReaderConfig{
			MaxSize: s.cfg.MaxPacketSize,
		}),
		writer: packet.NewWriter(c),
	}
}

func (c *conn) handle() error {
	s, ver, err := handleHandshake(c)
	if err != nil {
		return err
	}

	c.state = s
	switch s {
	case protocol.StateStatus:
		return handleStatus(c)
	case protocol.StateLogin:
		if ver != 47 {
			return c.Disconnect(&chat.Msg{Text: "Please use Minecraft 1.8."})
		}

		username, id, err := handleLogin(c)
		if err != nil {
			return err
		}

		c.state = protocol.StatePlay
		return handlePlay(c, username, id)
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

	return c.conn.Close()
}

// IsOpen returns whether the connection is still open
func (c *conn) IsOpen() bool {
	return atomic.LoadInt32(&c.closed) == 0
}

func (c *conn) readPacket() (packet.Inbound, error) {
	if err := c.conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(c.serv.cfg.Timeout))); err != nil {
		return nil, err
	}

	return c.reader.ReadPacket(c.state)
}

// WritePacket writes a single packet to the client. This function may be called concurrently.
func (c *conn) WritePacket(pk packet.Outbound) error {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()

	if err := c.conn.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(c.serv.cfg.Timeout))); err != nil {
		return err
	}

	return c.writer.WritePacket(pk)
}

// Disconnect kicks the player with a specified reason
func (c *conn) Disconnect(reason *chat.Msg) error {
	var err error
	switch c.state {
	case protocol.StatePlay:
		err = c.WritePacket(&login.Disconnect{
			Reason: reason,
		})
	case protocol.StateLogin:
		err = c.WritePacket(&login.Disconnect{
			Reason: reason,
		})
	default:
		return errActionUnsupportedState
	}

	if err != nil {
		return err
	}

	return c.Close()
}
