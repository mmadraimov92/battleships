package menu

import (
	"context"
)

type exit struct {
	cancel context.CancelFunc
}

func NewExit(cancel context.CancelFunc) *exit {
	return &exit{cancel}
}

func (t *exit) Select(_ context.Context) {
	t.cancel()
}

func (t *exit) Title() string {
	return "Exit"
}
