package app

import (
	"context"
	"fmt"
	"io"

	"tui/terminal"
)

type App struct {
	screen            io.Writer
	input             chan terminal.KeyEvent
	items             []Item
	selectedItemIndex uint8
}

func New(w io.Writer, inputChan chan terminal.KeyEvent, cancel context.CancelFunc) *App {
	timer := timer{
		screen: w,
		input:  inputChan,
	}

	exit := exit{
		screen: w,
		cancel: cancel,
	}

	return &App{
		screen: w,
		input:  inputChan,
		items:  []Item{&timer, &exit},
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
		case terminal.DownArrowKey:
			if m.selectedItemIndex == 0 {
				m.selectedItemIndex = uint8(len(m.items)) - 1
			} else {
				m.selectedItemIndex--
			}
		case terminal.UpArrowKey:
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

	terminal.ClearScreen(m.screen)

	for i, subMenus := range m.items {
		fmt.Fprint(m.screen, "* ")
		subMenus.Title()
		if i == int(m.selectedItemIndex) {
			fmt.Fprint(m.screen, " <-")
		}
		fmt.Fprint(m.screen, "\n")
	}
}
