package menu

import (
	"context"
	"fmt"

	"tui/cyclic"
	"tui/terminal"
)

type Item interface {
	Select(context.Context)
	Title() string
}

type menu struct {
	input        chan terminal.KeyEvent
	items        []Item
	selectedItem *cyclic.Number
}

func New(input chan terminal.KeyEvent, items []Item) *menu {
	return &menu{
		input:        input,
		items:        items,
		selectedItem: cyclic.NewNumber(int8(len(items) - 1)),
	}
}

func (m *menu) Run(ctx context.Context) {
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

func (m *menu) draw(ctx context.Context, pressedKey *terminal.KeyEvent) {
	if ctx.Err() != nil {
		return
	}

	if pressedKey != nil {
		switch *pressedKey {
		case terminal.UpArrowKey:
			m.selectedItem.Decrement()
		case terminal.DownArrowKey:
			m.selectedItem.Increment()
		case terminal.EnterKey:
			m.items[m.selectedItem.Current()].Select(ctx)
			return
		default:
			return
		}
	}

	terminal.ClearScreen()
	for i, item := range m.items {
		row := fmt.Sprintf("* %s", item.Title())
		if i == int(m.selectedItem.Current()) {
			row += " <-"
		}
		row += "\n"
		terminal.Draw(row)
	}
}
