package battleships

import (
	"context"
	"io"
	"log/slog"
	"math"
	"testing"
	"time"
	"tui/terminal"
)

func TestBattleships_SimulatorStarts(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()

	terminal.SetRendererOutput(io.Discard)
	simulator, conn := NewSimulator()

	input := make(chan terminal.KeyEvent)
	slog.SetLogLoggerLevel(slog.LevelDebug)
	b := New(input, conn, slog.Default())

	go b.start(ctx)
	setupBoard(ctx, b)

	incoming, err := simulator.Read(ctx)
	if err != nil {
		t.Error(err)
	}
	assertEqual(t, incoming.t, initiative)

	if simulator.Write(newInitiativeMessage(math.MaxInt8)) != nil {
		t.Error(err)
	}

	// send attack Hit
	if simulator.Write(newAttackMessage(0, 0)) != nil {
		t.Error(err)
	}
	incoming, err = simulator.Read(ctx)
	if err != nil {
		t.Error(err)
	}
	assertEqual(t, incoming.t, response)
	assertEqual(t, incoming.status, statusHit)

	// receive attack
	pressDown(ctx, b)
	pressRight(ctx, b)
	pressEnter(ctx, b)

	incoming, err = simulator.Read(ctx)
	if err != nil {
		t.Error(err)
	}
	assertEqual(t, incoming.t, attack)
	assertEqual(t, incoming.row, 1)
	assertEqual(t, incoming.col, 1)

	// send attack Miss
	if simulator.Write(newAttackMessage(7, 7)) != nil {
		t.Error(err)
	}
	incoming, err = simulator.Read(ctx)
	if err != nil {
		t.Error(err)
	}
	assertEqual(t, incoming.t, response)
	assertEqual(t, incoming.status, statusMiss)

	cancel()
	<-ctx.Done()
}

func setupBoard(ctx context.Context, b *Battleships) {
	for range 5 {
		select {
		case b.input <- terminal.EnterKey:
		case <-ctx.Done():
		}
	}
}

func pressDown(ctx context.Context, b *Battleships) {
	select {
	case b.input <- terminal.DownArrowKey:
	case <-ctx.Done():
	}
}

func pressRight(ctx context.Context, b *Battleships) {
	select {
	case b.input <- terminal.RightArrowKey:
	case <-ctx.Done():
	}
}

func pressEnter(ctx context.Context, b *Battleships) {
	select {
	case b.input <- terminal.EnterKey:
	case <-ctx.Done():
	}
}
