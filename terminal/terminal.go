//go:build !darwin && !linux && !windows

package terminal

import (
	"fmt"
)

type state struct {
	// Dummy struct for unsupported platforms
}

func makeRaw(fd int) (*state, error) {
	return nil, fmt.Errorf("terminal raw mode not supported on this platform")
}

func restore(fd int, state *state) error {
	return fmt.Errorf("terminal restore not supported on this platform")
}
