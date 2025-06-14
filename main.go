package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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

	logger := slog.New(
		slog.
			NewTextHandler(logFile, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}).
			WithAttrs([]slog.Attr{slog.String("instance", logInstance)}),
	)

	go func() {
		err := terminal.HandleKeyboardInput(ctx, inputChan)
		if err != nil {
			logger.Error(err.Error())
			cancel()
		}
	}()

	options := []func(*battleships.Battleships){
		battleships.WithAddress(*addr),
	}
	if *isServer {
		options = append(options, battleships.AsServer())
	}

	go func() {
		menu.New(
			inputChan,
			[]menu.Item{
				battleships.New(inputChan, logger, options...),
				menu.NewBattleshipsAI(),
				menu.NewExit(cancel),
			},
		).Run(ctx)
	}()

	<-ctx.Done()

	logger.Info("Exiting")
}
