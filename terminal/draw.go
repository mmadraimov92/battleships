package terminal

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"syscall"
	"time"
)

// https://en.wikipedia.org/wiki/ANSI_escape_code
var (
	control            = "\033["
	hideCursorSequence = control + "?25l"
	showCursorSequence = control + "?25h"
	cursorDownSequence = control + "1B"
	nextLineSequence   = control + "1E"
	clearSequence      = control + "H" + control + "2J"
	resetSequence      = control + "0m"
	underlineSequence  = control + "4m"
	invertSequence     = control + "7m"
)

type renderer struct {
	buf *bytes.Buffer
	w   io.Writer
	m   *sync.Mutex
}

var r = renderer{
	buf: bytes.NewBuffer([]byte{}),
	w:   os.Stdout,
	m:   &sync.Mutex{},
}

func SetRendererOutput(w io.Writer) {
	r.w = w
}

func Draw(s string) {
	r.m.Lock()
	defer r.m.Unlock()

	_, err := r.buf.WriteString(s)
	if err != nil {
		panic(err)
	}
	flush()
}

func flush() {
	for {
		_, err := r.buf.WriteTo(r.w)
		if err == nil {
			return
		}
		if errors.Is(err, syscall.EAGAIN) {
			time.Sleep(time.Millisecond)
			continue
		}

		panic(err)
	}
}

func CursorDown() {
	Draw(cursorDownSequence)
}

func CursorNextLine() {
	Draw(nextLineSequence)
}

func Underline() {
	Draw(underlineSequence)
}

func Invert() {
	Draw(invertSequence)
}

func ResetFormatting() {
	Draw(resetSequence)
}

func ClearScreen() {
	r.w.Write([]byte(clearSequence)) //nolint:errcheck
}

func HideCursor() {
	r.w.Write([]byte(hideCursorSequence)) //nolint:errcheck
}

func ShowCursor() {
	r.w.Write([]byte(showCursorSequence)) //nolint:errcheck
}

func MoveCursorTo(x, y int) {
	Draw(fmt.Sprintf("%s%d;%dH", control, y, x))
}
