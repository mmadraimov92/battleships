package battleships

import "tui/terminal"

type mode string

const (
	PlacingShips mode = "Place your ships"
	Ready        mode = "Ready"
	Attack       mode = "Attack"
	Waiting      mode = "Waiting"
)

type game struct {
	myBoard     *board
	targetBoard *board
	mode        mode
}

func newGame() *game {
	return &game{
		myBoard:     newBoard(),
		targetBoard: newBoard(),
		mode:        PlacingShips,
	}
}

func (g *game) handleKeyEvent(k terminal.KeyEvent) {
	switch g.mode {
	case Waiting, Ready:
		return
	case PlacingShips:
		g.myBoard.handleKeyEvent(k)
	case Attack:
		g.targetBoard.handleKeyEvent(k)
	}
}
