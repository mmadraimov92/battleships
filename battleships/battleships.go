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

	incomingMessages := make(chan message, 1)
	outgoingMessages := make(chan message, 1)

	g := newGame(outgoingMessages)

	// todo: think of something else
	// initial board placement drawing
	g.shipPlacement.placeInValidPosition(g.myBoard)

	draw(g)

	// todo: start goroutine which listens to outgoing messages from game obj
	// and relays to underlying net.Conn.
	// also listens to net.Conn for incoming messages and relays to incomingMessages chan

	for {
		select {
		case <-ctx.Done():
			return
		case c := <-incomingMessages:
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
