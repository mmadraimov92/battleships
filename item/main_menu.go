package item

import (
	"context"
	"fmt"
	"io"

	"tui/terminal"
)

type MainMenu struct {
	screen            io.Writer
	input             chan terminal.KeyEvent
	subMenus          []Item
	selectedItemIndex uint8
}

func (m *MainMenu) Render(ctx context.Context) {
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

func (m *MainMenu) draw(ctx context.Context, pressedKey *terminal.KeyEvent) {
	if pressedKey != nil {
		switch *pressedKey {
		case terminal.DownArrowKey:
			if m.selectedItemIndex == 0 {
				m.selectedItemIndex = uint8(len(m.subMenus)) - 1
			} else {
				m.selectedItemIndex--
			}
		case terminal.UpArrowKey:
			if m.selectedItemIndex == uint8(len(m.subMenus))-1 {
				m.selectedItemIndex = 0
			} else {
				m.selectedItemIndex++
			}
		case terminal.EnterKey:
			m.subMenus[m.selectedItemIndex].Render(ctx)
			m.draw(ctx, nil)
			return
		default:
			return
		}
	}

	terminal.Clear(m.screen)

	for i, subMenus := range m.subMenus {
		fmt.Fprint(m.screen, "* ")
		subMenus.Title()
		if i == int(m.selectedItemIndex) {
			fmt.Fprint(m.screen, " <-")
		}
		fmt.Fprint(m.screen, "\n")
	}
}
