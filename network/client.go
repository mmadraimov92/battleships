package network

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

func Connect(ctx context.Context, addr string) (io.WriteCloser, error) {
	var d net.Dialer
	d.Timeout = 5 * time.Second
	conn, err := d.DialContext(ctx, "tcp4", addr)
	if err != nil {
		return nil, fmt.Errorf("DialContext: %w", err)
	}

	return conn, nil
}
