package battleships

import (
	"context"
	"net"

	"tui/terminal"
)

type Battleships struct {
	input chan terminal.KeyEvent
	conn  net.Conn
}

func New(input chan terminal.KeyEvent) *Battleships {
	return &Battleships{
		input: input,
	}
}

func (*Battleships) Title() string {
	return "Battleships"
}

func (b *Battleships) Select(ctx context.Context) {
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
		case <-outgoingMessages:
			// todo: send message over network
			continue
		case keyEvent := <-b.input:
			if keyEvent == terminal.DeleteKey || keyEvent == terminal.EscapeKey {
				cancel()
				return
			}
			g.handleKeyEvent(keyEvent)
			draw(g)
		}
	}
}
