package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"server/internal/helper"
	"server/internal/models"
	"server/internal/service/repository"
	package_log "server/pkg/logging"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	logger  *package_log.Logger
	service repository.ChatService
}

func NewChatHandler(logger *package_log.Logger, service repository.ChatService) *ChatHandler {
	return &ChatHandler{
		logger:  logger,
		service: service,
	}
}

func (h *ChatHandler) ChatRegisterRoutes(r *gin.RouterGroup) {
	r.GET("/get-users-chat-history", h.GetUserChatHistroies)
	r.GET("/get-user-chats/:partnerId", h.GetUserChatsWithPartner)
	r.GET("/join-ws", h.joinWS)
	r.GET("/join-ws-proxy", h.joinWSProxy)
}

func (h *ChatHandler) GetUserChatHistroies(c *gin.Context) {
	userId, err := helper.IntId(c)
	if err != nil {
		return
	}

	data, err := h.service.GetUserChatHistroies(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *ChatHandler) GetUserChatsWithPartner(c *gin.Context) {
	userId, err := helper.IntId(c)
	if err != nil {
		return
	}

	partnerId := c.Param("partnerId")
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	if limit == 0 || offset == 0 {
		limit = 20
		offset = 0
	}

	dto := models.ChatWithUser{
		UserID:    userId,
		PartnerID: partnerId,
		Limit:     limit,
		Offset:    offset,
	}

	data, err := h.service.GetUserChats(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *ChatHandler) joinWS(c *gin.Context) {
	userId, err := helper.IntId(c)
	if err != nil {
		h.logger.Errorln("error getting user id:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// foydalanuvchi kim bilan chat qilmoqchi (partner ID)
	toUserID := c.Query("partnerId")
	if toUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "partnerId is required"})
		return
	}

	// 2-serverdagi websocket manzili
	wsURL := url.URL{
		Scheme: "ws",
		Host:   "localhost:8081",
		Path:   "api/v1/join-ws",
		RawQuery: url.Values{
			"fromUserId": []string{userId},
			"userId":     []string{toUserID},
		}.Encode(),
	}

	// WebSocket ulanishini o‘rnatamiz
	conn, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		h.logger.Errorln("error connecting to second WS server:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect ws"})
		return
	}

	h.logger.Infof("✅ User %s connected to second WS server with %s", userId, wsURL.String())

	// real vaqtda xabarlarni olish
	go func() {
		defer conn.Close()
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				h.logger.Errorln("error reading ws message:", err)
				break
			}
			h.logger.Infof("📩 Message from 8081: %s", string(msg))
		}
	}()

	// test uchun 8081 ga bitta salom yuboramiz
	message := fmt.Sprintf("User %s joined chat with %s", userId, toUserID)
	err = conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		h.logger.Errorln("error sending message:", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "WebSocket connection established with server 8081",
	})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *ChatHandler) joinWSProxy(c *gin.Context) {
	token := c.Query("token")
	toUserID := c.Query("userId")

	if token == "" || toUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing token or userId"})
		return
	}

	// --- 1. Frontend bilan WS o‘rnatamiz
	clientConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Errorln("upgrade error:", err)
		return
	}
	h.logger.Infoln("✅ Frontend connected to 8080 WebSocket")

	// --- 2. 8081 dagi WS serverga ulanamiz
	wsURL := url.URL{
		Scheme: "ws",
		Host:   "localhost:8081",
		Path:   "api/v1/join-ws",
		RawQuery: url.Values{
			"token":  []string{token},
			"userId": []string{toUserID},
		}.Encode(),
	}

	serverConn, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		h.logger.Errorln("failed to connect 8081 WS:", err)
		clientConn.WriteMessage(websocket.TextMessage, []byte("❌ Failed to connect to chat service"))
		clientConn.Close()
		return
	}
	h.logger.Infoln("🔗 Connected to 8081 WebSocket server")

	// --- 3. Forward xabarlarni front <-> 8081 o‘rtasida
	go func() {
		defer clientConn.Close()
		defer serverConn.Close()
		for {
			msgType, msg, err := clientConn.ReadMessage()
			if err != nil {
				h.logger.Errorln("client read error:", err)
				break
			}
			h.logger.Infof("➡️ Front -> 8081: %s", msg)
			_ = serverConn.WriteMessage(msgType, msg)
		}
	}()

	go func() {
		defer clientConn.Close()
		defer serverConn.Close()
		for {
			msgType, msg, err := serverConn.ReadMessage()
			if err != nil {
				h.logger.Errorln("server read error:", err)
				break
			}
			h.logger.Infof("⬅️ 8081 -> Front: %s", msg)
			_ = clientConn.WriteMessage(msgType, msg)
		}
	}()

	// --- 4. Ulanish ochiq tursin
	for {
		time.Sleep(30 * time.Second)
		if clientConn == nil {
			break
		}
	}
}
