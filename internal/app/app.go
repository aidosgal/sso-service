package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/ei-jobs/auth-service/internal/app/grpc"
	"github.com/ei-jobs/auth-service/internal/config"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, cfg config.DatabaseConfig, tokenTTL time.Duration) *App {
	grpcApp := grpcapp.NewApp(log, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
