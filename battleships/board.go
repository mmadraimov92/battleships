package battleships

import (
	"tui/cyclic"
)

type cellStatus int8

const (
	statusUndefined cellStatus = iota
	statusHit
	statusMiss
)

type board struct {
	cells       [10][10]cell
	selectedRow *cyclic.Number
	selectedCol *cyclic.Number
}

type cell struct {
	shipClass shipClass
	status    cellStatus
}

func newBoard() *board {
	return &board{
		cells:       [10][10]cell{},
		selectedRow: cyclic.NewNumber(9),
		selectedCol: cyclic.NewNumber(9),
	}
}

func (b *board) cellAt(row, col int8) *cell {
	return &b.cells[row][col]
}

func (b *board) markAsHit(row, col int8, cls shipClass) {
	b.cells[row][col].status = statusHit
	b.cells[row][col].shipClass = cls
}

func (b *board) markAsMiss(row, col int8) {
	b.cells[row][col].status = statusMiss
}

func (b *board) isDestroyed() bool {
	for _, row := range b.cells {
		for _, c := range row {
			if c.status == statusUndefined && c.shipClass != empty {
				return false
			}
		}
	}

	return true
}
