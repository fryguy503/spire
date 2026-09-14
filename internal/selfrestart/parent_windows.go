package selfrestart

import (
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"golang.org/x/sys/windows"
)

func inheritParentPipe(cmd *exec.Cmd, reader *os.File) error {
	handle := windows.Handle(reader.Fd())
	if err := windows.SetHandleInformation(handle, windows.HANDLE_FLAG_INHERIT, windows.HANDLE_FLAG_INHERIT); err != nil {
		return err
	}
	cmd.Env = append(cmd.Env, parentPipeEnv+"="+strconv.FormatUint(uint64(handle), 10))
	cmd.SysProcAttr = &syscall.SysProcAttr{AdditionalInheritedHandles: []syscall.Handle{syscall.Handle(handle)}}
	return nil
}

func openParentPipe(fd uintptr) (*os.File, error) {
	if err := windows.SetHandleInformation(windows.Handle(fd), windows.HANDLE_FLAG_INHERIT, 0); err != nil {
		return nil, err
	}
	return os.NewFile(fd, "launcher-liveness"), nil
}
