package app

import (
	"context"
)

type exit struct {
	cancel context.CancelFunc
}

func (t *exit) Render(_ context.Context) {
	t.cancel()
}

func (t *exit) Title() string {
	return "Exit"
}
