package helper

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errMessage struct {
	Message string `json:"message"`
}

func newErrResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, errMessage{Message: message})
	logrus.Errorf("Error: %s", message)
}
