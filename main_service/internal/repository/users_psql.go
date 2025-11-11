package repository

import (
	"context"
	"server/internal/models"
	package_psql "server/pkg/db"
	package_log "server/pkg/logging"

	"github.com/jackc/pgx/v5"
)

type UsersPsqlRepository struct {
	logger *package_log.Logger
	client package_psql.Client
}

func NewUsersPsqlRepository(logger *package_log.Logger, client package_psql.Client) *UsersPsqlRepository {
	return &UsersPsqlRepository{
		logger: logger,
		client: client,
	}
}

func (r *UsersPsqlRepository) CreateUser(ctx context.Context, dto models.Registration) (string, error) {
	var uuid string
	args := pgx.NamedArgs{
		"username": dto.Username,
		"password": dto.Password,
	}

	q := `INSERT INTO users (username, password) VALUES (@username, @password)	RETURNING uuid`

	err := r.client.QueryRow(ctx, q, args).Scan(&uuid)
	if err != nil {
		r.logger.Errorln("error:", err)
		return "", err
	}

	return uuid, nil
}

func (r *UsersPsqlRepository) CheckUser(ctx context.Context, dto models.Registration) (string, error) {
	var userId string
	args := pgx.NamedArgs{
		"username": dto.Username,
		"password": dto.Password,
	}

	q := `
		SELECT uuid FROM users WHERE username = @username AND password = @password`

	err := r.client.QueryRow(ctx, q, args).Scan(&userId)
	if err != nil {
		r.logger.Errorln("error", err)
		return "", err
	}

	return userId, nil
}
