package battleships

import (
	"context"
	"io"

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

	g := newGame()
	i.draw(g)

	for {
		select {
		case <-ctx.Done():
			return
		case keyEvent := <-i.input:
			if keyEvent == terminal.DeleteKey || keyEvent == terminal.EscapeKey {
				cancel()
				return
			}
			g.handleKeyEvent(keyEvent)
			i.draw(g)
		}
	}
}
