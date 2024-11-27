package service

import (
	"log/slog"
)

type UserRespository interface {
}

type UserService struct {
    log *slog.Logger
    repository UserRespository 
}

func NewUserService(log *slog.Logger, repository UserRespository) *UserService {
    return &UserService{
        log: log,
        repository: repository,
    }
}
