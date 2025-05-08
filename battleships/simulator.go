package battleships

import (
	"context"
	"errors"
	"io"
)

type Simulator struct {
	conn io.ReadWriter
}

func NewSimulator() (*Simulator, io.ReadWriter) {
	clientR, serverW := io.Pipe()
	serverR, clientW := io.Pipe()

	clientConn := struct {
		io.Reader
		io.Writer
	}{
		Reader: clientR,
		Writer: clientW,
	}

	simulator := &Simulator{
		conn: struct {
			io.Reader
			io.Writer
		}{
			Reader: serverR,
			Writer: serverW,
		},
	}

	return simulator, clientConn
}

func (s *Simulator) Read(ctx context.Context) (*message, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			buf := make([]byte, 2)
			n, err := s.conn.Read(buf[:cap(buf)])
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
	_, err := s.conn.Write(encodeMessage(msg))
	if err != nil {
		return err
	}

	return nil
}
