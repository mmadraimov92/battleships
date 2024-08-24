package app

import (
	"context"
	"fmt"
	"io"

	"tui/battleships"
	"tui/terminal"
)

type App struct {
	screen            io.Writer
	input             chan terminal.KeyEvent
	items             []Item
	selectedItemIndex uint8
}

func New(w io.Writer, inputChan chan terminal.KeyEvent, cancel context.CancelFunc) *App {
	battleships := battleships.New(w, inputChan)

	timer := timer{
		screen: w,
		input:  inputChan,
	}

	exit := exit{
		cancel: cancel,
	}

	return &App{
		screen: w,
		input:  inputChan,
		items:  []Item{battleships, &timer, &exit},
	}
}

func (m *App) Run(ctx context.Context) {
	m.draw(ctx, nil)

	for {
		select {
		case <-ctx.Done():
			return
		case keyEvent := <-m.input:
			m.draw(ctx, &keyEvent)
		}
	}
}

func (m *App) draw(ctx context.Context, pressedKey *terminal.KeyEvent) {
	if pressedKey != nil {
		switch *pressedKey {
		case terminal.UpArrowKey:
			if m.selectedItemIndex == 0 {
				m.selectedItemIndex = uint8(len(m.items)) - 1
			} else {
				m.selectedItemIndex--
			}
		case terminal.DownArrowKey:
			if m.selectedItemIndex == uint8(len(m.items))-1 {
				m.selectedItemIndex = 0
			} else {
				m.selectedItemIndex++
			}
		case terminal.EnterKey:
			m.items[m.selectedItemIndex].Render(ctx)
			m.draw(ctx, nil)
			return
		default:
			return
		}
	}

	// todo: call write only once per frame
	terminal.ClearScreen(m.screen)

	for i, item := range m.items {
		row := fmt.Sprintf("* %s", item.Title())
		if i == int(m.selectedItemIndex) {
			row += " <-"
		}
		fmt.Fprintln(m.screen, row)
	}
}
