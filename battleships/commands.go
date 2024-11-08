package battleships

type commandType int8

const (
	attack commandType = iota
	response
)

type message struct {
	row      int8
	col      int8
	t        commandType
	status   cellStatus
	ship     shipClass
	gameOver bool
}

func newAttackMessage(row, col int8) message {
	return message{
		t:   attack,
		row: row,
		col: col,
	}
}

func newResponseMessageMiss(row, col int8) message {
	return message{
		t:      response,
		row:    row,
		col:    col,
		status: statusMiss,
	}
}

func newResponseMessageHit(row, col int8, ship shipClass, gameOver bool) message {
	return message{
		t:        response,
		row:      row,
		col:      col,
		status:   statusHit,
		ship:     ship,
		gameOver: gameOver,
	}
}
