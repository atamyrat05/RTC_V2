package handler

import (
	"net/http"
	"server/internal/models"
	"server/internal/service/repository"
	package_log "server/pkg/logging"

	"github.com/gin-gonic/gin"
)

type UsersHandler struct {
	logger  *package_log.Logger
	service repository.UsersService
}

func NewUsersHandler(logger *package_log.Logger, service repository.UsersService) *UsersHandler {
	return &UsersHandler{
		logger:  logger,
		service: service,
	}
}

func (h *UsersHandler) UsersRegisterRoutes(r *gin.RouterGroup) {
	r.POST("/sign-up", h.signUp)
	r.POST("/sign-in", h.signIn)
}

func (h *UsersHandler) signUp(c *gin.Context) {
	var dto models.Registration
	if err := c.BindJSON(&dto); err != nil {
		h.logger.Errorln("error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})

		return
	}

	data, err := h.service.CreateUser(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	body := models.RegistrResponse{
		Token: data.Token,
		Id: models.Id{
			Ids: data.Id,
		},
	}

	c.JSON(http.StatusOK, body)
}

func (h *UsersHandler) signIn(c *gin.Context) {
	var dto models.Registration
	if err := c.BindJSON(&dto); err != nil {
		h.logger.Errorln("error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})

		return
	}

	data, err := h.service.CheckUser(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	body := models.RegistrResponse{
		Token: data.Token,
		Id: models.Id{
			Ids: data.Id,
		},
	}

	c.JSON(http.StatusOK, body)
}
