package mable

import (
	"github.com/rs/zerolog/log"
	"net"
)

type Server struct {
	cfg      *Config
	listener net.Listener
}

func NewServer(cfg *Config) (*Server, error) {
	l, err := net.Listen("tcp", cfg.Address)
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
	h := newConnHandler(c)
	defer h.Close()

	if err := h.handle(); err != nil {
		log.Debug().
			Err(err).
			Stack().
			Str("src", c.RemoteAddr().String()).
			Msg("Connection error")
	}
}

// Close stops the server
func (s *Server) Close() error {
	return s.listener.Close()
}
