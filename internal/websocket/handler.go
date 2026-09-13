package websocket

import (
	"errors"

	"golang.org/x/net/websocket"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HandleHello(ws *websocket.Conn, msg string) error {
	return websocket.Message.Send(ws, "Hello, Client!")
}

func (h *Handler) HandleUnauthorized(ws *websocket.Conn) error {
	err := websocket.Message.Send(ws, "{\"error\": \"Unauthorized\"}")
	if err != nil {
		return err
	}
	return errors.New("Unauthorized")
}
