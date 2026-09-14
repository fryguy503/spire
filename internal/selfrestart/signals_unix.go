//go:build !windows

package selfrestart

import (
	"os"
	"syscall"
)

func stopSignals() []os.Signal { return []os.Signal{os.Interrupt, syscall.SIGTERM} }

func forwardStop(process *os.Process, sig os.Signal) {
	_ = process.Signal(sig)
}
