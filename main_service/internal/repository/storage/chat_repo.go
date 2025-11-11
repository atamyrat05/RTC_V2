package storage

import (
	"context"
	"server/internal/models"
)

type ChatRepository interface {
	GetUserChatHistroies(ctx context.Context, user_id string) ([]models.Chat, error)
	GetUserChats(ctx context.Context, data models.ChatWithUser) ([]models.ChatsMessage, error)
	SaveChats(ctx context.Context, dto models.SaveChat) error
}
