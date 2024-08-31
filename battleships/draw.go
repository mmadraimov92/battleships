package battleships

import (
	"fmt"

	"tui/terminal"
)

var (
	space   = byte(0x20)
	newLine = byte(0x0a)
	cols    = []byte{0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39} // 0-9
	rows    = []byte{0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a} // A-J
	water   = byte(0x7e)
	bigX    = byte(0x58)
	bigO    = byte(0x4f)
	plus    = byte(0x2b)
	two     = byte(0x32)
	three   = byte(0x33)
	four    = byte(0x34)
	five    = byte(0x35)
)

func (i *Item) draw(g *game) {
	terminal.ClearScreen(i.screen)

	screenBuf := []byte{}
	for range 8 {
		screenBuf = append(screenBuf, space)
	}
	screenBuf = append(screenBuf, []byte("My board")...)
	screenBuf = append(screenBuf, newLine)

	screenBuf = append(screenBuf, []byte{space, space}...)
	screenBuf = append(screenBuf, terminal.UnderlineSequence...)
	for _, col := range cols {
		screenBuf = append(screenBuf, []byte{col, space}...)
	}
	screenBuf = append(screenBuf, terminal.ResetSequence...)
	screenBuf = append(screenBuf, newLine)

	for i, row := range rows {
		screenBuf = append(screenBuf, row)
		screenBuf = append(screenBuf, terminal.VerticalBar...)
		for j := range len(cols) {
			symbolToDraw := cellSymbol(g.myBoard.cellAt(i, j))
			if i == int(g.myBoard.selectedRow.Current()) && j == int(g.myBoard.selectedCol.Current()) {
				symbolToDraw = plus
			}
			screenBuf = append(screenBuf, []byte{symbolToDraw, space}...)
		}
		screenBuf = append(screenBuf, newLine)
	}

	fmt.Fprintln(i.screen, string(screenBuf))
	fmt.Fprintln(i.screen, g.mode)

	if g.mode == Attack {
		fmt.Fprintln(
			i.screen,
			"Select cell to attack:",
			string(rows[g.myBoard.selectedRow.Current()]),
			string(cols[g.myBoard.selectedCol.Current()]),
		)
	}
}

func cellSymbol(c *cell) byte {
	if c.status != nil {
		switch *c.status {
		case statusHit:
			return bigX
		case statusMiss:
			return bigO
		}
	}
	switch c.shipClass {
	case destroyer:
		return two
	case cruiser, submarine:
		return three
	case battleship:
		return four
	case carrier:
		return five
	default:
		return water
	}
}
