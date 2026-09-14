package selfrestart

import "os"

func stopSignals() []os.Signal { return []os.Signal{os.Interrupt} }

func forwardStop(_ *os.Process, _ os.Signal) {
	// Windows delivers console Ctrl+C to both processes. Process.Signal does
	// not support os.Interrupt on Windows; the child handles the console event.
}
