package app

import (
	"context"
	"fmt"
	"io"
)

type exit struct {
	screen io.Writer
	cancel context.CancelFunc
}

func (t *exit) Render(_ context.Context) {
	t.cancel()
}

func (t *exit) Title() {
	fmt.Fprint(t.screen, "Exit")
}
