package app

import (
	"context"
	"fmt"
	"time"

	"tui/terminal"
)

type timer struct {
	input chan terminal.KeyEvent
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
	interval := 50 * time.Millisecond
	for {
		select {
		case <-ctx.Done():
			t.draw(fmt.Sprintf("Timer stopped at %s", d.String()))
			time.Sleep(2 * time.Second)
			return
		case <-time.Tick(interval):
			t.draw(d.String())
			d += interval
		}
	}
}

func (t *timer) Title() string {
	return "Timer"
}

func (t *timer) draw(s string) {
	terminal.ClearScreen()
	terminal.Draw(s)
}
