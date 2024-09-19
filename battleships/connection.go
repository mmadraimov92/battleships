package battleships

import (
	"context"
	"io"
	"net"

	"tui/network"
)

type connection struct {
	listenFrom net.Addr
	connectTo  net.Addr
}

// todo: get IPs from flags
func newConnection() *connection {
	return &connection{
		listenFrom: &net.TCPAddr{
			IP:   net.IPv4(0, 0, 0, 0),
			Port: 1931,
		},
		connectTo: &net.TCPAddr{
			IP:   net.IPv4(127, 0, 0, 1),
			Port: 1931,
		},
	}
}

func (c *connection) handle(ctx context.Context) error {
	srv, err := network.NewServer(c.listenFrom.String())
	if err != nil {
		return err
	}

	go srv.Start(ctx, noopHandleRead)

	return nil
}

func noopHandleRead(io.Reader) {}
