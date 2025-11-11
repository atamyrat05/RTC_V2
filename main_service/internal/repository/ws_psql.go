package repository

import (
	"context"
	"server/internal/models"
	package_psql "server/pkg/db"
	package_log "server/pkg/logging"

	"github.com/jackc/pgx/v5"
)

type WsPsqlRepository struct {
	logger *package_log.Logger
	client package_psql.Client
}

func NewWsPsqlRepository(logger *package_log.Logger, client package_psql.Client) *WsPsqlRepository {
	return &WsPsqlRepository{
		logger: logger,
		client: client,
	}
}

func (r *WsPsqlRepository) SaveMessage(ctx context.Context, dto models.SaveChat) error {
	args := pgx.NamedArgs{
		"message":      dto.Message,
		"from_user_id": dto.From_user_id,
		"to_user_id":   dto.To_user_id,
	}

	q := `
		INSERT INTO messages (
			from_user_id, to_user_id, content
		) VALUES (
		 	@message, @from_user_id, @to_user_id
		)`

	_, err := r.client.Exec(ctx, q, args)
	if err != nil {
		r.logger.Errorln("error:", err)
		return err
	}

	return nil
}
