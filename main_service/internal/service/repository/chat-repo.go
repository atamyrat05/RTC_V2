package repository

import (
	"context"
	"server/internal/models"
)

type ChatService interface {
	GetUserChatHistroies(ctx context.Context, user_id string) ([]models.Chat, error)
	GetUserChats(ctx context.Context, dto models.ChatWithUser) ([]models.ChatsMessage, error)
	SaveChats(ctx context.Context, dto models.SaveChat) error
}
