package service

import (
	"context"
	"errors"
	"server/internal/models"
	"server/internal/repository/storage"
	"server/pkg/jwt"
	package_log "server/pkg/logging"
)

type UsersService struct {
	logger *package_log.Logger
	repo   storage.UsersRepository
}

func NewUsersService(logger *package_log.Logger, repo storage.UsersRepository) *UsersService {
	return &UsersService{
		logger: logger,
		repo:   repo,
	}
}

func (s *UsersService) CreateUser(ctx context.Context, dto models.Registration) (models.AuthResponse, error) {
	userId, err := s.repo.CreateUser(ctx, dto)
	if err != nil {
		s.logger.Errorln("error", err)
		return models.AuthResponse{}, err
	}

	token, err := jwt.GenerateToken(userId)
	if err != nil {
		s.logger.Errorln("error", err)
		return models.AuthResponse{}, err
	}

	body := models.AuthResponse{
		Token: token,
		Id:    userId,
	}
	return body, nil
}

func (s *UsersService) CheckUser(ctx context.Context, dto models.Registration) (models.AuthResponse, error) {
	userId, err := s.repo.CheckUser(ctx, dto)
	if err != nil {
		return models.AuthResponse{}, errors.New("unauthorizated")
	}

	if userId != "" {
		token, err := jwt.GenerateToken(userId)
		if err != nil {
			s.logger.Errorln("error", err)
			return models.AuthResponse{}, err
		}

		body := models.AuthResponse{
			Token: token,
			Id:    userId,
		}

		return body, nil
	}

	return models.AuthResponse{}, errors.New("unauthorizated")
}
