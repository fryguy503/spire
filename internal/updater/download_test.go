package updater

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

func TestUpdateDownloadDeadlineCancelsStalledResponse(t *testing.T) {
	for _, stage := range []string{"headers", "body"} {
		t.Run(stage, func(t *testing.T) {
			var attempts atomic.Int32
			canceled := make(chan struct{}, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				attempts.Add(1)
				if stage == "body" {
					w.Header().Set("Content-Length", "1000")
					_, _ = w.Write([]byte("partial download"))
					w.(http.Flusher).Flush()
				}
				<-r.Context().Done()
				canceled <- struct{}{}
			}))
			defer server.Close()
			ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
			defer cancel()
			destination := filepath.Join(t.TempDir(), "update.zip")
			if err := downloadUpdate(ctx, destination, server.URL); !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("stalled download returned %v", err)
			}
			select {
			case <-canceled:
			case <-time.After(time.Second):
				t.Fatal("download request survived its deadline")
			}
			if got := attempts.Load(); got != 1 {
				t.Fatalf("expired download retried %d times", got)
			}
			if stage == "body" {
				if err := os.Remove(destination); err != nil {
					t.Fatalf("partial download is still open: %v", err)
				}
			}
		})
	}
}

func TestUpdateDownloadRetriesWithClosedTruncatedFiles(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if attempts.Add(1) == 1 {
			w.Header().Set("Content-Length", "1000")
			_, _ = w.Write([]byte("partial download longer than the successful response"))
			return
		}
		_, _ = w.Write([]byte("complete"))
	}))
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	destination := filepath.Join(t.TempDir(), "update.zip")
	if err := downloadUpdate(ctx, destination, server.URL); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(destination)
	if err != nil || string(data) != "complete" || attempts.Load() != 2 {
		t.Fatalf("retried download = %q, attempts = %d, error = %v", data, attempts.Load(), err)
	}
}

func TestUpdateDownloadRejectsHTTPError(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		http.Error(w, "try later", http.StatusServiceUnavailable)
	}))
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	destination := filepath.Join(t.TempDir(), "update.zip")
	if err := downloadUpdate(ctx, destination, server.URL); err == nil {
		t.Fatal("HTTP 503 accepted as a release archive")
	}
	if attempts.Load() != 3 {
		t.Fatalf("attempts = %d", attempts.Load())
	}
	if _, err := os.Stat(destination); !os.IsNotExist(err) {
		t.Fatalf("HTTP error body was saved: %v", err)
	}
}
