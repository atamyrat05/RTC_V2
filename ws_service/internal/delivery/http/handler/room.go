package handler

import (
	"net/http"
	"server/internal/delivery/ws"
	"server/internal/helper"
	"server/internal/models"
	"server/internal/service/repository"
	package_log "server/pkg/logging"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	fromUserId := c.Query("fromUserId")
	toUserID := c.Query("userId")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cl := &ws.Client{
		Conn:     conn,
		Message:  make(chan *ws.Message, 10),
		ID:       fromUserId,
		RoomID:   1,
		ToUserId: toUserID,
		Service:  h.service,
		Logger:   h.logger,
	}

	m := &ws.Message{
		Content:    "A new user has joined the room",
		RoomID:     1,
		FromUserId: fromUserId,
		ToUserId:   toUserID,
	}

	h.hub.Register <- cl
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
	roomId, _ := strconv.Atoi(c.Param("roomId"))

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
