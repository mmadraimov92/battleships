package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tui/app"
	"tui/terminal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	w := os.Stdout
	defer w.Close()

	terminal.HideCursor(w)
	defer terminal.ShowCursor(w)

	inputChan := make(chan terminal.KeyEvent, 1)
	app := app.New(w, inputChan, cancel)

	go func() {
		err := terminal.HandleKeyboardInput(ctx, inputChan)
		if err != nil {
			fmt.Fprintln(w, err.Error())
			cancel()
		}
	}()

	go func() {
		app.Run(ctx)
	}()

	<-ctx.Done()

	fmt.Fprintln(w, "Exiting")
	time.Sleep(100 * time.Millisecond) // give some time for all stdout to print, todo: remove later
}
