package battleships

import "tui/terminal"

type mode int

const (
	preparationMode mode = iota
	readyMode
	attackMode
	waitingMode
	winMode
	loseMode
)

type game struct {
	myBoard       *board
	targetBoard   *board
	mode          mode
	shipPlacement shipPlacement
	messages      chan<- message
}

func newGame(messages chan<- message) *game {
	return &game{
		myBoard:       newBoard(),
		targetBoard:   newBoard(),
		mode:          preparationMode,
		shipPlacement: newShipPlacement(),
		messages:      messages,
	}
}

func (g *game) handleKeyEvent(k terminal.KeyEvent) {
	switch g.mode {
	case waitingMode, readyMode, winMode, loseMode:
		return
	case preparationMode:
		g.placeShips(k)
	case attackMode:
		g.selectCellToAttack(k)
	}
}

func (g *game) handleIncomingMessage(c message) {
	if g.mode != waitingMode {
		return
	}

	switch c.t {
	case attack:
		cell := g.myBoard.cellAt(c.row, c.col)
		if cell.shipClass == empty {
			g.messages <- newResponseMessageMiss(c.row, c.col)
			g.mode = attackMode
			return
		}

		g.myBoard.markAsHit(c.row, c.col, cell.shipClass)
		gameOver := g.myBoard.isDestroyed()
		g.messages <- newResponseMessageHit(c.row, c.col, cell.shipClass, gameOver)
		if gameOver {
			g.mode = loseMode
		} else {
			g.mode = attackMode
		}

	case response:
		if c.status == statusHit {
			g.targetBoard.markAsHit(c.row, c.col, c.ship)
			if c.gameOver {
				g.mode = winMode
			}
			return
		}

		g.targetBoard.markAsMiss(c.row, c.col)
		g.mode = waitingMode
	}
}

func (g *game) selectCellToAttack(k terminal.KeyEvent) {
	switch k {
	case terminal.UpArrowKey:
		g.targetBoard.selectedRow.Decrement()
	case terminal.DownArrowKey:
		g.targetBoard.selectedRow.Increment()
	case terminal.RightArrowKey:
		g.targetBoard.selectedCol.Increment()
	case terminal.LeftArrowKey:
		g.targetBoard.selectedCol.Decrement()
	case terminal.EnterKey:
		g.messages <- newAttackMessage(g.targetBoard.selectedRow.Current(), g.targetBoard.selectedCol.Current())
		g.mode = waitingMode
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
		if g.areAllShipsPlaced() {
			return
		}
	}

	g.shipPlacement.placeInValidPosition(g.myBoard)
}

func (g *game) areAllShipsPlaced() bool {
	if len(g.shipPlacement.placed) == 5 {
		g.mode = readyMode
		return true
	}
	return false
}
