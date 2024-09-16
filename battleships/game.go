package battleships

import "tui/terminal"

type mode int

const (
	preparation mode = iota
	ready
	attack
	waiting
)

type game struct {
	myBoard       *board
	targetBoard   *board
	mode          mode
	shipPlacement shipPlacement
}

func newGame() *game {
	return &game{
		myBoard:       newBoard(),
		targetBoard:   newBoard(),
		mode:          preparation,
		shipPlacement: newShipPlacement(),
	}
}

func (g *game) handleKeyEvent(k terminal.KeyEvent) {
	switch g.mode {
	case waiting, ready:
		return
	case preparation:
		g.placeShips(k)
	case attack:
		g.targetBoard.selectCellToAttack(k)
	}
}

func (g *game) shipPlacementInfo() string {
	switch g.shipPlacement.currentlyPlacing {
	case carrier:
		return "Carrier (5 cells)"
	case battleship:
		return "Battleship (4 cells)"
	case submarine:
		return "Submarine (3 cells)"
	case cruiser:
		return "Cruiser (3 cells)"
	case destroyer:
		return "Destroyer (2 cells)"
	default:
		return ""
	}
}

func (g *game) placeShips(k terminal.KeyEvent) {
	switch k {
	case terminal.UpArrowKey:
		g.myBoard.selectedRow.Decrement()
	case terminal.DownArrowKey:
		g.myBoard.selectedRow.Increment()
	case terminal.RightArrowKey:
		g.myBoard.selectedCol.Increment()
	case terminal.LeftArrowKey:
		g.myBoard.selectedCol.Decrement()
	case terminal.SmallRKey:
		g.shipPlacement.orientation.Increment()
	case terminal.EnterKey:
		g.shipPlacement.acceptPlacement()
		g.myBoard.selectedRow.Reset()
		g.myBoard.selectedCol.Reset()
		g.shipPlacement.orientation.Reset()
		g.checkIfAllPlaced()
		if g.mode == ready {
			return
		}
	}

	if g.shipPlacement.isValidPlacement(g.myBoard) {
		g.shipPlacement.placeOnBoard(g.myBoard)
	} else { // revert back
		switch k {
		case terminal.UpArrowKey:
			g.myBoard.selectedRow.Increment()
		case terminal.DownArrowKey:
			g.myBoard.selectedRow.Decrement()
		case terminal.RightArrowKey:
			g.myBoard.selectedCol.Decrement()
		case terminal.LeftArrowKey:
			g.myBoard.selectedCol.Increment()
		case terminal.SmallRKey:
			g.shipPlacement.orientation.Decrement()
		}
	}

}

func (g *game) checkIfAllPlaced() {
	if len(g.shipPlacement.placed) == 5 {
		g.mode = ready
	}
}
