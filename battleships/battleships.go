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

	// todo: think of something else
	// initial board placement drawing
	g.shipPlacement.placeInValidPosition(g.myBoard)

	draw(g)

	for {
		select {
		case <-ctx.Done():
			return
		case c := <-g.incomingMessages:
			g.handleIncomingMessage(c)
		case keyEvent := <-i.input:
			if keyEvent == terminal.DeleteKey || keyEvent == terminal.EscapeKey {
				cancel()
				return
			}
			g.handleKeyEvent(keyEvent)
			draw(g)
		}
	}
}
