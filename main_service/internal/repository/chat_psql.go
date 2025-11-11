package repository

import (
	"context"
	"errors"
	"server/internal/models"
	package_psql "server/pkg/db"
	package_log "server/pkg/logging"

	"github.com/jackc/pgx/v5"
)

type ChatPsqlRepository struct {
	logger *package_log.Logger
	client package_psql.Client
}

func NewChatPsqlRepository(logger *package_log.Logger, client package_psql.Client) *ChatPsqlRepository {
	return &ChatPsqlRepository{
		logger: logger,
		client: client,
	}
}

func (r *ChatPsqlRepository) GetUserChatHistroies(ctx context.Context, user_id string) ([]models.Chat, error) {
	var data []models.Chat

	//usering yazan adamlary
	q := `
		SELECT to_user_id, created_at 
		FROM messages WHERE from_user_id=@uuid 
		ORDER BY created_at DESC
		`

	args := pgx.NamedArgs{
		"uuid": user_id,
	}

	rows, err := r.client.Query(ctx, q, args)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		r.logger.Errorln("error:", err)
		return nil, err
	}
	defer rows.Close()

	var k = 1
	var count = -1
	for rows.Next() {
		var row models.Chat

		if err = rows.Scan(&row.User_Id, &row.LastTime); err != nil && !errors.Is(err, pgx.ErrNoRows) {
			r.logger.Errorln("error:", err)
			return nil, err
		}

		if k == 1 {
			data = append(data, row)
			count++
		} else {
			var w = 0
			for _, j := range data {
				if j.User_Id == row.User_Id {
					w = 1
				}
				if w == 1 {
					break
				}
			}
			if w == 0 {
				data = append(data, row)
				count++
			}
		}
		k = 0

		if data[count].Username == "" {
			args = pgx.NamedArgs{
				"user1": user_id,
				"uuid":  data[count].User_Id,
			}

			q = `SELECT username, user_image FROM users WHERE uuid=@uuid`

			err = r.client.QueryRow(ctx, q, args).Scan(&data[count].Username, &data[count].Image_url)
			if err != nil && !errors.Is(err, pgx.ErrNoRows) {
				r.logger.Errorln("error:", err)
				return nil, err
			}

			q = `
				SELECT uuid 
				FROM rooms 
				WHERE (user_1=@user1 AND user_2=@uuid) 
				OR (user_1=@uuid AND user_2=@user1)
			`

			err = r.client.QueryRow(ctx, q, args).Scan(&data[count].RoomId)
			if err != nil && !errors.Is(err, pgx.ErrNoRows) {
				r.logger.Errorln("error:", err)
				return nil, err
			}
		}
	}

	//basga user lar yazan bolsa
	q = `
		SELECT from_user_id, created_at 
		FROM messages 
		WHERE to_user_id=@uuid 
		ORDER BY created_at DESC
		`

	args = pgx.NamedArgs{
		"uuid": user_id,
	}

	rows, err = r.client.Query(ctx, q, args)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		r.logger.Errorln("error:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var row models.Chat

		if err = rows.Scan(&row.User_Id, &row.LastTime); err != nil && !errors.Is(err, pgx.ErrNoRows) {
			r.logger.Errorln("error:", err)
			return nil, err
		}

		var w = 0
		for _, j := range data {
			if j.User_Id == row.User_Id {
				w = 1
			}
			if w == 1 {
				break
			}
		}
		if w == 0 {
			data = append(data, row)
			count++
		}

		if data[count].Username == "" {
			args = pgx.NamedArgs{
				"user1": user_id,
				"uuid":  data[count].User_Id,
			}

			q = `SELECT username, user_image FROM users WHERE uuid=@uuid`

			err = r.client.QueryRow(ctx, q, args).Scan(&data[count].Username, &data[count].Image_url)
			if err != nil && !errors.Is(err, pgx.ErrNoRows) {
				r.logger.Errorln("error:", err)
				return nil, err
			}

			q = `
				SELECT uuid 
				FROM rooms 
				WHERE (user_1=@user1 AND user_2=@uuid) 
				OR (user_1=@uuid AND user_2=@user1)
			`

			err = r.client.QueryRow(ctx, q, args).Scan(&data[count].RoomId)
			if err != nil && !errors.Is(err, pgx.ErrNoRows) {
				r.logger.Errorln("error:", err)
				return nil, err
			}
		}
	}

	return data, nil
}

func (r *ChatPsqlRepository) GetUserChats(ctx context.Context, dto models.ChatWithUser) ([]models.ChatsMessage, error) {
	var data []models.ChatsMessage

	//user poluchayut sms history
	q := `
		SELECT from_user_id, content, created_at 
		FROM messages 
		WHERE (from_user_id=@user_id AND to_user_id=@partner_id) 
			OR (from_user_id=@partner_id AND to_user_id=@user_id)
		ORDER BY created_at
	`

	args := pgx.NamedArgs{
		"user_id":    dto.UserID,
		"partner_id": dto.PartnerID,
		"limit":      dto.Limit,
		"offset":     dto.Offset,
	}

	rows, err := r.client.Query(ctx, q, args)
	if err != nil {
		r.logger.Errorln("error:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var body models.ChatsMessage

		if err = rows.Scan(&body.Sender, &body.Message, &body.Time); err != nil {
			r.logger.Errorln("error:", err)
			return nil, err
		}

		data = append(data, body)
	}

	return data, nil
}

func (r *ChatPsqlRepository) SaveChats(ctx context.Context, dto models.SaveChat) error {
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
		r.logger.Errorln("error", err)
		return err
	}

	return nil
}
