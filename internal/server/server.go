package server

import (
	"github.com/rs/zerolog/log"
	"net"
	"runtime/debug"
)

type Config struct {
	Addr          string
	MaxPacketSize int
	Timeout       int
}

type Server struct {
	cfg      Config
	listener net.Listener
}

func NewServer(cfg Config) (*Server, error) {
	l, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}

	return &Server{
		cfg:      cfg,
		listener: l,
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

	h := newConn(s, c)
	defer h.Close()

	if err := h.handle(); err != nil {
		log.Debug().
			Err(err).
			Str("src", c.RemoteAddr().String()).
			Msg("Connection error")
	}
}

// Close stops the server
func (s *Server) Close() error {
	return s.listener.Close()
}
