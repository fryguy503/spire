// Package selfrestart replaces Spire and its launcher after application updates.
// Ordinary exits and crashes are left to the user's process manager.
package selfrestart

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	restartExitCode = 75
	stateEnv        = "SPIRE_INTERNAL_RESTART_STATE"
	portEnv         = "SPIRE_INTERNAL_RESTART_PORT"
	launcherPortEnv = "SPIRE_INTERNAL_LAUNCHER_PORT"
)

var statePath string
var desktopPort int
var resumedDesktop bool

type restartState struct {
	DesktopPort int  `json:"desktop_port"`
	Managed     bool `json:"managed"`
}

// Run must be called before application initialization. The launcher does not
// open database connections, listen on ports, or run background services.
func Run() (handled bool, exitCode int, err error) {
	if err := awaitLauncherHandoff(); err != nil {
		return true, 1, err
	}
	if statePath = os.Getenv(stateEnv); statePath != "" {
		if err := watchParent(); err != nil {
			return true, 1, err
		}
		desktopPort, _ = strconv.Atoi(os.Getenv(portEnv))
		resumedDesktop = desktopPort > 0
		// Commands launched by Spire must not inherit the worker protocol.
		_ = os.Unsetenv(stateEnv)
		_ = os.Unsetenv(portEnv)
		return false, 0, nil
	}
	if len(os.Args) > 1 && os.Args[1] != "http:serve" {
		return false, 0, nil
	}
	executable, err := os.Executable()
	if err != nil {
		return true, 1, err
	}
	// Resolve once, before an update renames the running executable.
	executable, err = filepath.EvalSymlinks(executable)
	if err != nil {
		return true, 1, err
	}
	directory, err := os.Getwd()
	if err != nil {
		return true, 1, err
	}
	port, _ := strconv.Atoi(os.Getenv(launcherPortEnv))
	_ = os.Unsetenv(launcherPortEnv)
	stopping := make(chan os.Signal, 2)
	signal.Notify(stopping, stopSignals()...)
	defer signal.Stop(stopping)
	exitCode, err = supervise(executable, os.Args[1:], os.Environ(), directory, port, stopping)
	return true, exitCode, err
}

func Enabled() bool { return statePath != "" }

// DesktopPort returns the port of the previous desktop instance, if any.
func DesktopPort() int        { return desktopPort }
func SetDesktopPort(port int) { desktopPort = port }
func ResumedDesktop() bool    { return resumedDesktop }

// Prepare verifies the restart handoff before an HTTP handler reports success.
func Prepare() error {
	if !Enabled() {
		return errors.New("automatic restart is unavailable; start Spire directly or with http:serve")
	}
	mode := strings.ToLower(strings.TrimSpace(os.Getenv("SPIRE_RESTART_MODE")))
	if mode != "" && mode != "self" && mode != "managed" {
		return errors.New("SPIRE_RESTART_MODE must be self or managed")
	}
	data, err := json.Marshal(restartState{DesktopPort: desktopPort, Managed: mode == "managed"})
	if err != nil {
		return err
	}
	return os.WriteFile(statePath, data, 0600)
}

// Exit asks the launcher to restart after this process has completely exited.
// Call Prepare first, and flush any HTTP response before calling Exit.
func Exit() { os.Exit(restartExitCode) }

func supervise(executable string, args, environment []string, directory string, port int, stopping <-chan os.Signal) (int, error) {
	stateDir, err := os.MkdirTemp("", "spire-restart-")
	if err != nil {
		return 1, err
	}
	defer os.RemoveAll(stateDir)
	stateFile := filepath.Join(stateDir, "state.json")
	select {
	case <-stopping:
		return 0, nil
	default:
	}
	cmd := exec.Command(executable, args...)
	cmd.Dir = directory
	cmd.Env = append(append([]string{}, environment...), stateEnv+"="+stateFile, portEnv+"="+strconv.Itoa(port))
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	parentReader, parentWriter, err := os.Pipe()
	if err != nil {
		return 1, err
	}
	defer parentReader.Close()
	defer parentWriter.Close()
	if err := inheritParentPipe(cmd, parentReader); err != nil {
		return 1, err
	}
	if err := cmd.Start(); err != nil {
		return 1, fmt.Errorf("start Spire: %w", err)
	}
	_ = parentReader.Close()
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	stopRequested := false
	var shutdownTimer *time.Timer
	var shutdownDeadline <-chan time.Time
	defer func() {
		if shutdownTimer != nil {
			shutdownTimer.Stop()
		}
	}()
	waiting := true
	for waiting {
		select {
		case sig := <-stopping:
			if stopRequested {
				_ = cmd.Process.Kill()
			} else {
				stopRequested = true
				forwardStop(cmd.Process, sig)
				shutdownTimer = time.NewTimer(10 * time.Second)
				shutdownDeadline = shutdownTimer.C
			}
		case <-shutdownDeadline:
			// A stalled log collector must not prevent forced shutdown.
			_ = cmd.Process.Kill()
			shutdownDeadline = nil
		case err = <-done:
			waiting = false
		}
	}
	if stopRequested {
		return 0, nil
	}
	if err == nil {
		return 0, nil
	}
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		return 1, err
	}
	if exitErr.ExitCode() != restartExitCode {
		code := exitErr.ExitCode()
		if code < 0 {
			code = 1
		}
		return code, nil
	}
	data, err := os.ReadFile(stateFile)
	if err != nil {
		return 1, fmt.Errorf("read Spire restart handoff: %w", err)
	}
	var state restartState
	if err := json.Unmarshal(data, &state); err != nil {
		return 1, fmt.Errorf("read Spire restart handoff: %w", err)
	}
	if err := os.RemoveAll(stateDir); err != nil {
		return 1, err
	}
	select {
	case <-stopping:
		return 0, nil
	default:
	}
	if state.Managed {
		return restartExitCode, nil
	}
	return replaceLauncher(executable, args, append(environment, launcherPortEnv+"="+strconv.Itoa(state.DesktopPort)), directory)
}
