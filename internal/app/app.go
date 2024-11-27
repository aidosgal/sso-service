package app

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	grpcapp "github.com/ei-jobs/auth-service/internal/app/grpc"
	"github.com/ei-jobs/auth-service/internal/config"
	repository "github.com/ei-jobs/auth-service/internal/repository/auth"
	service "github.com/ei-jobs/auth-service/internal/service/auth"
	_ "github.com/lib/pq"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, cfg config.DatabaseConfig, tokenTTL time.Duration) *App {
    connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/postgres?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.SSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
        panic(err)
	}
	defer db.Close()

	connStr = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)
	db, err = sql.Open("postgres", connStr)

	authRepository := repository.NewAuthRepository(db)

	authService := service.NewAuthService(log, authRepository)

	grpcApp := grpcapp.NewApp(log, grpcPort, authService)

	return &App{
		GRPCSrv: grpcApp,
	}
}
