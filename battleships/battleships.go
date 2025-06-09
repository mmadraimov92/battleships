package battleships

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"tui/terminal"
)

type Battleships struct {
	input  chan terminal.KeyEvent
	conn   io.ReadWriter
	logger *slog.Logger
	g      *game
}

func New(input chan terminal.KeyEvent, conn io.ReadWriter, logger *slog.Logger) *Battleships {
	return &Battleships{
		input:  input,
		conn:   conn,
		logger: logger,
	}
}

func (*Battleships) Title() string {
	return "Battleships PVP"
}

func (b *Battleships) Select(ctx context.Context) {
	b.start(ctx)
}

func (b *Battleships) start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	incomingMessages := make(chan message, 1)
	outgoingMessages := make(chan message, 1)

	b.g = newGame(outgoingMessages, b.logger)

	// todo: think of something else
	// initial board placement drawing
	b.g.shipPlacement.placeInValidPosition(b.g.myBoard)

	draw(b.g)

preparation:
	for {
		select {
		case <-ctx.Done():
			return
		case keyEvent := <-b.input:
			if keyEvent == terminal.DeleteKey || keyEvent == terminal.EscapeKey {
				cancel()
				return
			}
			b.g.handlePreparationInput(keyEvent)
			if b.g.areAllShipsPlaced() {
				b.g.mode = readyMode
				break preparation
			}
			draw(b.g)
		}
	}

	b.g.sendInitiative()
	go func(ctx context.Context) {
		buf := make([]byte, 8)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				n, err := b.conn.Read(buf[:cap(buf)])
				if err != nil && !errors.Is(err, io.EOF) {
					b.logger.Error(err.Error())
					terminal.Draw("Connection error...Exiting")
					time.Sleep(2 * time.Second)
					cancel()
				}
				if n == 0 {
					continue
				}
				buf = buf[:n]
				b.logger.Debug(fmt.Sprintf("received message: %+v", decodeMessage(buf)))
				incomingMessages <- decodeMessage(buf)
			}
		}
	}(ctx)

	draw(b.g)
	b.logger.Debug("Start initiative")
initiative:
	for {
		select {
		case <-ctx.Done():
			return
		case c := <-incomingMessages:
			done := b.g.handleInitiativeMessage(c)
			b.logger.Debug("handled initiative message")
			if done {
				break initiative
			}
			b.g.sendInitiative()
		case m := <-outgoingMessages:
			_, err := b.conn.Write(encodeMessage(m))
			if err != nil {
				b.logger.Error(err.Error())
				terminal.Draw("Connection error...Exiting")
				time.Sleep(2 * time.Second)
				cancel()
			}
		}
	}

	b.logger.Debug("Start main game loop")
	draw(b.g)
	for {
		select {
		case <-ctx.Done():
			return
		case c := <-incomingMessages:
			b.g.handleIncomingMessage(c)
		case m := <-outgoingMessages:
			_, err := b.conn.Write(encodeMessage(m))
			if err != nil {
				b.logger.Error(err.Error())
				terminal.Draw("Connection error...Exiting")
				time.Sleep(2 * time.Second)
				cancel()
			}
		case keyEvent := <-b.input:
			if keyEvent == terminal.DeleteKey || keyEvent == terminal.EscapeKey {
				cancel()
				return
			}
			b.g.handleAttack(keyEvent)
		}
		draw(b.g)
	}
}
