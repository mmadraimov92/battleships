package app

import "context"

type Item interface {
	Render(context.Context)
	Title()
}
