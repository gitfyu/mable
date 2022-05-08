package server

import (
	"errors"
	"net"
	"sync/atomic"
	"time"

	"github.com/gitfyu/mable/chat"
	"github.com/gitfyu/mable/internal/protocol"
	"github.com/gitfyu/mable/internal/protocol/packet"
	"github.com/gitfyu/mable/internal/protocol/packet/outbound/login"
	"github.com/gitfyu/mable/internal/protocol/packet/outbound/play"
)

// conn represents a client connection.
type conn struct {
	serv       *Server
	conn       net.Conn
	state      protocol.State
	reader     *packet.Reader
	writer     *packet.Writer
	writeQueue chan packet.Outbound
	closed     int32
	flushed    chan struct{}
}

func newConn(s *Server, c net.Conn) *conn {
	return &conn{
		serv:  s,
		conn:  c,
		state: protocol.StateHandshake,
		reader: packet.NewReader(c, packet.ReaderConfig{
			MaxSize: s.cfg.MaxPacketSize,
		}),
		writer:     packet.NewWriter(c),
		writeQueue: make(chan packet.Outbound, 100), // TODO configurable size
		flushed:    make(chan struct{}),
	}
}

// dispatchPackets reads packets from conn.writeQueue and dispatches them until the connection is closed or an error
// occurs. When Close is called, this function will still dispatch packets that have been queued but not sent yet.
func (c *conn) dispatchPackets() {
	var err error
	for p := range c.writeQueue {
		if err = c.conn.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(c.serv.cfg.Timeout))); err != nil {
			break
		}

		c.writer.WritePacket(p)

		// TODO in the future, as an optimization flushing should probs not be done for every packet
		if err = c.writer.Flush(); err != nil {
			break
		}
	}

	close(c.flushed)

	if err != nil {
		c.serv.logger.Debug("Failed to dispatch packet(s)").Err(err).Log()
		c.Close()
	}
}

func (c *conn) handle() error {
	go c.dispatchPackets()

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
			c.Disconnect(&chat.Msg{Text: "Please use Minecraft 1.8."})
			return nil
		}

		username, id, err := handleLogin(c)
		if err != nil {
			return err
		}

		c.serv.logger.Info("Player logged in").
			Str("name", username).
			Stringer("id", id).
			Log()

		c.state = protocol.StatePlay
		defer func() {
			c.serv.logger.Info("Player disconnected").
				Str("name", username).
				Stringer("id", id).
				Log()
		}()
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

	close(c.writeQueue)
	<-c.flushed
	return c.conn.Close()
}

// IsOpen returns whether the connection is still open
func (c *conn) IsOpen() bool {
	return atomic.LoadInt32(&c.closed) == 0
}

// readPacket reads a single packet from the client.
func (c *conn) readPacket() (packet.Inbound, error) {
	if err := c.conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(c.serv.cfg.Timeout))); err != nil {
		return nil, err
	}

	return c.reader.ReadPacket(c.state)
}

// WritePacket writes a single packet to the client. This function may be called concurrently.
func (c *conn) WritePacket(pk packet.Outbound) {
	c.writeQueue <- pk
}

// Disconnect kicks the player with a specified reason.
func (c *conn) Disconnect(reason *chat.Msg) {
	c.serv.logger.Debug("Disconnecting").
		Stringer("reason", reason).
		Log()

	switch c.state {
	case protocol.StatePlay:
		c.WritePacket(&play.Disconnect{
			Reason: reason,
		})
	case protocol.StateLogin:
		c.WritePacket(&login.Disconnect{
			Reason: reason,
		})
	}

	c.Close()
}
