package mable

import (
	"io"
	"net"
)

type Server struct {
	closer io.Closer
}

// ListenAndServe will run the Server using the provided Config
func (s *Server) ListenAndServe(cfg *Config) error {
	l, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}

	defer l.Close()
	s.closer = l

	for {
		c, err := l.Accept()
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
	return s.closer.Close()
}
