package service

import (
	"context"
	"log/slog"

	"github.com/ei-jobs/auth-service/internal/domain/model"
)

type AuthRepository interface {
	StoreUser(ctx context.Context, phone string, name string, appId int32, password string) (int64, error)
	GetUserByPhone(ctx context.Context, phone string, app_id int32) (model.User, error)
	UpdatePassword(ctx context.Context, phone string, app_id int32, password string) (model.User, error)
}

type AuthService struct {
	log        *slog.Logger
	repository AuthRepository
}

func NewAuthService(log *slog.Logger, repository AuthRepository) *AuthService {
	return &AuthService{
		log:        log,
		repository: repository,
	}
}

func (s *AuthService) Login(ctx context.Context, phone string, password string, appId int32) (string, error) {
	return "", nil
}

func (s *AuthService) Register(ctx context.Context, name string, phone string, password string, appId int32) (string, error) {
	return "", nil
}

func (s *AuthService) ForgetPassword(ctx context.Context, phone string, password string, app_id int32) (string, error) {
	return "", nil
}

func (s *AuthService) ChangePassword(ctx context.Context, phone string, password string, app_id int32) (string, error) {
	return "", nil
}
