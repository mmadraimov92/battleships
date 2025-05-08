package battleships

import (
	"encoding/binary"
)

type messageType int8

const (
	initiative messageType = iota
	attack
	response
)

type message struct {
	row      int8
	col      int8
	t        messageType
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

func newInitiativeMessage(i int8) message {
	return message{
		t:   initiative,
		row: i,
	}
}

func encodeMessage(m message) []byte {
	var encoded uint16
	encoded |= uint16(m.row&0b1111) << 12 // 4 bits for row
	encoded |= uint16(m.col&0b1111) << 8  // 4 bits for col
	encoded |= uint16(m.t&0b11) << 6      // 2 bit for command
	encoded |= uint16(m.status&0b1) << 5  // 1 bit for hit/miss
	encoded |= uint16(m.ship&0b111) << 2  // 3 bits for ship class
	if m.gameOver {
		encoded |= 1 // last 1 bit for gameOver
	}

	result := make([]byte, 2)
	binary.BigEndian.PutUint16(result, encoded)
	return result
}

func decodeMessage(encoded []byte) message {
	if len(encoded) < 2 {
		return message{}
	}
	value := binary.BigEndian.Uint16(encoded)
	msg := message{
		row:      int8((value >> 12) & 0b1111),
		col:      int8((value >> 8) & 0b1111),
		t:        messageType((value >> 6) & 0b11),
		status:   cellStatus((value >> 5) & 0b1),
		ship:     shipClass((value >> 2) & 0b111),
		gameOver: (value & 0x1) != 0,
	}

	if msg.status == statusUndefined && msg.t == response {
		msg.status = statusMiss
	}

	return msg
}
