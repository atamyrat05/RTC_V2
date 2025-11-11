package service

import (
	"context"
	"server/internal/models"
	"server/internal/repository/storage"
	package_log "server/pkg/logging"
)

type ChatService struct {
	logger *package_log.Logger
	repo   storage.ChatRepository
}

func NewChatService(logger *package_log.Logger, repo storage.ChatRepository) *ChatService {
	return &ChatService{
		logger: logger,
		repo:   repo,
	}
}

func (s *ChatService) GetUserChatHistroies(ctx context.Context, user_id string) ([]models.Chat, error) {
	return s.repo.GetUserChatHistroies(ctx, user_id)
}

func (s *ChatService) GetUserChats(ctx context.Context, dto models.ChatWithUser) ([]models.ChatsMessage, error) {
	return s.repo.GetUserChats(ctx, dto)
}

func (s *ChatService) SaveChats(ctx context.Context, dto models.SaveChat) error {
	return s.repo.SaveChats(ctx, dto)
}
