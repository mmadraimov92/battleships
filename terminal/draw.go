package terminal

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// https://en.wikipedia.org/wiki/ANSI_escape_code
var (
	control            = "\033["
	hideCursorSequence = control + "?25l"
	showCursorSequence = control + "?25h"
	clearSequence      = control + "H" + control + "2J"
	ResetSequence      = control + "0m"
	UnderlineSequence  = control + "4m"
	VerticalBar        = string(rune('\u2502'))
)

type renderer struct {
	buf *bytes.Buffer
	w   io.Writer
}

var r = renderer{
	buf: bytes.NewBuffer([]byte{}),
	w:   os.Stdout,
}

func Draw(s string) {
	_, err := r.buf.WriteString(s)
	if err != nil {
		panic(err)
	}
}

func Flush() {
	fmt.Fprintln(r.w, r.buf.String())
	r.buf.Reset()
}

func ClearScreen() {
	r.w.Write([]byte(clearSequence))
}

func HideCursor() {
	r.w.Write([]byte(hideCursorSequence))
}

func ShowCursor() {
	r.w.Write([]byte(showCursorSequence))
}

func MoveCursorTo(x, y int8) {
	Draw(fmt.Sprintf("%s%d;%dH", control, y, x))
}
