package desktop

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	spiremiddleware "github.com/EQEmuTools/spire/internal/http/middleware"
	"github.com/labstack/echo/v4"
)

func TestReadinessAcceptsBasicAuthChallenge(t *testing.T) {
	e := echo.New()
	e.Use(spiremiddleware.BasicAuth(func(username, password string, _ echo.Context) (bool, error) {
		return username == "test-user" && password == "test-password", nil
	}))
	e.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "ready") })
	server := httptest.NewServer(e)
	defer server.Close()
	if err := waitForSiteToBeAvailable(server.URL, time.Second); err != nil {
		t.Fatalf("healthy Basic Auth server rejected by desktop readiness: %v", err)
	}
}

func TestReadinessStopsAfterSuccess(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	if err := waitForSiteToBeAvailable(server.URL, time.Second); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)
	if got := requests.Load(); got != 1 {
		t.Fatalf("readiness continued requesting after success: %d requests", got)
	}
}

func TestReadinessRejectsServerErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusServiceUnavailable) }))
	defer server.Close()
	if err := waitForSiteToBeAvailable(server.URL, 50*time.Millisecond); err == nil {
		t.Fatal("HTTP 503 was treated as ready")
	}
}

func TestReadinessCancelsSlowRequest(t *testing.T) {
	canceled := make(chan struct{}, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			canceled <- struct{}{}
		case <-time.After(time.Second):
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()
	if err := waitForSiteToBeAvailable(server.URL, 50*time.Millisecond); err == nil {
		t.Fatal("slow server did not time out")
	}
	select {
	case <-canceled:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed-out readiness request is still running")
	}
}
