package battleships

import (
	"context"
	"errors"
	"io"
	"net"
)

type Simulator struct {
	server net.Conn
	client net.Conn
}

func NewSimulator() (*Simulator, error) {
	simulator := &Simulator{}

	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	go func() {
		defer ln.Close()
		server, err := ln.Accept()
		if err != nil {
			return
		}
		simulator.server = server
	}()

	client, err := net.Dial("tcp4", ln.Addr().String())
	if err != nil {
		return nil, err
	}
	simulator.client = client

	return simulator, nil
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

func (s *Simulator) Close() {
	s.server.Close()
	s.client.Close()
}
