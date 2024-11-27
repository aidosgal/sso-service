package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/ei-jobs/auth-service/internal/app/grpc"
	"github.com/ei-jobs/auth-service/internal/config"
	repository "github.com/ei-jobs/auth-service/internal/repository/auth"
	service "github.com/ei-jobs/auth-service/internal/service/auth"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, cfg config.DatabaseConfig, tokenTTL time.Duration) *App {
	authRepository, err := repository.NewAuthRepository(cfg)
	if err != nil {
		panic(err)
	}

	authService := service.NewAuthService(log, authRepository)

	grpcApp := grpcapp.NewApp(log, grpcPort, authService)

	return &App{
		GRPCSrv: grpcApp,
	}
}
