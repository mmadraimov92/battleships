package battleships

import (
	"tui/terminal"
)

var (
	space            = " "
	cols             = "0123456789"
	rows             = "ABCDEFGHIJ"
	water            = "~"
	hitSymbol        = "X"
	missSymbol       = "O"
	plus             = "+"
	destroyerSymbol  = "2"
	cruiserSymbol    = "3"
	submarineSymbol  = "3"
	battleshipSymbol = "4"
	carrierSymbol    = "5"
	verticalBar      = string(rune('\u2502'))
)

func draw(g *game) {
	terminal.ClearScreen()

	drawMyBoard(g.myBoard)
	drawTargetBoard(g.targetBoard)
	terminal.CursorNextLine()
	drawInfo(g)

	terminal.Flush()
}

func drawMyBoard(board *board) {
	for range 8 {
		terminal.Draw(space)
	}
	terminal.Draw("My board")
	drawCoordinates(1, 2)
	curX, curY := 3, 3
	terminal.MoveCursorTo(curX, curY)
	for i := range rows {
		for j := range len(cols) {
			symbolToDraw := cellSymbol(board.cellAt(i, j))
			if i == int(board.selectedRow.Current()) && j == int(board.selectedCol.Current()) {
				symbolToDraw = plus
			}
			terminal.Draw(symbolToDraw + space)
		}
		terminal.CursorDown()
		curY++
		terminal.MoveCursorTo(curX, curY)
	}
}

func drawTargetBoard(board *board) {
	offsetX := 25
	terminal.MoveCursorTo(offsetX, 1)
	for range 6 {
		terminal.Draw(space)
	}
	terminal.Draw("Target board")
	drawCoordinates(offsetX, 2)

	curX, curY := offsetX+2, 3
	terminal.MoveCursorTo(curX, curY)
	for i := range rows {
		for j := range len(cols) {
			symbolToDraw := cellSymbol(board.cellAt(i, j))
			if i == int(board.selectedRow.Current()) && j == int(board.selectedCol.Current()) {
				symbolToDraw = missSymbol
			}
			terminal.Draw(symbolToDraw + space)
		}
		terminal.CursorDown()
		curY++
		terminal.MoveCursorTo(curX, curY)
	}
}

func drawInfo(g *game) {
	terminal.Draw("Mode: " + string(g.mode))
	terminal.CursorNextLine()

	if g.mode == Attack {
		terminal.Draw(
			"Select cell to attack: " +
				string(rows[g.myBoard.selectedRow.Current()]) +
				string(cols[g.myBoard.selectedCol.Current()]))
	}
}

func drawCoordinates(offsetX, offsetY int) {
	curX, curY := offsetX, offsetY
	terminal.MoveCursorTo(curX, curY)
	terminal.Draw(space)
	terminal.Draw(space)
	terminal.Underline()
	for _, col := range cols {
		terminal.Draw(string(col))
		terminal.Draw(space)
	}
	terminal.ResetFormatting()
	terminal.CursorDown()
	curY++
	terminal.MoveCursorTo(curX, curY)
	for _, row := range rows {
		terminal.Draw(string(row) + verticalBar)
		terminal.CursorDown()
		curY++
		terminal.MoveCursorTo(curX, curY)
	}
}

func cellSymbol(c *cell) string {
	if c.status != nil {
		switch *c.status {
		case statusHit:
			return hitSymbol
		case statusMiss:
			return missSymbol
		}
	}
	switch c.shipClass {
	case destroyer:
		return destroyerSymbol
	case cruiser:
		return cruiserSymbol
	case submarine:
		return submarineSymbol
	case battleship:
		return battleshipSymbol
	case carrier:
		return carrierSymbol
	default:
		return water
	}
}
