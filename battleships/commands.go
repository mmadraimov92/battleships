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

func encodeMessage(m message) uint16 {
	var encoded uint16
	encoded |= uint16(m.row&0xF) << 12   // 4 bits for row
	encoded |= uint16(m.col&0xF) << 8    // 4 bits for col
	encoded |= uint16(m.t&0x1) << 7      // 1 bit for command
	encoded |= uint16(m.status&0x1) << 6 // 1 bit for hit/miss
	encoded |= uint16(m.ship&0x7) << 3   // 3 bits for ship class
	if m.gameOver {
		encoded |= 1 // last 1 bit for gameOver
	}
	return encoded
}

func decodeMessage(encoded uint16) message {
	msg := message{
		row:      int8((encoded >> 12) & 0xF),
		col:      int8((encoded >> 8) & 0xF),
		t:        commandType((encoded >> 7) & 0x1),
		status:   cellStatus((encoded >> 6) & 0x1),
		ship:     shipClass((encoded >> 3) & 0x7),
		gameOver: (encoded & 0x1) != 0,
	}

	if msg.status == statusUndefined && msg.t == response {
		msg.status = statusMiss
	}

	return msg
}
