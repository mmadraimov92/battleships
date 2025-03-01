package battleships

import (
	"testing"

	"tui/terminal"
)

func TestGame_PlayerLoses(t *testing.T) {
	messages := make(chan message, 1)
	game := newGame(messages)

	myBoard := newBoard()
	myBoard.cells[0][0] = cell{shipClass: destroyer, status: statusUndefined}
	myBoard.cells[1][0] = cell{shipClass: destroyer, status: statusUndefined}

	game.myBoard = myBoard
	game.mode = waitingMode

	// opponent attacks first
	game.handleIncomingMessage(newAttackMessage(0, 0))

	got := <-messages

	assertEqual(t, got.t, response)
	assertEqual(t, got.status, statusHit)
	assertEqual(t, got.ship, destroyer)
	assertEqual(t, got.gameOver, false)
	assertEqual(t, game.mode, attackMode)

	// player attacks 0,0
	game.selectCellToAttack(terminal.EnterKey)

	got = <-messages

	assertEqual(t, got.t, attack)
	assertEqual(t, got.row, 0)
	assertEqual(t, got.col, 0)
	assertEqual(t, game.mode, waitingMode)

	// opponent response with hit
	game.handleIncomingMessage(newResponseMessageHit(0, 0, destroyer, false))

	assertEqual(t, game.targetBoard.cellAt(0, 0).shipClass, destroyer)
	assertEqual(t, game.targetBoard.cellAt(0, 0).status, statusHit)
	assertEqual(t, game.mode, waitingMode)

	// opponent attacks 1,0
	game.handleIncomingMessage(newAttackMessage(1, 0))

	got = <-messages

	assertEqual(t, got.t, response)
	assertEqual(t, got.status, statusHit)
	assertEqual(t, got.ship, destroyer)
	assertEqual(t, got.gameOver, true)
	assertEqual(t, game.mode, loseMode)
}

func TestGame_PlayerWins(t *testing.T) {
	messages := make(chan message, 1)
	game := newGame(messages)

	myBoard := newBoard()
	myBoard.cells[0][0] = cell{shipClass: destroyer, status: statusUndefined}
	myBoard.cells[1][0] = cell{shipClass: destroyer, status: statusUndefined}

	game.myBoard = myBoard
	game.mode = waitingMode

	// opponent attacks first
	game.handleIncomingMessage(newAttackMessage(0, 0))

	got := <-messages

	assertEqual(t, got.t, response)
	assertEqual(t, got.status, statusHit)
	assertEqual(t, got.ship, destroyer)
	assertEqual(t, got.gameOver, false)
	assertEqual(t, game.mode, attackMode)

	// player attacks 0,0
	game.selectCellToAttack(terminal.EnterKey)

	got = <-messages

	assertEqual(t, got.t, attack)
	assertEqual(t, got.row, 0)
	assertEqual(t, got.col, 0)
	assertEqual(t, game.mode, waitingMode)

	// opponent response with hit
	game.handleIncomingMessage(newResponseMessageHit(0, 0, destroyer, true))

	assertEqual(t, game.targetBoard.cellAt(0, 0).shipClass, destroyer)
	assertEqual(t, game.targetBoard.cellAt(0, 0).status, statusHit)
	assertEqual(t, game.mode, winMode)
}

func assertEqual[T comparable](t *testing.T, got, expected T) {
	if got != expected {
		t.Errorf("got: %v, want: %v", got, expected)
	}
}
