package battleships

import (
	"context"
	"net"
	"time"
	"tui/terminal"
)

func (b *Battleships) connect(ctx context.Context) net.Conn {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var conn net.Conn

	if b.isServer {
		listener, err := net.Listen("tcp4", b.address)
		if err != nil {
			b.logger.Error(err.Error())
			cancel()
		}
		defer listener.Close()

		terminal.ClearScreen()
		terminal.Draw("Waiting for other player to connect")
		b.logger.Info("Waiting for other player to connect")
		ready := make(chan struct{})
		go func() {
			conn, err = listener.Accept()
			if err != nil {
				b.logger.Error(err.Error())
				cancel()
			}
			close(ready)
		}()

		select {
		case <-ready:
		case <-ctx.Done():
		}
		return conn
	}

	for {
		select {
		case <-ctx.Done():
			return conn
		case <-time.Tick(time.Second):
			var d net.Dialer
			var err error
			d.Timeout = 5 * time.Second
			conn, err = d.DialContext(ctx, "tcp4", b.address)
			if err != nil {
				terminal.ClearScreen()
				terminal.Draw("Waiting for the server connection")
				terminal.CursorNextLine()
				b.logger.Error(err.Error())
				continue
			}
			return conn
		}
	}
}
