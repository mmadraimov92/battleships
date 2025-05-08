package battleships

import (
	"testing"
)

func TestMessageEncodeDecode(t *testing.T) {
	tests := []struct {
		name        string
		m           message
		description string
	}{
		{
			name:        "Attack message",
			m:           newAttackMessage(5, 7),
			description: "Row 5, Col 7, Type attack",
		},
		{
			name:        "Miss response",
			m:           newResponseMessageMiss(3, 2),
			description: "Row 3, Col 2, Type response, Status miss",
		},
		{
			name:        "Hit response, not game over",
			m:           newResponseMessageHit(9, 4, destroyer, false),
			description: "Row 9, Col 4, Type response, Status hit, Ship destroyer, Not game over",
		},
		{
			name:        "Hit response, game over",
			m:           newResponseMessageHit(15, 15, carrier, true),
			description: "Row 15, Col 15, Type response, Status hit, Ship carrier, Game over",
		},
		{
			name:        "Initiative message",
			m:           newInitiativeMessage(5),
			description: "Initiative value 3, Type initiative",
		},
	}

	for _, tt := range tests {
		encoded := encodeMessage(tt.m)
		decoded := decodeMessage(encoded)

		if decoded.row != tt.m.row ||
			decoded.col != tt.m.col ||
			decoded.t != tt.m.t ||
			decoded.status != tt.m.status ||
			decoded.ship != tt.m.ship ||
			decoded.gameOver != tt.m.gameOver {
			t.Errorf("%s: Original: %+v, After encode/decode: %+v", tt.name, tt.m, decoded)
		}
	}
}
