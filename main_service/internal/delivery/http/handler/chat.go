package handler

import (
	"net/http"
	"server/internal/helper"
	"server/internal/models"
	"server/internal/service/repository"
	package_log "server/pkg/logging"
	"strconv"

	"github.com/gin-gonic/gin"
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
