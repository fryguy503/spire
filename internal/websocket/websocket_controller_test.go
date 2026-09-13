package websocket

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/EQEmuTools/spire/internal/logger"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

func TestWebsocketRejectsExecutionAndPreservesNotifications(t *testing.T) {
	// Buffered lifecycle channels let this test observe registration and cleanup
	// without starting the application's permanent broadcast loop.
	manager := &ClientManager{register: make(chan *Client, 1), unregister: make(chan *Client, 1)}
	controller := NewController(nil, NewHandler(), manager, logger.NewAppLogger())
	e := echo.New()
	for _, route := range controller.Routes() {
		e.Add(route.Method(), "/api/v1/"+route.Route(), route.Handler())
	}
	server := httptest.NewServer(e)
	t.Cleanup(server.Close)
	ws, err := websocket.Dial("ws"+strings.TrimPrefix(server.URL, "http")+"/api/v1/websocket", "", server.URL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = ws.Close() })
	if err := ws.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
		t.Fatal(err)
	}
	var client *Client
	select {
	case client = <-manager.register:
	case <-time.After(5 * time.Second):
		t.Fatal("websocket client was not registered")
	}
	for _, tc := range []struct{ message, response string }{
		{`{"action":"hello"}`, "Hello, Client!"},
		{`{"action":"exec_server_bin","command":"world","args":[]}`, `unsupported websocket action "exec_server_bin"`},
		{`{"action":"unknown"}`, `unsupported websocket action "unknown"`},
		{`{invalid`, "Error:"},
		{`{"action":"hello"}`, "Hello, Client!"},
	} {
		if err := websocket.Message.Send(ws, tc.message); err != nil {
			t.Fatal(err)
		}
		var response string
		if err := websocket.Message.Receive(ws, &response); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(response, tc.response) {
			t.Fatalf("message %s: got %q, want %q", tc.message, response, tc.response)
		}
	}
	// The same connection still accepts pushed server notifications.
	const notification = `{"type":"stopTimer","message":"1"}`
	if err := websocket.Message.Send(client.WS, notification); err != nil {
		t.Fatal(err)
	}
	var response string
	if err := websocket.Message.Receive(ws, &response); err != nil || response != notification {
		t.Fatalf("notification: %q, error: %v", response, err)
	}
	if err := ws.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case disconnected := <-manager.unregister:
		if disconnected != client {
			t.Fatal("unregistered a different client")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("disconnected client was not unregistered")
	}
	if len(manager.register) != 0 {
		t.Fatal("messages must not register the same client again")
	}
}
