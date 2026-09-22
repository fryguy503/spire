package selfrestart

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"golang.org/x/sys/windows"
)

const handoffHandleEnv = "SPIRE_INTERNAL_LAUNCHER_HANDLE"

func awaitLauncherHandoff() error {
	value := os.Getenv(handoffHandleEnv)
	if value == "" {
		return nil
	}
	_ = os.Unsetenv(handoffHandleEnv)
	number, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid launcher handoff: %w", err)
	}
	handle := windows.Handle(number)
	defer windows.CloseHandle(handle)
	// Use an inherited handle, not a PID that could have exited and been reused.
	event, err := windows.WaitForSingleObject(handle, 30000)
	if err != nil {
		return fmt.Errorf("wait for previous launcher: %w", err)
	}
	if event != windows.WAIT_OBJECT_0 {
		if event == uint32(windows.WAIT_TIMEOUT) {
			return fmt.Errorf("previous launcher did not exit within 30 seconds")
		}
		return fmt.Errorf("unexpected launcher wait result: %d", event)
	}
	return nil
}

func replaceLauncher(executable string, args, environment []string, directory string) (int, error) {
	// The replacement launcher waits for this process to exit before starting
	// its server. Service wrappers use SPIRE_RESTART_MODE=managed instead.
	handle, err := windows.OpenProcess(windows.SYNCHRONIZE, true, uint32(os.Getpid()))
	if err != nil {
		return 1, err
	}
	defer windows.CloseHandle(handle)
	cmd := exec.Command(executable, args...)
	cmd.Dir = directory
	cmd.Env = append(environment, handoffHandleEnv+"="+strconv.FormatUint(uint64(handle), 10))
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{AdditionalInheritedHandles: []syscall.Handle{syscall.Handle(handle)}}
	if err := cmd.Start(); err != nil {
		return 1, fmt.Errorf("start updated launcher: %w", err)
	}
	_ = cmd.Process.Release()
	return 0, nil
}
