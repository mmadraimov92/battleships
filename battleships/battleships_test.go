package battleships

import (
	"context"
	"io"
	"log/slog"
	"math"
	"os"
	"testing"
	"time"
	"tui/terminal"
)

func TestBattleships_SimulatorStarts(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	terminal.SetRendererOutput(io.Discard)
	simulator, err := NewSimulator()
	if err != nil {
		t.Fatal(err)
	}
	defer simulator.Close()

	input := make(chan terminal.KeyEvent)
	slog.SetLogLoggerLevel(slog.LevelDebug)
	b := New(input, slog.New(
		slog.
			NewTextHandler(os.Stderr, &slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			}),
	), WithAddress(simulator.Addr()))

	conn := b.connect(ctx)
	if conn == nil {
		t.Fatal("could not connect to simulator")
	}
	defer conn.Close()
	b.conn = conn

	go b.start(ctx)
	setupBoard(ctx, b)

	incoming, err := simulator.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, incoming.t, initiative)

	if simulator.Write(newInitiativeMessage(math.MaxInt8)) != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Millisecond)

	// send attack Hit
	if simulator.Write(newAttackMessage(0, 0)) != nil {
		t.Fatal(err)
	}
	incoming, err = simulator.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, incoming.t, response)
	assertEqual(t, incoming.status, statusHit)

	// receive attack
	pressDown(ctx, b)
	pressRight(ctx, b)
	pressEnter(ctx, b)

	incoming, err = simulator.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, incoming.t, attack)
	assertEqual(t, incoming.row, 1)
	assertEqual(t, incoming.col, 1)

	// send attack Miss
	if simulator.Write(newAttackMessage(7, 7)) != nil {
		t.Fatal(err)
	}
	incoming, err = simulator.Read(ctx)
	if err != nil {
		t.Fatal(err)
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
