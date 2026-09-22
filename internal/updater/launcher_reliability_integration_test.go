//go:build integration

package updater

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

type launcherInstance struct {
	PID    int32 `json:"pid"`
	Parent int32 `json:"parent"`
}

type launcherFixture struct {
	base    string
	cmd     *exec.Cmd
	done    chan error
	current launcherInstance
	client  *http.Client
	stopped bool
}

func startLauncherFixture(t *testing.T, binary string, direct bool, stderr *os.File) *launcherFixture {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	_ = listener.Close()
	f := &launcherFixture{base: "http://127.0.0.1:" + strconv.Itoa(port), done: make(chan error, 1), client: &http.Client{Timeout: time.Second}}
	f.cmd = exec.Command(binary, "http:serve", "--port", strconv.Itoa(port))
	f.cmd.Dir = t.TempDir()
	f.cmd.Env = append(os.Environ(), "APP_ENV=local", "SPIRE_RESTART_MODE=self", "SPIRE_TEST_RELEASE_URL=http://127.0.0.1:1", "SPIRE_TEST_DIRECT="+strconv.Itoa(boolInt(direct)))
	log, err := os.Create(filepath.Join(f.cmd.Dir, "launcher.log"))
	if err != nil {
		t.Fatal(err)
	}
	f.cmd.Stdout, f.cmd.Stderr = log, log
	if stderr != nil {
		f.cmd.Stderr = stderr
	}
	if err := f.cmd.Start(); err != nil {
		_ = log.Close()
		t.Fatal(err)
	}
	go func() { f.done <- f.cmd.Wait() }()
	t.Cleanup(func() {
		f.stop()
		_ = log.Close()
		if t.Failed() {
			content, _ := os.ReadFile(log.Name())
			t.Logf("launcher log:\n%s", content)
		}
	})
	return f
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func (f *launcherFixture) waitReady(previous int32) (launcherInstance, error) {
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		response, err := f.client.Get(f.base + "/api/v1/app/env")
		if err == nil {
			var body struct {
				Data launcherInstance `json:"data"`
			}
			err = json.NewDecoder(response.Body).Decode(&body)
			_ = response.Body.Close()
			if err == nil && body.Data.PID != 0 && body.Data.PID != previous {
				f.current = body.Data
				return body.Data, nil
			}
		}
		time.Sleep(5 * time.Millisecond)
	}
	return launcherInstance{}, fmt.Errorf("worker did not become ready")
}

func (f *launcherFixture) stop() {
	if f.stopped {
		return
	}
	f.stopped = true
	response, err := f.client.Post(f.base+"/stop", "application/json", nil)
	if err == nil {
		_ = response.Body.Close()
	}
	if f.current.PID != 0 && !waitProcessExit(f.current.PID, time.Second) {
		killTestProcess(f.current.PID)
	}
	// A direct fixture's parent is the test runner, which we must never stop.
	if f.current.Parent != 0 && f.current.Parent != int32(os.Getpid()) && !waitProcessExit(f.current.Parent, time.Second) {
		killTestProcess(f.current.Parent)
	}
	select {
	case <-f.done:
	default:
		_ = f.cmd.Process.Kill()
		<-f.done
	}
}

func killTestProcess(pid int32) {
	if p, err := os.FindProcess(int(pid)); err == nil {
		_ = p.Kill()
		_ = p.Release()
	}
}

func waitProcessExit(pid int32, timeout time.Duration) bool {
	p, err := process.NewProcess(pid)
	if err != nil {
		return true
	}
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		running, err := p.IsRunning()
		if err != nil || !running {
			return true
		}
		// An exited child awaiting its parent's Wait is already stopped.
		status, _ := p.Status()
		for _, value := range status {
			if value == "zombie" {
				return true
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	return false
}

func TestLauncherReliability(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "launcher-fixture")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	build := exec.Command("go", "build", "-o", binary, "./testdata/restart")
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build fixture: %v\n%s", err, output)
	}
	t.Run("parent_death", func(t *testing.T) {
		f := startLauncherFixture(t, binary, false, nil)
		state, err := f.waitReady(0)
		if err != nil {
			t.Fatal(err)
		}
		if err := f.cmd.Process.Kill(); err != nil {
			t.Fatal(err)
		}
		if !waitProcessExit(state.PID, 3*time.Second) {
			t.Fatal("worker survived the forced termination of its launcher")
		}
	})
	for _, operation := range []string{"parent_death", "restart"} {
		t.Run(operation+"_with_blocked_log", func(t *testing.T) {
			reader, writer, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			defer reader.Close()
			defer writer.Close()
			f := startLauncherFixture(t, binary, false, writer)
			state, err := f.waitReady(0)
			if err != nil {
				t.Fatal(err)
			}
			response, err := f.client.Post(f.base+"/block-log", "application/json", nil)
			if err != nil {
				t.Fatal(err)
			}
			_ = response.Body.Close()
			if _, err := reader.Read(make([]byte, 1)); err != nil {
				t.Fatal(err)
			}
			// The fixture writes much more than the pipe can hold. Stop draining it
			// to model a service log collector that has stalled.
			time.Sleep(50 * time.Millisecond)
			if operation == "parent_death" {
				if err := f.cmd.Process.Kill(); err != nil {
					t.Fatal(err)
				}
			} else {
				response, err := f.client.Post(f.base+"/restart", "application/json", nil)
				if err != nil {
					t.Fatal(err)
				}
				_ = response.Body.Close()
				if response.StatusCode != 200 {
					t.Fatalf("restart status %d", response.StatusCode)
				}
				if _, err := f.waitReady(state.PID); err != nil {
					t.Fatal("blocked logging prevented launcher restart: ", err)
				}
			}
			if !waitProcessExit(state.PID, 3*time.Second) {
				t.Fatal("blocked logging prevented worker shutdown")
			}
		})
	}
	t.Run("performance_and_restarts", func(t *testing.T) {
		report := map[string]interface{}{"os": runtime.GOOS, "architecture": runtime.GOARCH}
		var directTimes, launchedTimes []float64
		for i := 0; i < 10; i++ {
			for _, direct := range []bool{true, false} {
				start := time.Now()
				f := startLauncherFixture(t, binary, direct, nil)
				if _, err := f.waitReady(0); err != nil {
					t.Fatal(err)
				}
				elapsed := float64(time.Since(start).Microseconds()) / 1000
				if direct {
					directTimes = append(directTimes, elapsed)
				} else {
					launchedTimes = append(launchedTimes, elapsed)
				}
				f.stop()
			}
		}
		sort.Float64s(directTimes)
		sort.Float64s(launchedTimes)
		report["direct_start_median_ms"], report["launcher_start_median_ms"] = directTimes[5], launchedTimes[5]
		report["startup_overhead_median_ms"] = launchedTimes[5] - directTimes[5]
		f := startLauncherFixture(t, binary, false, nil)
		state, err := f.waitReady(0)
		if err != nil {
			t.Fatal(err)
		}
		parent, err := process.NewProcess(state.Parent)
		if err != nil {
			t.Fatal(err)
		}
		memory, err := parent.MemoryInfo()
		if err != nil {
			t.Fatal(err)
		}
		before, err := parent.Times()
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(2 * time.Second)
		after, err := parent.Times()
		if err != nil {
			t.Fatal(err)
		}
		idleCPU := (after.User + after.System - before.User - before.System) * 1000
		report["launcher_idle_cpu_ms_over_2s"], report["launcher_rss_bytes"] = idleCPU, memory.RSS
		if idleCPU > 100 {
			t.Fatalf("idle launcher consumed %.1fms CPU in 2 seconds", idleCPU)
		}
		initialFDs, fdErr := parent.NumFDs()
		var restartTimes []float64
		for i := 0; i < 50; i++ {
			previous := state
			start := time.Now()
			response, err := f.client.Post(f.base+"/restart", "application/json", nil)
			if err != nil {
				t.Fatal(err)
			}
			_ = response.Body.Close()
			if response.StatusCode != 200 {
				t.Fatalf("restart status %d", response.StatusCode)
			}
			state, err = f.waitReady(previous.PID)
			if err != nil {
				t.Fatal(err)
			}
			restartTimes = append(restartTimes, float64(time.Since(start).Microseconds())/1000)
			if !waitProcessExit(previous.PID, time.Second) {
				t.Fatal("previous worker survived restart")
			}
			if runtime.GOOS == "windows" && !waitProcessExit(previous.Parent, time.Second) {
				t.Fatal("previous launcher survived handoff")
			}
			if runtime.GOOS != "windows" && previous.Parent != state.Parent {
				t.Fatal("launcher PID changed on Linux")
			}
		}
		parent, err = process.NewProcess(state.Parent)
		if err != nil {
			t.Fatal(err)
		}
		finalMemory, err := parent.MemoryInfo()
		if err != nil {
			t.Fatal(err)
		}
		finalFDs, finalFDErr := parent.NumFDs()
		if fdErr == nil && finalFDErr == nil {
			report["initial_fds"], report["final_fds"] = initialFDs, finalFDs
			if finalFDs > initialFDs+2 {
				t.Fatalf("file descriptor growth: %d -> %d", initialFDs, finalFDs)
			}
		}
		sort.Float64s(restartTimes)
		report["restart_cycles"], report["restart_median_ms"], report["restart_p95_ms"], report["final_launcher_rss_bytes"] = len(restartTimes), restartTimes[25], restartTimes[47], finalMemory.RSS
		if finalMemory.RSS > memory.RSS+16*1024*1024 {
			t.Fatalf("launcher memory grew: %d -> %d", memory.RSS, finalMemory.RSS)
		}
		data, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("launcher measurements:\n%s", data)
		if output := os.Getenv("SPIRE_LAUNCHER_QA_REPORT"); output != "" {
			if err := os.WriteFile(output, data, 0644); err != nil {
				t.Fatal(err)
			}
		}
	})
}
