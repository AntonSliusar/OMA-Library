package service

import (
	"log/slog"
	"oma-library/internal/models"
	"oma-library/pkg/storage"
)

type UserService struct {
	storage *storage.Storage
}

func NewUserService(s *storage.Storage) *UserService {
	return &UserService{storage: s}
}

func(s UserService) CheckExist(email string) (bool, error) {
	exist, err := s.storage.CheckExist(email)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func(s UserService) RegisterUser(req models.SignUpRequest) error {
	err := s.RegisterUser(req)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func(s UserService) LoginUser(req models.SignInRequset) (models.User, error) {
	user, err := s.storage.GetByEmail(req.Email)
	if err != nil {
		slog.Error(err.Error())
		return models.User{}, err
	}
	return user, nil
}
