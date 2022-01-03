package server

import (
	"github.com/gitfyu/mable/log"
	"net"
	"runtime/debug"
)

type Config struct {
	Addr          string
	MaxPacketSize int
	Timeout       int
	LogLevel      string
}

type Server struct {
	cfg      Config
	listener net.Listener
	logger   log.Logger
}

func NewServer(cfg Config) (*Server, error) {
	l, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}

	return &Server{
		cfg:      cfg,
		listener: l,
		logger: log.Logger{
			Name:     "SERVER",
			MinLevel: log.LevelFromString(cfg.LogLevel),
		},
	}, nil
}

func (s *Server) Addr() net.Addr {
	return s.listener.Addr()
}

// ListenAndServe will run the Server
func (s *Server) ListenAndServe() error {
	for {
		c, err := s.listener.Accept()
		if err != nil {
			return err
		}

		go s.handleConn(c)
	}
}

func (s *Server) handleConn(c net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
	}()

	s.logger.Debug("New connection").
		Stringer("src", c.RemoteAddr()).
		Log()

	h := newConn(s, c)
	defer h.Close()

	if err := h.handle(); err != nil {
		s.logger.Debug("Connection error").
			Err(err).
			Stringer("src", c.RemoteAddr()).
			Log()
	} else {
		s.logger.Debug("Connection closed").
			Stringer("src", c.RemoteAddr()).
			Log()
	}
}

// Close stops the server
func (s *Server) Close() error {
	return s.listener.Close()
}
