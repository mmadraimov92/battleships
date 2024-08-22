package item

import (
	"context"
	"fmt"
	"io"

	"tui/terminal"
)

type exit struct {
	screen io.Writer
	input  chan terminal.KeyEvent
	cancel context.CancelFunc
}

func (t *exit) Render(_ context.Context) {
	t.cancel()
}

func (t *exit) Title() {
	fmt.Fprint(t.screen, "Exit")
}
