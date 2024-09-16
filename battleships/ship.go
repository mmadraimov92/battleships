package battleships

type shipClass int

const (
	empty shipClass = iota
	destroyer
	cruiser
	submarine
	battleship
	carrier
)

func (c shipClass) shipSize() int {
	shipSize := 0
	switch c {
	case carrier:
		shipSize = 5
	case battleship:
		shipSize = 4
	case submarine:
		shipSize = 3
	case cruiser:
		shipSize = 3
	case destroyer:
		shipSize = 2
	}

	return shipSize
}
