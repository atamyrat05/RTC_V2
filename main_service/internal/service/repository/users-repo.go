package repository

import (
	"context"
	"server/internal/models"
)

type UsersService interface {
	CreateUser(ctx context.Context, dto models.Registration) (models.AuthResponse, error)
	CheckUser(ctx context.Context, dto models.Registration) (models.AuthResponse, error)
}
