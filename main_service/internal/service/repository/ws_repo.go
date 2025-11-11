package repository

import (
	"context"
	"server/internal/models"
)

type WsService interface {
	SaveMessage(ctx context.Context, dto models.SaveChat) error
}
