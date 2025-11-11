package service

import (
	"context"
	"server/internal/models"
	"server/internal/repository/storage"
	package_log "server/pkg/logging"
)

type RoomService struct {
	logger *package_log.Logger
	repo   storage.RoomRepository
}

func NewRoomService(logger *package_log.Logger, repo storage.RoomRepository) *RoomService {
	return &RoomService{
		logger: logger,
		repo:   repo,
	}
}

func (s *RoomService) CreateRoom(ctx context.Context, data models.SingleRoom) (int, error) {
	return s.repo.CreateRoom(ctx, data)
}

func (s *RoomService) GetAllRooms(ctx context.Context) ([]string, error) {
	return s.repo.GetAllRooms(ctx)
}

func (s *RoomService) SaveMessage(ctx context.Context, dto models.SaveChat) error {
	return s.repo.SaveMessage(ctx, dto)
}
