package battleships

import (
	"log"
	"math/rand"
	"tui/terminal"
)

type mode int8

const (
	preparationMode mode = iota
	readyMode
	attackMode
	waitingMode
	winMode
	loseMode

	maxInitiative = 0b1110 // 14
)

func (m mode) String() string {
	switch m {
	case preparationMode:
		return "preparation"
	case readyMode:
		return "ready"
	case attackMode:
		return "attack"
	case waitingMode:
		return "waiting"
	case winMode:
		return "win"
	case loseMode:
		return "lose"
	default:
		return "unknown"
	}
}

type game struct {
	myBoard       *board
	targetBoard   *board
	mode          mode
	shipPlacement shipPlacement
	messages      chan<- message
	initiative    int8
}

func newGame(messages chan<- message) *game {
	return &game{
		myBoard:       newBoard(),
		targetBoard:   newBoard(),
		mode:          preparationMode,
		shipPlacement: newShipPlacement(),
		messages:      messages,
		initiative:    int8(rand.Intn(maxInitiative)),
	}
}

func (g *game) handleAttack(k terminal.KeyEvent) {
	if g.mode != attackMode {
		return
	}
	g.selectCellToAttack(k)
}

func (g *game) handlePreparationInput(k terminal.KeyEvent) {
	if g.mode != preparationMode {
		return
	}
	g.placeShips(k)
}

func (g *game) handleInitiativeMessage(c message) bool {
	defer log.Println("game initiative", g.initiative)
	if c.t != initiative {
		return false
	}

	if g.initiative > c.row {
		g.mode = attackMode
		log.Println("game mode set: attack")
		return true
	}

	if g.initiative < c.row {
		g.mode = waitingMode
		log.Println("game mode set: waiting")
		return true
	}

	return false
}

func (g *game) sendInitiative() {
	g.messages <- newInitiativeMessage(g.initiative)
}

func (g *game) handleIncomingMessage(c message) {
	defer log.Println("game mode:", g.mode)
	if g.mode != waitingMode {
		return
	}

	switch c.t {
	case attack:
		defer log.Printf("handled attack message")
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
		defer log.Printf("handled response message")
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
		return true
	}
	return false
}
