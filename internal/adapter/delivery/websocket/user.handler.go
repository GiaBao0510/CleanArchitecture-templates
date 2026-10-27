package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	domain "github.com/GiaBao0510/CleanArchitecture-templates/internal/domain/user"
	"github.com/coder/websocket"
)

type Handler struct{
	hub *Hub
}

func NewHandler(hub *Hub) *Handler{
	return &Handler{hub: hub}
}

func(h *Handler) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		// Configure OriginPatterns for trusted frontend origins
		// when deploying across origins.
	})
	if err != nil {
		return
	}
	defer conn.CloseNow()

	client := newClient(conn)
	h.hub.Register(client)
	defer h.hub.Unregister(client)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		client.writeLoop(ctx)
	}()

	for {
		if _, _, err := conn.Read(ctx); err != nil {
			break
		}
	}

	cancel()
	<- done
}

// Chuyển sự kiện từ Kafka sang WebSocket Hub
type UserEventHandler struct {
	hub *Hub
}

func NewUserEventHandler(hub *Hub) *UserEventHandler{
	return &UserEventHandler{hub: hub}
}

func(h *UserEventHandler) HandlerUserEvent(
	ctx context.Context,
	event domain.Event,
) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal websocket event: %w", err)
	}

	h.hub.Broadcast(ctx, data)
	return nil
}
