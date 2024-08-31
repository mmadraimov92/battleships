package terminal

import (
	"syscall"

	"golang.org/x/sys/unix"
)

type state struct {
	termios unix.Termios
}

const ioctlReadTermios = unix.TIOCGETA
const ioctlWriteTermios = unix.TIOCSETA

func makeRaw(fd int) (*state, error) {
	termios, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	if err != nil {
		return nil, err
	}
	if err = syscall.SetNonblock(fd, true); err != nil {
		return nil, err
	}
	oldState := state{termios: *termios}

	termios.Iflag &^= unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON
	termios.Lflag &^= unix.ECHO | unix.ECHONL | unix.ICANON | unix.IEXTEN
	termios.Cflag &^= unix.CSIZE | unix.PARENB
	termios.Cflag |= unix.CS8
	termios.Cc[unix.VMIN] = 1
	termios.Cc[unix.VTIME] = 0
	if err := unix.IoctlSetTermios(fd, ioctlWriteTermios, termios); err != nil {
		return nil, err
	}

	return &oldState, nil
}

func restore(fd int, state *state) error {
	return unix.IoctlSetTermios(fd, ioctlWriteTermios, &state.termios)
}
