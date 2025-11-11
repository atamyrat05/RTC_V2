package ws

import (
	"context"
	"fmt"
	"log"
	"server/internal/models"
	"server/internal/service/repository"
	package_log "server/pkg/logging"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	ID       string `json:"id"`
	RoomID   int `json:"roomId"`
	ToUserId string `json:"to_user_id"`
	Logger   *package_log.Logger
	Service  repository.RoomService
}

type Message struct {
	Content    string `json:"content"`
	RoomID     int `json:"roomId"`
	FromUserId string `json:"from_user_id"`
	ToUserId   string `json:"to_user_id"`
}

func (c *Client) WriteMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}

		if message.Content == "user left the chat" {
			logrus.Printf("user left chat! roomId=%d, userId=%s", message.RoomID, message.FromUserId)
		}

		c.Conn.WriteJSON(message)
	}
}

func (c *Client) ReadMessage(hub *Hub, ctx context.Context) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := &Message{
			Content:    string(m),
			RoomID:     c.RoomID,
			FromUserId: c.ID,
			ToUserId:   c.ToUserId,
		}

		dto := models.SaveChat{
			Message:      string(m),
			From_user_id: c.ID,
			To_user_id:   c.ToUserId,
		}

		fmt.Println("message", dto.Message)
		err = c.Service.SaveMessage(ctx, dto)
		if err != nil {
			c.Logger.Errorln("error", err)
		}

		hub.Broadcast <- msg
	}
}
