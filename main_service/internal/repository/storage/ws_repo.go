package storage

import (
	"context"
	"server/internal/models"
)

type WsRepository interface {
	SaveMessage(ctx context.Context, dto models.SaveChat) error
}
