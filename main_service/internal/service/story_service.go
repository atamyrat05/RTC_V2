package service

import (
	"context"
	"server/internal/models"
	"server/internal/repository/storage"
	package_log "server/pkg/logging"
)

type StoryService struct {
	logger *package_log.Logger
	repo   storage.StoryRepository
}

func NewStoryService(logger *package_log.Logger, repo storage.StoryRepository) *StoryService {
	return &StoryService{
		logger: logger,
		repo:   repo,
	}
}

func (s *StoryService) SaveStory(ctx context.Context, id string, url string) error {
	return s.repo.SaveStory(ctx, id, url)
}

func (s *StoryService) GetUserStory(ctx context.Context, id string) (models.Story, error) {
	return s.repo.GetUserStory(ctx, id)
}

func (s *StoryService) GetAllStoryForUser(ctx context.Context, dto models.GetStoryDto) ([]models.Story, error) {
	return s.repo.GetAllStoryForUser(ctx, dto)
}

func (s *StoryService) GetStoryById(ctx context.Context, id string) ([]string, error) {
	return s.repo.GetStoryById(ctx, id)
}
