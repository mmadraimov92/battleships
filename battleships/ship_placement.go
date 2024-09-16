package battleships

import "tui/cyclic"

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
		orientation:      cyclic.NewNumber(0, 3),
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

func (sp shipPlacement) placeOnBoard(b *board) {
	// clear previous cells with given class
	for row := range rows {
		for col := range cols {
			if sp.currentlyPlacing == b.cellAt(row, col).shipClass {
				b.cellAt(row, col).shipClass = empty
			}
		}
	}

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

func (sp shipPlacement) isValidPlacement(b *board) bool {
	row := b.selectedRow.Current()
	col := b.selectedCol.Current()

	shipSize := sp.currentlyPlacing.shipSize()

	// check borders
	switch sp.orientation.Current() {
	case 0:
		if row+shipSize > 10 {
			return false
		}
	case 1:
		if col-shipSize+1 < 0 {
			return false
		}
	case 2:
		if row-shipSize+1 < 0 {
			return false
		}
	case 3:
		if col+shipSize > 10 {
			return false
		}
	}

	// check collisions
	switch sp.orientation.Current() {
	case 0:
		for range shipSize {
			shipAtCell := b.cellAt(row, col).shipClass
			if shipAtCell != empty && shipAtCell != sp.currentlyPlacing {
				return false
			}
			row++
		}
		return true
	case 1:
		for range shipSize {
			shipAtCell := b.cellAt(row, col).shipClass
			if shipAtCell != empty && shipAtCell != sp.currentlyPlacing {
				return false
			}
			col--
		}
		return true
	case 2:
		for range shipSize {
			shipAtCell := b.cellAt(row, col).shipClass
			if shipAtCell != empty && shipAtCell != sp.currentlyPlacing {
				return false
			}
			row--
		}
		return true

	case 3:
		for range shipSize {
			shipAtCell := b.cellAt(row, col).shipClass
			if shipAtCell != empty && shipAtCell != sp.currentlyPlacing {
				return false
			}
			col++
		}
		return true
	}

	return true
}
