package battleships

import (
	"context"
	"errors"
	"io"
	"net"
)

type Simulator struct {
	ln     net.Listener
	server net.Conn
}

func NewSimulator() (*Simulator, error) {
	s := &Simulator{}

	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	s.ln = ln

	go func() {
		defer ln.Close() //nolint:errcheck
		server, err := ln.Accept()
		if err != nil {
			return
		}
		s.server = server
	}()

	return s, nil
}

func (s *Simulator) Read(ctx context.Context) (*message, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			buf := make([]byte, 2)
			n, err := s.server.Read(buf[:cap(buf)])
			if err != nil {
				if errors.Is(err, io.EOF) {
					continue
				}
				return nil, err
			}
			buf = buf[:n]
			msg := decodeMessage(buf)
			return &msg, nil
		}
	}
}

func (s *Simulator) Write(msg message) error {
	_, err := s.server.Write(encodeMessage(msg))
	if err != nil {
		return err
	}

	return nil
}

func (s *Simulator) Addr() string {
	return s.ln.Addr().String()
}

func (s *Simulator) Close() {
	if s.server != nil {
		s.server.Close() //nolint:errcheck
	}
}
