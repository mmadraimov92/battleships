package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"tui/item"
	"tui/terminal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	w := os.Stdout
	defer w.Close()

	inputChan := make(chan terminal.KeyEvent, 1)
	mainMenu := item.NewMainMenu(w, inputChan, cancel)

	go func() {
		mainMenu.Render(ctx)
	}()

	go func() {
		terminal.HandleKeyboardInput(ctx, inputChan)
	}()

	terminal.HideCursor(w)
	defer terminal.ShowCursor(w)

	<-ctx.Done()

	fmt.Fprint(w, "Exiting")
}
