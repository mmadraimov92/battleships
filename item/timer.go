package item

import (
	"context"
	"fmt"
	"io"
	"time"

	"tui/terminal"
)

type timer struct {
	screen io.Writer
	input  chan terminal.KeyEvent
}

func (t *timer) Render(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case keyEvent := <-t.input:
				if keyEvent == terminal.DeleteKey || keyEvent == terminal.EscapeKey {
					cancel()
					return
				}
			}
		}
	}()

	var d time.Duration
	interval := 10 * time.Millisecond
	for {
		select {
		case <-ctx.Done():
			fmt.Fprintln(t.screen, "Timer stopped")
			time.Sleep(2 * time.Second)
			return
		case <-time.Tick(interval):
			terminal.ClearScreen(t.screen)
			fmt.Fprintln(t.screen, d)
			d += interval
		}
	}
}

func (t *timer) Title() {
	fmt.Fprint(t.screen, "Timer")
}
