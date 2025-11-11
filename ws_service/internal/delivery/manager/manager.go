package handler

import (
	"server/internal/configs"
	"server/internal/delivery/http/handler"
	"server/internal/delivery/ws"
	"server/internal/repository"
	"server/internal/service"
	package_psql "server/pkg/db"
	package_log "server/pkg/logging"

	"github.com/gin-gonic/gin"
)

const (
	baseURL = "/api/v1/"
	roomURL = baseURL + "/room"
)

func Manager(logger *package_log.Logger, clientPsql package_psql.Client, cfg *configs.Config, hub *ws.Hub) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Static("/api/v1/uploads", "./uploads")

	roomGroup := r.Group(roomURL)
	{
		roomRepo := repository.NewRoomPsqlRepository(logger, clientPsql)
		roomService := service.NewRoomService(logger, roomRepo)
		roomHandler := handler.NewRoomHandler(logger, roomService, hub)
		roomHandler.RoomRegisterRoutes(roomGroup)
	}

	return r
}
