package battleships

import (
	"tui/cyclic"
	"tui/terminal"
)

var (
	statusHit  = "hit"
	statusMiss = "miss"
)

type board struct {
	cells       [10][10]cell
	selectedRow *cyclic.Number
	selectedCol *cyclic.Number
}

type cell struct {
	shipClass class
	status    *string
}

func newBoard() *board {
	return &board{
		cells:       [10][10]cell{},
		selectedRow: cyclic.NewNumber(0, 9),
		selectedCol: cyclic.NewNumber(0, 9),
	}
}

func (b *board) handleKeyEvent(k terminal.KeyEvent) {
	switch k {
	case terminal.UpArrowKey:
		b.selectedRow.Decrement()
	case terminal.DownArrowKey:
		b.selectedRow.Increment()
	case terminal.RightArrowKey:
		b.selectedCol.Increment()
	case terminal.LeftArrowKey:
		b.selectedCol.Decrement()
	}
}

func (b *board) cellAt(row, col int) *cell {
	return &b.cells[row][col]
}
