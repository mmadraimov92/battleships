package terminal

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"
)

type KeyEvent int8

const (
	empty KeyEvent = iota

	EnterKey
	EscapeKey
	DeleteKey
	UpArrowKey
	DownArrowKey
	LeftArrowKey
	RightArrowKey
	SmallRKey
)

var keyMap = map[byte]KeyEvent{
	0x41: UpArrowKey,
	0x42: DownArrowKey,
	0x43: RightArrowKey,
	0x44: LeftArrowKey,
	0x7f: DeleteKey,
	0x1b: EscapeKey,
	0x0a: EnterKey,
	0x0d: EnterKey,
	0x72: SmallRKey,
}

func HandleKeyboardInput(ctx context.Context, input chan KeyEvent) error {
	oldState, err := makeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("HandleKeyboardInput: %w", err)
	}
	defer restore(int(os.Stdin.Fd()), oldState)

	buf := make([]byte, 3)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.Tick(20 * time.Millisecond):
			n, err := os.Stdin.Read(buf)
			if err != nil {
				if errors.Is(err, syscall.EAGAIN) {
					continue
				}
				fmt.Fprint(os.Stdout, err.Error())
			}
			keyEvent := processInput(buf, n)
			if keyEvent != empty {
				select {
				case input <- keyEvent:
				default:
				}
			}
		}
	}
}

func processInput(input []byte, bytesRead int) KeyEvent {
	if bytesRead == 1 {
		key, ok := keyMap[input[0]]
		if !ok {
			return empty
		}
		return key
	}

	if bytesRead == 3 && hasControlSequence(input) {
		key, ok := keyMap[input[2]]
		if !ok {
			return empty
		}
		return key
	}

	return empty
}

func hasControlSequence(input []byte) bool {
	return input[0] == control[0] && input[1] == control[1]
}
