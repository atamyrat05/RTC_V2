package repository

import (
	"context"
	"errors"
	"server/internal/models"
	package_psql "server/pkg/db"
	package_log "server/pkg/logging"

	"github.com/jackc/pgx/v5"
)

type StoryPsqlRepository struct {
	logger *package_log.Logger
	client package_psql.Client
}

func NewStoryPsqlRepository(logger *package_log.Logger, client package_psql.Client) *StoryPsqlRepository {
	return &StoryPsqlRepository{
		logger: logger,
		client: client,
	}
}

func (r *StoryPsqlRepository) SaveStory(ctx context.Context, id string, url string) error {
	q := `INSERT INTO stories (user_id, file_url) VALUES (@uuid, @file_url)`

	args := pgx.NamedArgs{
		"uuid":     id,
		"file_url": url,
	}

	_, err := r.client.Exec(ctx, q, args)
	if err != nil {
		r.logger.Errorln("error:", err)
		return err
	}

	return nil
}

func (r *StoryPsqlRepository) GetUserStory(ctx context.Context, id string) (models.Story, error) {
	var data models.Story

	args := pgx.NamedArgs{
		"uuid": id,
	}

	q := `
		SELECT username, COALESCE(user_image, '') as user_image 
		FROM users 
		WHERE uuid=@uuid
		`

	err := r.client.QueryRow(ctx, q, args).Scan(&data.Username, &data.Image_url)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		r.logger.Errorln("error:", err)
		return models.Story{}, err
	}

	q = `SELECT file_url FROM stories WHERE user_id=@uuid ORDER BY created_at DESC LIMIT 1`

	err = r.client.QueryRow(ctx, q, args).Scan(&data.File_url)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		r.logger.Errorln("error:", err)
		return models.Story{}, err
	}

	return data, nil
}

// func (r *StoryPsqlRepository) GetAllStoryForUser(ctx context.Context, dto models.GetStoryDto) ([]models.Story, error) {
// 	var data []models.Story

// 	q := `
// 		SELECT following_user_id
// 		FROM following
// 		WHERE user_id=@uuid
// 		LIMIT @limit
// 		OFFSET @offset
// 		`

// 	args := pgx.NamedArgs{
// 		"uuid":   dto.User_Id,
// 		"limit":  dto.Limit,
// 		"offset": dto.Offset,
// 	}

// 	rows, err := r.client.Query(ctx, q, args)
// 	if err != nil {
// 		r.logger.Errorln("error:", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var followingUserId string
// 		var body models.Story
// 		k := 0

// 		if err = rows.Scan(&followingUserId); err != nil {
// 			r.logger.Errorln("error:", err)
// 			return nil, err
// 		}

// 		args = pgx.NamedArgs{
// 			"uuid": followingUserId,
// 		}

// 		q = `SELECT file_url FROM stories WHERE user_id=@uuid ORDER BY created_at DESC LIMIT 1`

// 		err = r.client.QueryRow(ctx, q, args).Scan(&body.File_url)
// 		if err != nil {
// 			if errors.Is(err, pgx.ErrNoRows) {
// 				k = 1
// 			} else {
// 				r.logger.Errorln("error", err)
// 				return nil, err
// 			}
// 		}

// 		if k == 0 {
// 			q = `SELECT username, user_image FROM users WHERE uuid=@uuid`

// 			err = r.client.QueryRow(ctx, q, args).Scan(&body.Username, &body.Image_url)
// 			if err != nil {
// 				r.logger.Errorln("error:", err)
// 				return nil, err
// 			}

// 			body.User_id = followingUserId

// 			data = append(data, body)
// 		}
// 	}

// 	return data, nil
// }

func (r *StoryPsqlRepository) GetStoryById(ctx context.Context, id string) ([]string, error) {
	var files []string

	q := `SELECT file_url FROM stories WHERE user_id=@uuid ORDER BY created_at DESC`

	args := pgx.NamedArgs{
		"uuid": id,
	}

	rows, err := r.client.Query(ctx, q, args)
	if err != nil {
		r.logger.Errorln("error:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var file string
		if err = rows.Scan(&file); err != nil {
			r.logger.Errorln("error:", err)
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

func (r *StoryPsqlRepository) GetAllStoryForUser(ctx context.Context, dto models.GetStoryDto) ([]models.Story, error) {
	var body []models.Story
	data, err := r.GetUserChatHistoriesIds(ctx, dto.User_Id)
	if err != nil {
		r.logger.Errorln("error", err)
		return nil, err
	}

	for _, userId := range data {
		var row models.Story
		args := pgx.NamedArgs{
			"id": userId,
		}

		q := `
			SELECT 
				u.uuid AS user_id,
				u.username,
				u.user_image AS image_url,
				COALESCE(s.file_url, '') AS file_url
			FROM users u
			LEFT JOIN stories s
				ON u.uuid = s.user_id
			WHERE u.uuid = @id
			ORDER BY s.created_at DESC
			LIMIT 1
			`

		err = r.client.QueryRow(ctx, q, args).Scan(&row.User_id, &row.Username, &row.Image_url, &row.File_url)
		if err != nil {
			r.logger.Errorln("error", err)
			return nil, err
		}

		body = append(body, row)
	}

	return body, nil
}

func (r *StoryPsqlRepository) GetUserChatHistoriesIds(ctx context.Context, user_id string) ([]string, error) {
	var data []string

	//usering yazan adamlary
	q := `
		SELECT to_user_id 
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
		var row string

		if err = rows.Scan(&row); err != nil && !errors.Is(err, pgx.ErrNoRows) {
			r.logger.Errorln("error:", err)
			return nil, err
		}

		if k == 1 {
			data = append(data, row)
			count++
		} else {
			var w = 0
			for i := range data {
				if data[i] == row {
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
	}

	//basga user lar yazan bolsa
	q = `
		SELECT from_user_id
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
		var row string

		if err = rows.Scan(&row); err != nil && !errors.Is(err, pgx.ErrNoRows) {
			r.logger.Errorln("error:", err)
			return nil, err
		}

		var w = 0
		for i := range data {
			if data[i] == row {
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

	return data, nil
}
