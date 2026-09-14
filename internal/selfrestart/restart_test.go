package selfrestart

import (
	"encoding/json"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestPrepareRestartHandoff(t *testing.T) {
	previousPath, previousPort, previousResumed := statePath, desktopPort, resumedDesktop
	t.Cleanup(func() { statePath, desktopPort, resumedDesktop = previousPath, previousPort, previousResumed })
	file := filepath.Join(t.TempDir(), "state.json")
	statePath, desktopPort = file, 8094
	t.Setenv("SPIRE_RESTART_MODE", "managed")
	if err := Prepare(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	var state restartState
	if err := json.Unmarshal(data, &state); err != nil {
		t.Fatal(err)
	}
	if state.DesktopPort != 8094 || !state.Managed {
		t.Fatalf("restart state = %+v", state)
	}
	t.Setenv("SPIRE_RESTART_MODE", "unknown")
	if err := Prepare(); err == nil {
		t.Fatal("invalid restart mode was accepted")
	}
	statePath = ""
	if err := Prepare(); err == nil {
		t.Fatal("unmanaged process accepted a restart")
	}
}

func TestPrepareAfterTemporaryDirectoryCleanup(t *testing.T) {
	previous := statePath
	t.Cleanup(func() { statePath = previous })
	t.Setenv("SPIRE_RESTART_MODE", "self")
	directory := filepath.Join(t.TempDir(), "spire-restart-test")
	if err := os.Mkdir(directory, 0700); err != nil {
		t.Fatal(err)
	}
	statePath = filepath.Join(directory, "state.json")
	if err := os.Remove(directory); err != nil {
		t.Fatal(err)
	}
	if err := Prepare(); err != nil {
		t.Fatalf("temporary-directory cleanup disabled future updates: %v", err)
	}
	data, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	var state restartState
	if err := json.Unmarshal(data, &state); err != nil {
		t.Fatalf("recreated handoff is unreadable: %v", err)
	}
}

func TestSupervisorPropagatesExitWithoutRestart(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	for _, code := range []int{0, 23, restartExitCode} {
		environment := append(os.Environ(), "SPIRE_TEST_EXIT_CODE="+strconv.Itoa(code))
		got, err := supervise(executable, []string{"-test.run=^TestExitHelper$"}, environment, t.TempDir(), 0, make(chan os.Signal))
		if code == restartExitCode {
			if err == nil || got != 1 {
				t.Fatalf("restart without handoff = %d, %v", got, err)
			}
		} else if err != nil || got != code {
			t.Fatalf("exit %d = %d, %v", code, got, err)
		}
	}
}

func TestExitHelper(t *testing.T) {
	if ready := os.Getenv("SPIRE_TEST_WAIT_READY"); ready != "" {
		_, _, _ = Run()
		ignored := make(chan os.Signal, 2)
		signal.Notify(ignored, stopSignals()...)
		if err := os.WriteFile(ready, []byte(strconv.Itoa(os.Getpid())), 0600); err != nil {
			t.Fatal(err)
		}
		for {
			time.Sleep(time.Hour)
		}
	}
	if value := os.Getenv("SPIRE_TEST_EXIT_CODE"); value != "" {
		code, err := strconv.Atoi(value)
		if err != nil {
			t.Fatal(err)
		}
		os.Exit(code)
	}
}

func TestStopBeforeLaunchDoesNotStartWorker(t *testing.T) {
	stopping := make(chan os.Signal, 1)
	stopping <- os.Interrupt
	if code, err := supervise("missing-executable", nil, nil, t.TempDir(), 0, stopping); code != 0 || err != nil {
		t.Fatalf("stopped launcher = %d, %v", code, err)
	}
}
