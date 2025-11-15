package handler

import (
	"server/internal/configs"
	"server/internal/delivery/http/handler"
	"server/internal/helper"
	"server/internal/repository"
	"server/internal/service"
	package_psql "server/pkg/db"
	package_log "server/pkg/logging"

	"github.com/gin-gonic/gin"
)

const (
	baseURL  = "/api/v1/"
	usersURL = baseURL + "/users"
	storyURL = baseURL + "/stories"
	chatURL  = baseURL + "/chat"
)

func Manager(logger *package_log.Logger, clientPsql package_psql.Client, cfg *configs.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Static("/api/v1/uploads", "./uploads")

	usersGroup := r.Group(usersURL)
	{
		userRepo := repository.NewUsersPsqlRepository(logger, clientPsql)
		userService := service.NewUsersService(logger, userRepo)
		userHandler := handler.NewUsersHandler(logger, userService)
		userHandler.UsersRegisterRoutes(usersGroup)
	}

	storyGroup := r.Group(storyURL, helper.UserIdentity)
	{
		storyRepo := repository.NewStoryPsqlRepository(logger, clientPsql)
		storyService := service.NewStoryService(logger, storyRepo)
		storyHandler := handler.NewStoryHandler(logger, storyService)
		storyHandler.StoryRegisterRoutes(storyGroup)
	}

	chatGroup := r.Group(chatURL)
	{
		chatRepo := repository.NewChatPsqlRepository(logger, clientPsql)
		chatService := service.NewChatService(logger, chatRepo)
		chatHandler := handler.NewChatHandler(logger, chatService)
		chatHandler.ChatRegisterRoutes(chatGroup)
	}

	return r
}
