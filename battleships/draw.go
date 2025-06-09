package battleships

import (
	"tui/terminal"
)

const (
	space            = " "
	cols             = "0123456789"
	rows             = "ABCDEFGHIJ"
	water            = "~"
	hitSymbol        = "X"
	missSymbol       = "O"
	selectedCell     = "X"
	destroyerSymbol  = "2"
	cruiserSymbol    = "3"
	submarineSymbol  = "3"
	battleshipSymbol = "4"
	carrierSymbol    = "5"
	verticalBar      = string(rune('\u2502'))
)

func draw(g *game) {
	terminal.ClearScreen()

	drawMyBoard(g.myBoard, g.mode)
	drawTargetBoard(g.targetBoard, g.mode)
	drawInfo(g)
}

func drawMyBoard(board *board, m mode) {
	for range 8 {
		terminal.Draw(space)
	}
	terminal.Draw("My board")
	drawCoordinates(1, 2)
	curX, curY := 3, 3
	terminal.MoveCursorTo(curX, curY)

	for i := range rows {
		for j := range len(cols) {
			if m == preparationMode && i == int(board.selectedRow.Current()) && j == int(board.selectedCol.Current()) {
				terminal.Invert()
				terminal.Draw(cellSymbol(board.cellAt(int8(i), int8(j))))
				terminal.ResetFormatting()
				terminal.Draw(space)
				continue
			}

			terminal.Draw(cellSymbol(board.cellAt(int8(i), int8(j))) + space)
		}
		terminal.CursorDown()
		curY++
		terminal.MoveCursorTo(curX, curY)
	}
}

func drawTargetBoard(board *board, m mode) {
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
			if m == attackMode && i == int(board.selectedRow.Current()) && j == int(board.selectedCol.Current()) {
				terminal.Invert()
				terminal.Draw(cellSymbol(board.cellAt(int8(i), int8(j))))
				terminal.ResetFormatting()
				terminal.Draw(space)
				continue
			}

			terminal.Draw(cellSymbol(board.cellAt(int8(i), int8(j))) + space)
		}
		terminal.CursorDown()
		curY++
		terminal.MoveCursorTo(curX, curY)
	}
}

func drawInfo(g *game) {
	terminal.CursorNextLine()

	if g.mode == readyMode {
		terminal.Draw("Waiting for game to start")
		terminal.CursorNextLine()
	}

	if g.mode == waitingMode {
		terminal.Draw("Waiting for the attack")
		terminal.CursorNextLine()
	}

	if g.mode == preparationMode {
		terminal.Draw("Place your ships: ")
		terminal.Draw(g.shipPlacementInfo())
		terminal.CursorNextLine()
		terminal.CursorNextLine()
		terminal.Draw("Controls:")
		terminal.CursorNextLine()
		terminal.Draw("Arrow keys: move ship")
		terminal.CursorNextLine()
		terminal.Draw("R: rotate")
		terminal.CursorNextLine()
		terminal.Draw("Enter: place ship")
		terminal.CursorNextLine()
	}

	if g.mode == attackMode {
		terminal.Draw(
			"Select cell to attack: " +
				string(rows[g.targetBoard.selectedRow.Current()]) +
				string(cols[g.targetBoard.selectedCol.Current()]))
	}

	if g.mode == winMode {
		terminal.Draw("YOU WIN!!!")
	}

	if g.mode == loseMode {
		terminal.Draw("YOU LOST!!!")
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
	if c.status != statusUndefined {
		switch c.status {
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
