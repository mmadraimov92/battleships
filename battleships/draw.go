package battleships

import (
	"strings"

	"tui/terminal"
)

var (
	space            = " "
	newLine          = "\n"
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
)

func (i *Item) draw(g *game) {
	terminal.ClearScreen()

	str := strings.Builder{}
	for range 8 {
		str.WriteString(space)
	}
	str.WriteString("My board")
	str.WriteString(newLine)
	terminal.Draw(str.String())
	terminal.MoveCursorTo(1, 2)
	i.drawCoordinates()

	str = strings.Builder{}
	for i, row := range rows {
		str.WriteString(string(row))
		str.WriteString(terminal.VerticalBar)
		for j := range len(cols) {
			symbolToDraw := cellSymbol(g.myBoard.cellAt(i, j))
			if i == int(g.myBoard.selectedRow.Current()) && j == int(g.myBoard.selectedCol.Current()) {
				symbolToDraw = plus
			}
			str.WriteString(symbolToDraw)
			str.WriteString(space)
		}
		str.WriteString(newLine)
	}
	terminal.Draw(str.String())
	terminal.Draw(newLine)
	terminal.Draw("Mode: " + string(g.mode))
	terminal.Draw(newLine)

	if g.mode == Attack {
		terminal.Draw(
			"Select cell to attack: " +
				string(rows[g.myBoard.selectedRow.Current()]) +
				string(cols[g.myBoard.selectedCol.Current()]))
	}
	terminal.Flush()
}

func (i *Item) drawCoordinates() {
	str := strings.Builder{}
	str.WriteString(space)
	str.WriteString(space)
	str.WriteString(terminal.UnderlineSequence)
	for _, col := range cols {
		str.WriteString(string(col))
		str.WriteString(space)
	}
	str.WriteString(terminal.ResetSequence)
	str.WriteString(newLine)
	terminal.Draw(str.String())
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
