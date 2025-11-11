package storage

import (
	"context"
	"server/internal/models"
)

type RoomRepository interface {
	CreateRoom(ctx context.Context, data models.SingleRoom) (string, error)
	GetAllRooms(ctx context.Context) ([]string, error)
	SaveMessage(ctx context.Context, dto models.SaveChat) error
}
