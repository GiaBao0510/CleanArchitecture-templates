package websocket

import (
	"context"
	"time"

	"github.com/coder/websocket"
)

// Struct này để kết nối với server websocket
type Client struct {
	conn *websocket.Conn
	send chan []byte
}

func newClient(conn *websocket.Conn) *Client{
	return &Client{
		conn: conn,
		send: make(chan []byte, 32),
	}
}

// Hàm writeLoop:
func(c *Client) writeLoop(ctx context.Context) {
	for{
		select {
		case <- ctx.Done():
			return
		case msg := <- c.send:
			writeCtx, cancel := context.WithTimeout(ctx, 5 *time.Second)
			err := c.conn.Write(writeCtx, websocket.MessageText, msg)
			cancel()
			if err != nil {
				return
			}
		}		
	}
}