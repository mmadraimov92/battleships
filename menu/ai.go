package menu

import (
	"context"
	"time"
	"tui/terminal"
)

type battleshipsAI struct {
}

func NewBattleshipsAI() *battleshipsAI {
	return &battleshipsAI{}
}

func (t *battleshipsAI) Select(_ context.Context) {
	terminal.ClearScreen()
	terminal.Draw("Coming soon... WIP")
	time.Sleep(2 * time.Second)
}

func (t *battleshipsAI) Title() string {
	return "Battleships AI"
}
