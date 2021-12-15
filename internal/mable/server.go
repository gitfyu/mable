package mable

import (
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

// ListenAndServe will run the Server using the provided Config
func (s *Server) ListenAndServe(cfg *Config) error {
	for {
		c, err := s.listener.Accept()
		if err != nil {
			return err
		}

		go s.handleConn(c)
	}
}

func (s *Server) handleConn(c net.Conn) {
	defer c.Close()

	// TODO
	c.Write([]byte("test"))
}

// Close stops the server
func (s *Server) Close() error {
	return s.listener.Close()
}
