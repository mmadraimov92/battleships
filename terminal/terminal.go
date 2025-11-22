//go:build !darwin && !linux && !windows

package terminal

import (
	"fmt"
)

type state struct {
	// Dummy struct for unsupported platforms
}

func makeRaw(int) (*state, error) {
	return nil, fmt.Errorf("terminal raw mode not supported on this platform")
}

func restore(int, *state) error {
	return fmt.Errorf("terminal restore not supported on this platform")
}
