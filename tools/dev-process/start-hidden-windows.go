//go:build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

// Node starts this short-lived helper detached. CREATE_NO_WINDOW then gives
// the actual service and its descendants a console with no visible window.
// No supervisor remains running, and standard output goes directly to the logs.
func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "expected PID file, executable, and arguments")
		os.Exit(1)
	}
	// Check the handoff destination before starting anything. A missing or
	// unwritable directory must not leave an unreported process running.
	pidFile, err := os.OpenFile(os.Args[1], os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer pidFile.Close()
	command := exec.Command(os.Args[2], os.Args[3:]...)
	command.Stdin, command.Stdout, command.Stderr = os.Stdin, os.Stdout, os.Stderr
	command.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // CREATE_NO_WINDOW
	if err := command.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if _, err := pidFile.WriteString(strconv.Itoa(command.Process.Pid)); err != nil {
		_ = command.Process.Kill()
		_ = command.Wait()
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	_ = command.Process.Release()
}
