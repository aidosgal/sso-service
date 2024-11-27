package service

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/ei-jobs/auth-service/internal/domain/model"
	"github.com/ei-jobs/auth-service/internal/lib/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	StoreUser(ctx context.Context, phone string, name string, appId int32, password []byte) (int64, error)
	GetUserByPhone(ctx context.Context, phone string, app_id int32) (model.User, error)
	UpdatePassword(ctx context.Context, phone string, app_id int32, password []byte) (model.User, error)
    GetAppById(ctx context.Context, app_id int32) (model.App, error)
    GetUserById(ctx context.Context, user_id int64) (model.User, error)
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
    const op = "authservice.Login"

    user, err := s.repository.GetUserByPhone(ctx, phone, appId)
    if err != nil {
        return "", fmt.Errorf("%s: %w", op, err) 
    }

    if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		s.log.Info("invalid credentials", slog.String("error", err.Error())) 

		return "", fmt.Errorf("%s: %s", op, "Invalid credentials")
	}

    app, err := s.repository.GetAppById(ctx, appId)
    if err != nil {
        return "", fmt.Errorf("%s: %w", op, err)
    }

    token, err := jwt.NewToken(&user, &app, time.Hour * 24 * 365)
	return token, nil
}

func (s *AuthService) Register(ctx context.Context, name string, phone string, password string, appId int32) (string, error) {
    const op = "authservice.Regsiter"

    passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) 
    if err != nil {
        return "", fmt.Errorf("%s: %w", op, err)
    }

    user_id, err := s.repository.StoreUser(ctx, phone, name, appId, passHash)
    if err != nil {
        return "", fmt.Errorf("%s: %w", op, err)
    } 
    
    var wg sync.WaitGroup
    errChan := make(chan error, 2)
    userChan := make(chan *model.User, 1)
    appChan := make(chan *model.App, 1)
    
    wg.Add(1)
    go func() {
        defer wg.Done()
        user, err := s.repository.GetUserById(ctx, user_id)
        if err != nil {
            errChan <- err
            return
        }
        userChan <- &user
    }()

    wg.Add(1)
    go func() {
        defer wg.Done()
        app, err := s.repository.GetAppById(ctx, appId)
        if err != nil {
            errChan <- err
            return
        }
        appChan <- &app
    }()

    go func() {
        wg.Wait()
        close(errChan)
        close(userChan)
        close(appChan)
    }()

    var user *model.User
    var app *model.App
    for {
        select {
        case err := <-errChan:
            if err != nil {
                return "", err
            }
        case u := <-userChan:
            if u != nil {
                user = u
            }
        case a := <-appChan:
            if a != nil {
                app = a
            }
        case <-ctx.Done():
            return "", fmt.Errorf("%s: context cancelled", op)
        }

        if user != nil && app != nil {
            break
        }
    }

    token, err := jwt.NewToken(user, app, time.Hour * 24 * 365)

	return token, nil
}

func (s *AuthService) ForgetPassword(ctx context.Context, phone string, password string, app_id int32) (string, error) {
	return "", nil
}

func (s *AuthService) ChangePassword(ctx context.Context, phone string, password string, app_id int32) (string, error) {
	return "", nil
}
