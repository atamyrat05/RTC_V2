package handler

import (
	"net/http"
	"server/internal/delivery/ws"
	"server/internal/helper"
	"server/internal/models"
	"server/internal/service/repository"
	"server/pkg/jwt"
	package_log "server/pkg/logging"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type RoomHandler struct {
	logger  *package_log.Logger
	service repository.RoomService
	hub     *ws.Hub
}

func NewRoomHandler(logger *package_log.Logger, service repository.RoomService, hub *ws.Hub) *RoomHandler {
	return &RoomHandler{
		logger:  logger,
		service: service,
		hub:     hub,
	}
}

func (h *RoomHandler) RoomRegisterRoutes(r *gin.RouterGroup) {
	room := r.Group("/room", helper.UserIdentity)
	{
		room.POST("/create-room", h.CreateRoom)
		room.GET("/get-room", h.GetRooms)
		room.GET("/get-clients/:roomId", h.GetClients)
	}

	r.GET("/join-ws", h.JoinRoom)
	r.POST("/user/:id", h.keyp)

}

func (h *RoomHandler) GetAllRooms() error {
	return nil
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var req models.CreateRoomReq

	userId, err := helper.IntId(c)
	if err != nil {
		h.logger.Errorln("error", err)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	body := models.SingleRoom{
		User1_id: userId,
		User2_id: req.ID,
	}

	roomId, err := h.service.CreateRoom(c.Request.Context(), body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	h.hub.Rooms[roomId] = &ws.Room{
		ID:      roomId,
		Clients: make(map[string]*ws.Client),
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"roomId": roomId,
	})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *RoomHandler) JoinRoom(c *gin.Context) {
	token := c.Query("token")

	fromUserId, err := jwt.ParseToken(token)
	if err != nil {
		h.logger.Errorln("error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	toUserID := c.Query("userId")

	cl := &ws.Client{
		Conn:     conn,
		Message:  make(chan *ws.Message, 10),
		ID:       fromUserId,
		RoomID:   "global",
		ToUserId: toUserID,
		Service:  h.service,
		Logger:   h.logger,
	}

	m := &ws.Message{
		Content:    "A new user has joined the room",
		RoomID:     "global",
		FromUserId: fromUserId,
		ToUserId:   toUserID,
	}

	h.hub.Register <- cl
	logrus.Printf("new user connected! roomId=%s, userId=%s", "global", fromUserId)
	h.hub.Broadcast <- m

	go cl.WriteMessage()
	cl.ReadMessage(h.hub, c.Request.Context())
}

func (h *RoomHandler) GetRooms(c *gin.Context) {
	rooms := make([]models.RoomRes, 0)

	for _, r := range h.hub.Rooms {
		rooms = append(rooms, models.RoomRes{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	c.JSON(http.StatusOK, rooms)
}

func (h *RoomHandler) GetClients(c *gin.Context) {
	var clients []models.ClientRes
	roomId := c.Param("roomId")

	if _, ok := h.hub.Rooms[roomId]; !ok {
		clients = make([]models.ClientRes, 0)
		c.JSON(http.StatusOK, clients)
	}

	for _, c := range h.hub.Rooms[roomId].Clients {
		clients = append(clients, models.ClientRes{
			ID:       c.ID,
			Username: "test",
		})
	}

	c.JSON(http.StatusOK, clients)
}

func (h *RoomHandler) keyp(c *gin.Context) {
	userId := c.Param("id")

	token, err := jwt.GenerateToken(userId)
	if err != nil {
		h.logger.Errorln("error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
