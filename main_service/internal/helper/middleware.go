package helper

import (
	"errors"
	"net/http"
	"server/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	auth    = "Authorization"
	UserCtx = "userId"
)

func UserIdentity(c *gin.Context) {
	header := c.GetHeader(auth)
	if header == "" {
		newErrResponse(c, http.StatusUnauthorized, "Auth header is empty")
		return
	}

	headerPart := strings.Split(header, " ")
	if len(headerPart) != 2 {
		newErrResponse(c, http.StatusUnauthorized, "Invalid auth header")
		return
	}
	userId, err := jwt.ParseToken(headerPart[1])
	if err != nil {
		newErrResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(UserCtx, userId)
}

func IntId(c *gin.Context) (string, error) {
	id, ok := c.Get(UserCtx)
	if !ok {
		newErrResponse(c, http.StatusInternalServerError, "user id not found!")
		return "", errors.New("user id not found")
	}

	userId, ok := id.(string)
	if !ok {
		newErrResponse(c, http.StatusInternalServerError, "user id is invalid type!")
		return "", errors.New("user id is invalid type")
	}

	return userId, nil
}
