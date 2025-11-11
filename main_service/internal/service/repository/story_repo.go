package repository

import (
	"context"
	"server/internal/models"
)

type StoryService interface {
	SaveStory(ctx context.Context, id string, url string) error
	GetUserStory(ctx context.Context, id string) (models.Story, error)
	GetAllStoryForUser(ctx context.Context, dto models.GetStoryDto) ([]models.Story, error)
	GetStoryById(ctx context.Context, id string) ([]string, error)
}
