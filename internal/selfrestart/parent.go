package selfrestart

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

const parentPipeEnv = "SPIRE_INTERNAL_PARENT_PIPE"

func watchParent() error {
	value := os.Getenv(parentPipeEnv)
	_ = os.Unsetenv(parentPipeEnv)
	number, err := strconv.ParseUint(value, 10, strconv.IntSize)
	if err != nil || number == 0 {
		return fmt.Errorf("invalid launcher liveness pipe")
	}
	pipe, err := openParentPipe(uintptr(number))
	if err != nil {
		return err
	}
	go func() {
		// Only the launcher owns the write end. EOF therefore means it exited,
		// including forced termination. This blocks without polling or timers.
		_, _ = io.Copy(io.Discard, pipe)
		_ = pipe.Close()
		// Shutdown must not wait for an inherited log pipe to drain.
		os.Exit(1)
	}()
	return nil
}
