package main

import (
	"context"
	"fmt"
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

	terminal.HideCursor()
	defer terminal.ShowCursor()

	inputChan := make(chan terminal.KeyEvent, 1)
	defer close(inputChan)

	// todo: setup logging to file

	go func() {
		err := terminal.HandleKeyboardInput(ctx, inputChan)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			cancel()
		}
	}()

	go func() {
		menu.New(
			inputChan,
			[]menu.Item{
				// todo: initialize conn
				battleships.New(inputChan, nil),
				menu.NewExit(cancel),
			},
		).Run(ctx)
	}()

	<-ctx.Done()

	fmt.Fprintln(os.Stdout, "Exiting")
}
