package network

import (
	"context"
	"fmt"
	"io"
	"net"
)

type Server struct {
	l net.Listener
}

func NewServer(addr string) (*Server, error) {
	listener, err := net.Listen("tcp4", addr)
	if err != nil {
		return nil, fmt.Errorf("net.Listen: %w", err)
	}

	return &Server{listener}, nil
}

func (s *Server) Start(ctx context.Context, handler func(io.Reader)) error {
	conn, err := s.l.Accept()
	if err != nil {
		return fmt.Errorf("listener.Accept: %w", err)
	}
	conn.Close()

	go handler(conn)

	<-ctx.Done()

	return nil
}
