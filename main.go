package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tui/battleships"
	"tui/menu"
	"tui/terminal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	isServer := flag.Bool("server", false, "start instance as server")
	addr := flag.String("addr", "127.0.0.1:1337", "address of server instance. Default 127.0.0.1:1337")

	flag.Parse()

	terminal.SetRendererOutput(os.Stdout)

	terminal.HideCursor()
	defer terminal.ShowCursor()

	inputChan := make(chan terminal.KeyEvent, 1)
	defer close(inputChan)

	logInstance := "client"
	if *isServer {
		logInstance = "server"
	}
	logFile, err := os.Create(fmt.Sprintf("log_%s.txt", logInstance))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		cancel()
	}
	defer logFile.Close()

	logger := slog.New(slog.NewTextHandler(logFile, nil).WithAttrs([]slog.Attr{slog.String("instance", logInstance)}))

	go func() {
		err := terminal.HandleKeyboardInput(ctx, inputChan)
		if err != nil {
			logger.Error(err.Error())
			cancel()
		}
	}()

	var conn net.Conn
	if *isServer {
		listener, err := net.Listen("tcp4", *addr)
		if err != nil {
			logger.Error(err.Error())
			cancel()
		}
		defer listener.Close()

		fmt.Fprintln(os.Stdout, "Waiting for other player to connect")
		logger.Info("Waiting for other player to connect")
		ready := make(chan struct{})
		go func() {
			conn, err = listener.Accept()
			if err != nil {
				logger.Error(err.Error())
				cancel()
			}
			close(ready)
		}()

		select {
		case <-ready:
			defer conn.Close()
		case <-ctx.Done():
		}
	} else {
		var d net.Dialer
		d.Timeout = 5 * time.Second
		conn, err = d.DialContext(ctx, "tcp4", *addr)
		if err != nil {
			logger.Error(err.Error())
			cancel()
		}
		if conn != nil {
		defer conn.Close()
		}
	}

	go func() {
		menu.New(
			inputChan,
			[]menu.Item{
				battleships.New(inputChan, conn, logger),
				menu.NewExit(cancel),
			},
		).Run(ctx)
	}()

	<-ctx.Done()

	logger.Info("Exiting")
}
