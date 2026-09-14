//go:build !windows

package selfrestart

import "syscall"

func awaitLauncherHandoff() error { return nil }

func replaceLauncher(executable string, args, environment []string, _ string) (int, error) {
	// Re-exec the updated launcher without changing the PID tracked by systemd,
	// Docker, or another parent. The old server process has already exited.
	err := syscall.Exec(executable, append([]string{executable}, args...), environment)
	return 1, err
}
