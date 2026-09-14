//go:build integration

package updater

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

// Run on each target OS: go test -tags=integration ./internal/updater -run TestUpdateAndRestart -v
// The fixture uses the real HTTP update handler, release selection, binary
// replacement, and launcher. Only GitHub and the database are substituted.
func TestUpdateAndRestart(t *testing.T) {
	for _, mode := range []string{"self", "managed"} {
		t.Run(mode, func(t *testing.T) { testUpdateAndRestart(t, mode) })
	}
}

func testUpdateAndRestart(t *testing.T, mode string) {
	dir := filepath.Join(t.TempDir(), "Spire install with spaces")
	if err := os.Mkdir(dir, 0755); err != nil {
		t.Fatal(err)
	}
	suffix := ""
	if runtime.GOOS == "windows" {
		suffix = ".exe"
	}
	live := filepath.Join(dir, "spire"+suffix)
	assetName := targetReleaseAssetName(runtime.GOOS, runtime.GOARCH)
	archives := map[string][]byte{}
	for _, version := range []string{"1.0.0", "2.0.0", "3.0.0"} {
		binary := filepath.Join(t.TempDir(), "spire"+suffix)
		if version == "1.0.0" {
			binary = live
		}
		build := exec.Command("go", "build", "-o", binary, "-ldflags", "-X main.version="+version, "./testdata/restart")
		if output, err := build.CombinedOutput(); err != nil {
			t.Fatalf("build fixture: %v\n%s", err, output)
		}
		content, err := os.ReadFile(binary)
		if err != nil {
			t.Fatal(err)
		}
		var buffer bytes.Buffer
		archive := zip.NewWriter(&buffer)
		header := &zip.FileHeader{Name: fmt.Sprintf("spire-%s-%s%s", runtime.GOOS, runtime.GOARCH, suffix), Method: zip.Store}
		header.SetMode(0755)
		entry, err := archive.CreateHeader(header)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := entry.Write(content); err != nil {
			t.Fatal(err)
		}
		if err := archive.Close(); err != nil {
			t.Fatal(err)
		}
		archives[version] = buffer.Bytes()
	}
	var mu sync.Mutex
	offered, invalidArchive := "1.0.0", false
	var hold, entered chan struct{}
	var releaseServer *httptest.Server
	releaseServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		version, invalid, blocked, received := offered, invalidArchive, hold, entered
		mu.Unlock()
		if r.URL.Path == "/asset.zip" {
			if invalid {
				_, _ = w.Write([]byte("invalid archive"))
				return
			}
			_, _ = w.Write(archives[version])
			return
		}
		if blocked != nil {
			received <- struct{}{}
			<-blocked
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]map[string]interface{}{{
			"tag_name": "v" + version, "prerelease": false,
			"assets": []map[string]string{{"name": assetName, "browser_download_url": releaseServer.URL + "/asset.zip"}},
		}})
	}))
	defer releaseServer.Close()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	_ = listener.Close()
	base := fmt.Sprintf("http://127.0.0.1:%d", port)
	args := []string{"http:serve", "--port", strconv.Itoa(port), "argument with spaces & quotes \" '"}
	command := exec.Command(live, args...)
	command.Dir = dir
	command.Env = append(os.Environ(), "APP_ENV=local", "SPIRE_RELEASE_REPO=Fixture/Spire", "SPIRE_TEST_RELEASE_URL="+releaseServer.URL, "SPIRE_TEST_MARKER=preserved")
	command.Env = append(command.Env, "SPIRE_RESTART_MODE="+mode)
	logFile, err := os.Create(filepath.Join(dir, "fixture.log"))
	if err != nil {
		t.Fatal(err)
	}
	defer logFile.Close()
	command.Stdout, command.Stderr = logFile, logFile
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- command.Wait() }()
	client := &http.Client{Timeout: 10 * time.Second}
	workerPID := 0
	launcherPID := command.Process.Pid
	finished := false
	t.Cleanup(func() {
		if !finished {
			_, _ = client.Post(base+"/stop", "application/json", nil)
			select {
			case <-done:
			case <-time.After(10 * time.Second):
				if workerPID > 0 {
					if p, err := os.FindProcess(workerPID); err == nil {
						_ = p.Kill()
					}
				}
				_ = command.Process.Kill()
				<-done
			}
		}
		// Windows self handoffs are no longer children of the test process.
		if runtime.GOOS == "windows" && mode == "self" && launcherPID != command.Process.Pid {
			_, _ = client.Post(base+"/stop", "application/json", nil)
			if p, err := os.FindProcess(launcherPID); err == nil {
				_, _ = p.Wait()
			}
		}
		if t.Failed() {
			log, _ := os.ReadFile(logFile.Name())
			t.Logf("fixture log:\n%s", log)
		}
	})
	type instance struct {
		Version         string   `json:"version"`
		PID             int      `json:"pid"`
		Parent          int      `json:"parent"`
		Cwd             string   `json:"cwd"`
		Args            []string `json:"args"`
		Port            int      `json:"port"`
		Resumed         bool     `json:"resumed"`
		Marker          string   `json:"marker"`
		WorkerEnv       string   `json:"worker_env"`
		LauncherVersion string   `json:"launcher_version"`
	}
	readInstance := func() (instance, error) {
		response, err := client.Get(base + "/api/v1/app/env")
		if err != nil {
			return instance{}, err
		}
		defer response.Body.Close()
		var body struct {
			Data instance `json:"data"`
		}
		err = json.NewDecoder(response.Body).Decode(&body)
		return body.Data, err
	}
	waitVersion := func(version string) instance {
		t.Helper()
		deadline := time.Now().Add(30 * time.Second)
		for time.Now().Before(deadline) {
			state, err := readInstance()
			if err == nil && state.Version == version {
				workerPID = state.PID
				return state
			}
			time.Sleep(50 * time.Millisecond)
		}
		t.Fatalf("Spire did not start version %s", version)
		return instance{}
	}
	postUpdate := func() (int, []byte, error) {
		response, err := client.Post(base+"/api/v1/app/update", "application/json", nil)
		if err != nil {
			return 0, nil, err
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		return response.StatusCode, body, err
	}
	original := waitVersion("1.0.0")
	if code, body, err := postUpdate(); err != nil || code != 200 || !bytes.Contains(body, []byte(`"updated":false`)) {
		t.Fatalf("no update: %d %s %v", code, body, err)
	}
	mu.Lock()
	offered, invalidArchive = "2.0.0", true
	mu.Unlock()
	if code, body, err := postUpdate(); err != nil || code != 500 {
		t.Fatalf("invalid archive: %d %s %v", code, body, err)
	}
	if state := waitVersion("1.0.0"); state.PID != original.PID {
		t.Fatal("failed update restarted Spire")
	}
	for _, version := range []string{"2.0.0", "3.0.0"} {
		mu.Lock()
		offered, invalidArchive = version, false
		hold, entered = make(chan struct{}), make(chan struct{}, 1)
		blocked, received := hold, entered
		mu.Unlock()
		type result struct {
			code int
			body []byte
			err  error
		}
		response := make(chan result, 1)
		go func() { code, body, err := postUpdate(); response <- result{code, body, err} }()
		select {
		case <-received:
		case <-time.After(10 * time.Second):
			t.Fatal("update did not reach release server")
		}
		if code, body, err := postUpdate(); err != nil || code != 409 {
			t.Fatalf("duplicate update: %d %s %v", code, body, err)
		}
		mu.Lock()
		hold, entered = nil, nil
		mu.Unlock()
		close(blocked)
		r := <-response
		if r.err != nil || r.code != 200 || !bytes.Contains(r.body, []byte(`"version":"`+version+`"`)) {
			t.Fatalf("install %s: %d %s %v", version, r.code, r.body, r.err)
		}
		if mode == "managed" {
			select {
			case err := <-done:
				if exit, ok := err.(*exec.ExitError); !ok || exit.ExitCode() != 75 {
					t.Fatalf("managed restart must request exit 75: %v", err)
				}
			case <-time.After(10 * time.Second):
				t.Fatal("managed launcher did not exit")
			}
			// Simulate a service manager starting its configured executable again.
			previous := command
			command = exec.Command(live, args...)
			command.Dir, command.Env = previous.Dir, previous.Env
			command.Stdout, command.Stderr = logFile, logFile
			if err := command.Start(); err != nil {
				t.Fatal(err)
			}
			go func() { done <- command.Wait() }()
		}
		state := waitVersion(version)
		launcherPID = state.Parent
		if state.PID == original.PID || state.Cwd != dir || !reflect.DeepEqual(state.Args, args) || state.Marker != "preserved" || state.WorkerEnv != "" || state.Port != port || state.LauncherVersion != version {
			t.Fatalf("restart did not preserve launch context: %+v", state)
		}
		if mode == "self" && !state.Resumed {
			t.Fatal("desktop restart state was not preserved")
		}
		if mode == "self" && runtime.GOOS != "windows" && state.Parent != original.Parent {
			t.Fatal("Linux launcher PID changed")
		}
		if mode == "self" && runtime.GOOS == "windows" && state.Parent == original.Parent {
			t.Fatal("Windows launcher was not replaced")
		}
		t.Logf("%s: worker PID %d -> %d, launcher PID %d, port %d", version, original.PID, state.PID, state.Parent, state.Port)
		original = state
	}
	if runtime.GOOS == "windows" && mode == "self" {
		// The original launcher exits successfully after handing off. Open a
		// handle to the current launcher before requesting its final shutdown.
		if err := <-done; err != nil {
			t.Fatalf("first launcher handoff: %v", err)
		}
		finished = true
		current, err := os.FindProcess(launcherPID)
		if err != nil {
			t.Fatal(err)
		}
		_, _ = client.Post(base+"/stop?code=23", "application/json", nil)
		state, err := current.Wait()
		if err != nil || state.ExitCode() != 23 {
			t.Fatalf("current launcher ordinary exit: %v %v", state, err)
		}
		launcherPID = command.Process.Pid
		return
	}
	_, _ = client.Post(base+"/stop?code=23", "application/json", nil)
	select {
	case err := <-done:
		finished = true
		if exit, ok := err.(*exec.ExitError); !ok || exit.ExitCode() != 23 {
			t.Fatalf("launcher must propagate ordinary exit: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("launcher restarted an ordinary exit")
	}
}
