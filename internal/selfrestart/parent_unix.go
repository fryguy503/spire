//go:build !windows

package selfrestart

import (
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

func inheritParentPipe(cmd *exec.Cmd, reader *os.File) error {
	cmd.Env = append(cmd.Env, parentPipeEnv+"="+strconv.Itoa(3+len(cmd.ExtraFiles)))
	cmd.ExtraFiles = append(cmd.ExtraFiles, reader)
	return nil
}

func openParentPipe(fd uintptr) (*os.File, error) {
	syscall.CloseOnExec(int(fd))
	return os.NewFile(fd, "launcher-liveness"), nil
}
