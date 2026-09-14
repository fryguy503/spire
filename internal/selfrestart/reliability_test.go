package selfrestart

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestShutdownOfUnresponsiveWorkerIsBounded(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	ready := filepath.Join(t.TempDir(), "ready")
	stopping := make(chan os.Signal, 2)
	done := make(chan error, 1)
	go func() {
		_, err := supervise(executable, []string{"-test.run=^TestExitHelper$"}, append(os.Environ(), "SPIRE_TEST_WAIT_READY="+ready), t.TempDir(), 0, stopping)
		done <- err
	}()
	var worker *os.Process
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if data, err := os.ReadFile(ready); err == nil {
			pid, err := strconv.Atoi(string(data))
			if err != nil {
				t.Fatal(err)
			}
			worker, err = os.FindProcess(pid)
			if err != nil {
				t.Fatal(err)
			}
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if worker == nil {
		stopping <- os.Interrupt
		stopping <- os.Interrupt
		<-done
		t.Fatal("worker did not start")
	}
	defer worker.Release()
	stopping <- os.Interrupt
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(12 * time.Second):
		_ = worker.Kill()
		<-done
		t.Fatal("launcher waited indefinitely after a single stop request")
	}
}
