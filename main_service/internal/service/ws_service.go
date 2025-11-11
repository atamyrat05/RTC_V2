package service

import (
	"context"
	"server/internal/models"
	"server/internal/repository/storage"
	package_log "server/pkg/logging"
)

type WsService struct {
	logger *package_log.Logger
	repo   storage.WsRepository
}

func NewWsService(logger *package_log.Logger, repo storage.WsRepository) *WsService {
	return &WsService{
		logger: logger,
		repo:   repo,
	}
}

func (s *WsService) SaveMessage(ctx context.Context, dto models.SaveChat) error {
	return s.repo.SaveMessage(ctx, dto)
}
