package repository

import (
	"context"
	"errors"
	"server/internal/models"
	package_psql "server/pkg/db"
	package_log "server/pkg/logging"

	"github.com/jackc/pgx/v5"
)

type RoomPsqlRepository struct {
	logger *package_log.Logger
	client package_psql.Client
}

func NewRoomPsqlRepository(logger *package_log.Logger, client package_psql.Client) *RoomPsqlRepository {
	return &RoomPsqlRepository{
		logger: logger,
		client: client,
	}
}

func (r *RoomPsqlRepository) CreateRoom(ctx context.Context, data models.SingleRoom) (string, error) {
	var roomId string
	args := pgx.NamedArgs{
		"user1_id": data.User1_id,
		"user2_id": data.User2_id,
	}

	q := `
		SELECT uuid
		FROM rooms
		WHERE (user_1 = @user1_id AND user_2 = @user2_id)
   		OR (user_1 = @user2_id AND user_2 = @user1_id)
		`

	err := r.client.QueryRow(ctx, q, args).Scan(&roomId)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		r.logger.Errorln("error", err)
		return "", err
	}

	if roomId == "" {
		q = `INSERT INTO rooms (user_1, user_2) VALUES (@user1_id, @user2_id) RETURNING uuid`

		err = r.client.QueryRow(ctx, q, args).Scan(&roomId)
		if err != nil {
			r.logger.Errorln("error", err)
			return "", err
		}
	}

	return roomId, nil
}

func (r *RoomPsqlRepository) GetAllRooms(ctx context.Context) ([]string, error) {
	var data []string
	q := `SELECT uuid FROM rooms`

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		r.logger.Errorln("error", err)
		return nil, err
	}

	for rows.Next() {
		var uuid string

		err = rows.Scan(&uuid)
		if err != nil {
			r.logger.Errorln("error", err)
			return nil, err
		}

		data = append(data, uuid)
	}

	return data, nil
}

func (r *RoomPsqlRepository) SaveMessage(ctx context.Context, dto models.SaveChat) error {
	args := pgx.NamedArgs{
		"message":      dto.Message,
		"from_user_id": dto.From_user_id,
		"to_user_id":   dto.To_user_id,
	}

	q := `
		INSERT INTO messages (
			from_user_id, to_user_id, content
		) VALUES (
		 	@from_user_id, @to_user_id, @message
		)`

	_, err := r.client.Exec(ctx, q, args)
	if err != nil {
		r.logger.Errorln("error:", err)
		return err
	}

	return nil
}
