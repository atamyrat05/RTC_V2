package storage

import (
	"context"
	"server/internal/models"
)

type UsersRepository interface {
	CreateUser(ctx context.Context, dto models.Registration) (string, error)
	CheckUser(ctx context.Context, dto models.Registration) (string, error)
}
