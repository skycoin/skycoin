//go:build tinygo && linux

package coloredcobra

import (
	"syscall"
	"unsafe"
)

// isTerminal reports whether fd refers to a terminal. TinyGo's host (native
// Linux) target uses the standard-library syscall package, so a TCGETS ioctl —
// which succeeds only on a tty — works here without pulling in go-isatty (whose
// Linux implementation is excluded under the tinygo build tag).
func isTerminal(fd uintptr) bool {
	var termios syscall.Termios
	_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, fd, syscall.TCGETS,
		uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	return errno == 0
}
