package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ei-jobs/auth-service/internal/config"
	"github.com/ei-jobs/auth-service/internal/domain/model"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(cfg config.DatabaseConfig) (*AuthRepository, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/postgres?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.SSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}
	defer db.Close()

	connStr = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the newly created database: %w", err)
	}

	return &AuthRepository{db: db}, nil
}

func (r *AuthRepository) StoreUser(ctx context.Context, phone string, name string, appId int32, password string) (int64, error) {
	return 0, nil
}

func (r *AuthRepository) GetUserByPhone(ctx context.Context, phone string) (model.User, error) {
	var user model.User
	return user, nil
}
