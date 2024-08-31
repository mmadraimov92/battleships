package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"tui/app"
	"tui/terminal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	terminal.HideCursor()
	defer terminal.ShowCursor()

	inputChan := make(chan terminal.KeyEvent, 1)
	defer close(inputChan)

	app := app.New(inputChan, cancel)

	go func() {
		err := terminal.HandleKeyboardInput(ctx, inputChan)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			cancel()
		}
	}()

	go func() {
		app.Run(ctx)
	}()

	<-ctx.Done()

	fmt.Fprintln(os.Stdout, "Exiting")
}
