package item

import (
	"context"
	"io"

	"tui/terminal"
)

type Item interface {
	Render(context.Context)
	Title()
}

func NewMainMenu(w io.Writer, inputChan chan terminal.KeyEvent, cancel context.CancelFunc) *MainMenu {
	timer := timer{
		screen: w,
		input:  inputChan,
	}

	exit := exit{
		screen: w,
		cancel: cancel,
	}

	return &MainMenu{
		screen:   w,
		input:    inputChan,
		subMenus: []Item{&timer, &exit},
	}
}
