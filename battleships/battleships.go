package battleships

import (
	"context"
	"fmt"
	"io"
	"time"

	"tui/terminal"
)

type Item struct {
	screen io.Writer
	input  chan terminal.KeyEvent
}

func New(w io.Writer, input chan terminal.KeyEvent) *Item {
	return &Item{
		screen: w,
		input:  input,
	}
}

func (*Item) Title() string {
	return "Battleships"
}

func (i *Item) Render(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case keyEvent := <-i.input:
				if keyEvent == terminal.DeleteKey || keyEvent == terminal.EscapeKey {
					cancel()
					return
				}
			}
		}
	}()

	terminal.ClearScreen(i.screen)
	fmt.Fprintln(i.screen, "Battleships started")

	for {
		select {
		case <-ctx.Done():
			terminal.ClearScreen(i.screen)
			fmt.Fprintln(i.screen, "Battleships stopped")
			time.Sleep(time.Second)
			return
		}
	}
}
