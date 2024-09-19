package battleships

import (
	"tui/cyclic"
)

var shipsToPlaceInOrder = []shipClass{
	carrier,
	battleship,
	submarine,
	cruiser,
	destroyer,
}

type shipPlacement struct {
	currentlyPlacing shipClass
	orientation      *cyclic.Number
	placed           []shipClass
}

func newShipPlacement() shipPlacement {
	return shipPlacement{
		currentlyPlacing: carrier,
		orientation:      cyclic.NewNumber(3),
	}
}

func (sp *shipPlacement) acceptPlacement() {
	sp.placed = append(sp.placed, sp.currentlyPlacing)

	// return if all 5 ships are placed
	if len(sp.placed) == 5 {
		return
	}

	// chose next ship to place
	for i, ship := range shipsToPlaceInOrder {
		if ship == sp.currentlyPlacing {
			sp.currentlyPlacing = shipsToPlaceInOrder[i+1]
			break
		}
	}
}

// todo: better ux for placement
// Place ship in given selected row/col and orientation.
// If the ship is colliding with another ship then shift col and row until free spot is found.
// If out of bounds do not do anything.
func (sp shipPlacement) placeInValidPosition(b *board) {
	// clear previous cells with given class
	for row := range rows {
		for col := range cols {
			if sp.currentlyPlacing == b.cellAt(int8(row), int8(col)).shipClass {
				b.cellAt(int8(row), int8(col)).shipClass = empty
			}
		}
	}

	for sp.isColliding(b) {
		if b.selectedCol.Current() == 9 {
			b.selectedRow.Increment()
			b.selectedCol.Increment()
			continue
		}
		if b.selectedRow.Current() == 9 {
			b.selectedRow.Increment()
			b.selectedCol.Increment()
			continue
		}

		b.selectedCol.Increment()
	}

	sp.placeOnBoard(b)
}

func (sp shipPlacement) placeOnBoard(b *board) {
	switch sp.orientation.Current() {
	case 0:
		sp.placeDown(b)
	case 1:
		sp.placeLeft(b)
	case 2:
		sp.placeUp(b)
	case 3:
		sp.placeRight(b)
	}
}

func (sp shipPlacement) placeDown(b *board) {
	row := b.selectedRow.Current()
	col := b.selectedCol.Current()
	shipSize := sp.currentlyPlacing.shipSize()

	b.cellAt(row, col).shipClass = sp.currentlyPlacing

	for range shipSize {
		b.cellAt(row, col).shipClass = sp.currentlyPlacing
		row++
	}
}

func (sp shipPlacement) placeUp(b *board) {
	row := b.selectedRow.Current()
	col := b.selectedCol.Current()
	shipSize := sp.currentlyPlacing.shipSize()

	b.cellAt(row, col).shipClass = sp.currentlyPlacing

	for range shipSize {
		b.cellAt(row, col).shipClass = sp.currentlyPlacing
		row--
	}
}

func (sp shipPlacement) placeLeft(b *board) {
	row := b.selectedRow.Current()
	col := b.selectedCol.Current()
	shipSize := sp.currentlyPlacing.shipSize()

	b.cellAt(row, col).shipClass = sp.currentlyPlacing

	for range shipSize {
		b.cellAt(row, col).shipClass = sp.currentlyPlacing
		col--
	}
}

func (sp shipPlacement) placeRight(b *board) {
	row := b.selectedRow.Current()
	col := b.selectedCol.Current()
	shipSize := sp.currentlyPlacing.shipSize()

	b.cellAt(row, col).shipClass = sp.currentlyPlacing

	for range shipSize {
		b.cellAt(row, col).shipClass = sp.currentlyPlacing
		col++
	}
}

func (sp shipPlacement) isOutOfBounds(b *board) bool {
	row := b.selectedRow.Current()
	col := b.selectedCol.Current()
	shipSize := sp.currentlyPlacing.shipSize()

	switch sp.orientation.Current() {
	case 0:
		if row+shipSize > 10 {
			return true
		}
	case 1:
		if col-shipSize+1 < 0 {
			return true
		}
	case 2:
		if row-shipSize+1 < 0 {
			return true
		}
	case 3:
		if col+shipSize > 10 {
			return true
		}
	}

	return false
}

func (sp shipPlacement) isColliding(b *board) bool {
	row := b.selectedRow.Current()
	col := b.selectedCol.Current()
	shipSize := sp.currentlyPlacing.shipSize()

	if sp.isOutOfBounds(b) {
		return true
	}

	switch sp.orientation.Current() {
	case 0:
		for range shipSize {
			shipAtCell := b.cellAt(row, col).shipClass
			if shipAtCell != empty && shipAtCell != sp.currentlyPlacing {
				return true
			}
			row++
		}
		return false
	case 1:
		for range shipSize {
			shipAtCell := b.cellAt(row, col).shipClass
			if shipAtCell != empty && shipAtCell != sp.currentlyPlacing {
				return true
			}
			col--
		}
		return false
	case 2:
		for range shipSize {
			shipAtCell := b.cellAt(row, col).shipClass
			if shipAtCell != empty && shipAtCell != sp.currentlyPlacing {
				return true
			}
			row--
		}
		return false

	case 3:
		for range shipSize {
			shipAtCell := b.cellAt(row, col).shipClass
			if shipAtCell != empty && shipAtCell != sp.currentlyPlacing {
				return true
			}
			col++
		}
		return false
	}

	return true
}
