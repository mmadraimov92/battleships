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
)

// https://en.wikipedia.org/wiki/ANSI_escape_code
var (
	control            = []byte{'\x1b', '\x5b'}
	hideCursorSequence = append(control, []byte{'\x3f', '\x32', '\x35', '\x6c'}...)
	showCursorSequence = append(control, []byte{'\x3f', '\x32', '\x35', '\x68'}...)
	clearSequence      = append(control, []byte{'\x48', '\x1b', '\x5b', '\x32', '\x4a'}...)

	keyMap = map[byte]KeyEvent{
		'\x41': UpArrowKey,
		'\x42': DownArrowKey,
		'\x43': RightArrowKey,
		'\x44': LeftArrowKey,
		'\x7f': DeleteKey,
		'\x1b': EscapeKey,
		'\r':   EnterKey,
		'\n':   EnterKey,
	}
)

func HandleKeyboardInput(ctx context.Context, input chan KeyEvent) {
	oldState, err := makeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer restore(int(os.Stdin.Fd()), oldState)

	buf := make([]byte, 3)

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.Tick(50 * time.Millisecond):
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
