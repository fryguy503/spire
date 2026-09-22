package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/EQEmuTools/spire/internal/http/routes"
	"github.com/EQEmuTools/spire/internal/logger"
	"github.com/EQEmuTools/spire/internal/pathmgmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
	"net/http"
)

type Controller struct {
	pathmgmt *pathmgmt.PathManagement
	handler  *Handler
	manager  *ClientManager
	logger   *logger.AppLogger
}

func NewController(
	pathmgmt *pathmgmt.PathManagement,
	handler *Handler,
	manager *ClientManager,
	logger *logger.AppLogger,
) *Controller {
	return &Controller{
		pathmgmt: pathmgmt,
		handler:  handler,
		manager:  manager,
		logger:   logger,
	}
}

func (a *Controller) Routes() []*routes.Route {
	return []*routes.Route{
		routes.RegisterRoute(http.MethodGet, "websocket", a.websocketHandler, nil),
	}
}

type SpireWebsocketMessage struct {
	Action string `json:"action"`
}

func (a *Controller) websocketHandler(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		client := &Client{WS: ws, ID: ws.Request().RemoteAddr}
		a.manager.register <- client
		defer func() { a.manager.unregister <- client }()

		for {
			var msg string
			if err := websocket.Message.Receive(ws, &msg); err != nil {
				return
			}
			if msg == "" {
				continue
			}

			a.logger.Debug().Any("msg", msg).Msg("Received message")
			var m SpireWebsocketMessage
			err := json.Unmarshal([]byte(msg), &m)
			if err == nil {
				switch m.Action {
				case "hello":
					err = a.handler.HandleHello(ws, msg)
				default:
					// Server mutations belong to the permission-checked HTTP APIs.
					err = fmt.Errorf("unsupported websocket action %q", m.Action)
				}
			}
			if err != nil {
				if sendErr := websocket.Message.Send(ws, fmt.Sprintf("Error: %v", err)); sendErr != nil {
					return
				}
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
