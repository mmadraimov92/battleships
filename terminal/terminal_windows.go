//go:build windows

package terminal

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

type state struct {
	mode uint32
}

var (
	kernel32           = windows.NewLazySystemDLL("kernel32.dll")
	procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode = kernel32.NewProc("SetConsoleMode")
	procGetStdHandle   = kernel32.NewProc("GetStdHandle")
)

const (
	stdInputHandle  = ^uint32(10) // STD_INPUT_HANDLE
	stdOutputHandle = ^uint32(11) // STD_OUTPUT_HANDLE

	// Input mode flags
	enableEchoInput            = 0x0004
	enableInsertMode           = 0x0020
	enableLineInput            = 0x0002
	enableMouseInput           = 0x0010
	enableProcessedInput       = 0x0001
	enableQuickEditMode        = 0x0040
	enableWindowInput          = 0x0008
	enableVirtualTerminalInput = 0x0200

	// Output mode flags
	enableProcessedOutput           = 0x0001
	enableWrapAtEolOutput          = 0x0002
	enableVirtualTerminalProcessing = 0x0004
	disableNewlineAutoReturn       = 0x0008
	enableLvbGridWorldwide         = 0x0010
)

func getConsoleHandle(handleType uint32) (windows.Handle, error) {
	handle, _, err := procGetStdHandle.Call(uintptr(handleType))
	if handle == uintptr(windows.InvalidHandle) {
		return windows.InvalidHandle, err
	}
	return windows.Handle(handle), nil
}

func getConsoleMode(handle windows.Handle) (uint32, error) {
	var mode uint32
	ret, _, err := procGetConsoleMode.Call(uintptr(handle), uintptr(unsafe.Pointer(&mode)))
	if ret == 0 {
		return 0, err
	}
	return mode, nil
}

func setConsoleMode(handle windows.Handle, mode uint32) error {
	ret, _, err := procSetConsoleMode.Call(uintptr(handle), uintptr(mode))
	if ret == 0 {
		return err
	}
	return nil
}

func makeRaw(fd int) (*state, error) {
	// On Windows, fd is typically not used the same way as Unix
	// We work directly with console handles

	inputHandle, err := getConsoleHandle(stdInputHandle)
	if err != nil {
		return nil, fmt.Errorf("failed to get input handle: %w", err)
	}

	// Get current input mode
	currentMode, err := getConsoleMode(inputHandle)
	if err != nil {
		return nil, fmt.Errorf("failed to get console mode: %w", err)
	}

	// Save the original state
	oldState := &state{mode: currentMode}

	// Disable line input, echo input, and processed input to get raw input
	rawMode := currentMode
	rawMode &^= enableEchoInput
	rawMode &^= enableLineInput
	rawMode &^= enableProcessedInput
	rawMode &^= enableMouseInput
	rawMode &^= enableQuickEditMode
	rawMode &^= enableWindowInput

	// Enable virtual terminal input for better key handling
	rawMode |= enableVirtualTerminalInput

	// Set the new mode
	if err := setConsoleMode(inputHandle, rawMode); err != nil {
		return nil, fmt.Errorf("failed to set raw mode: %w", err)
	}

	// Also enable virtual terminal processing for output if possible
	outputHandle, err := getConsoleHandle(stdOutputHandle)
	if err == nil {
		if outputMode, err := getConsoleMode(outputHandle); err == nil {
			newOutputMode := outputMode | enableVirtualTerminalProcessing
			setConsoleMode(outputHandle, newOutputMode)
		}
	}

	// Note: On Windows, we don't need to set non-blocking mode like on Unix
	// The console API handles input differently

	return oldState, nil
}

func restore(fd int, state *state) error {
	inputHandle, err := getConsoleHandle(stdInputHandle)
	if err != nil {
		return fmt.Errorf("failed to get input handle: %w", err)
	}

	if err := setConsoleMode(inputHandle, state.mode); err != nil {
		return fmt.Errorf("failed to restore console mode: %w", err)
	}

	return nil
}
