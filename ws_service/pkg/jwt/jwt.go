package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var securityKey = "qwerty123"
var tokenTTL = time.Hour * 24

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId string
}

func GenerateToken(userId string) (string, error) {
	claims := tokenClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(securityKey))
	if err != nil {
		return "", err
	}

	return t, nil
}

func ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid sign-in method")
		}
		return []byte(securityKey), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", errors.New("token claims are not type of *tokenClaims")
	}
	return claims.UserId, nil
}
