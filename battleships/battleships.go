package battleships

import (
	"context"

	"tui/terminal"
)

type Item struct {
	input chan terminal.KeyEvent
}

func New(input chan terminal.KeyEvent) *Item {
	return &Item{
		input: input,
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
