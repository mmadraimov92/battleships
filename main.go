package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"tui/app"
	"tui/terminal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	w := os.Stdout
	defer w.Close()

	inputChan := make(chan terminal.KeyEvent, 1)
	app := app.New(w, inputChan, cancel)

	go func() {
		app.Run(ctx)
	}()

	go func() {
		terminal.HandleKeyboardInput(ctx, inputChan)
	}()

	terminal.HideCursor(w)
	defer terminal.ShowCursor(w)

	<-ctx.Done()

	fmt.Fprint(w, "Exiting")
}
