package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ei-jobs/auth-service/internal/domain/model"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) StoreUser(ctx context.Context, phone string, name string, appId int32, password []byte) (int64, error) {
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

func (r *AuthRepository) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	const op = "repository.UpdateUser"

	_, err := r.db.ExecContext(ctx, `
		UPDATE users
		SET name = $1, avatar_url = $2, description = $3
		WHERE id = $4
	`, user.Name, user.AvatarUrl, user.Description, user.Id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (r *AuthRepository) GetUserById(ctx context.Context, user_id int64) (*model.User, error) {
	const op = "repository.GetUser"
	var user model.User

	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, password, phone, app_id, avatar_url, description, balance
		FROM users
		WHERE id = $1
	`, user_id).Scan(&user.Id, &user.Name, &user.PassHash, &user.Phone, &user.AppId, &user.AvatarUrl, &user.Description, &user.Balance)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (r *AuthRepository) DeleteUser(ctx context.Context, user_id int64) (bool, error) {
	const op = "repository.DeleteUser"

	result, err := r.db.ExecContext(ctx, `
		DELETE FROM users
		WHERE id = $1
	`, user_id)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return rowsAffected > 0, nil
}

func (r *AuthRepository) UpdatePassword(ctx context.Context, phone string, app_id int32, password []byte) (model.User, error) {
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

func (r *AuthRepository) GetAppById(ctx context.Context, app_id int32) (model.App, error) {
    const op = "repository.GetAppById"
    var app model.App

    err := r.db.QueryRow(`
        SELECT id, name, secret 
        FROM apps
        WHERE id = $1
    `, app_id).Scan(&app.Id, &app.Name, &app.Secret)
    if err != nil {
        return app, fmt.Errorf("%s: %w", op, err)
    }

    return app, nil
} 

