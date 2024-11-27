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
	const op = "repository.StoreUser"

	var user_id int64
	err := r.db.QueryRow(`
		INSERT INTO users (
			name,
			phone,
			password,
			app_id
		) VALUES ($1, $2, $3, $4)
		RETURNING id;
	`, name, phone, password, appId).Scan(&user_id)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return user_id, nil
}

func (r *AuthRepository) GetUserByPhone(ctx context.Context, phone string, app_id int32) (model.User, error) {
	const op = "repository.GetUserByPhone"
	var user model.User

	err := r.db.QueryRow(`
		SELECT id, name, password, phone, app_id
		FROM users
		WHERE phone = $1 AND app_id = $2
	`, phone, app_id).Scan(&user.Id, &user.Name, &user.PassHash, &user.Phone, &user.AppId)
	if err != nil {
		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (r *AuthRepository) UpdatePassword(ctx context.Context, phone string, app_id int32, password string) (model.User, error) {
	const op = "repository.GetUserByPhone"
	var user model.User

	_, err := r.db.Exec(`
			UPDATE users
			SET password = $1
			WHERE phone = $2 AND app_id = $3
		`, password, phone, app_id)
	if err != nil {
		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
