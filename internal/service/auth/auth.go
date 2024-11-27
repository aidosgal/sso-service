package service

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/ei-jobs/auth-service/internal/domain/model"
	"github.com/ei-jobs/auth-service/internal/lib/jwt"
	ssov1 "github.com/ei-jobs/protos/gen/go/sso"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	StoreUser(ctx context.Context, phone string, name string, appId int32, password []byte) (int64, error)
    UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
    DeleteUser(ctx context.Context, user_id int64) (bool, error)
	GetUserByPhone(ctx context.Context, phone string, app_id int32) (model.User, error)
	UpdatePassword(ctx context.Context, phone string, app_id int32, password []byte) (model.User, error)
    GetAppById(ctx context.Context, app_id int32) (model.App, error)
    GetUserById(ctx context.Context, user_id int64) (*model.User, error)
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
        userChan <- user
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
    const op = "authservice.ForgetPassword"

    //ToDo: implement to logic sending the sms code and receiving it

    return s.ChangePassword(ctx, phone, password, app_id)
}

func (s *AuthService) UpdateUser(ctx context.Context, user *ssov1.User) (*ssov1.User, error) {
    userRequest := &model.User{
        Id:     user.GetId(),
        Name:   user.GetName(),
        Phone:  user.GetPhone(),
        AppId:  user.GetAppId(),
        AvatarUrl: stringToPointer(user.GetAvatarUrl()),
        Description: stringToPointer(user.GetDescription()),
    }

    newUser, err := s.repository.UpdateUser(ctx, userRequest)
    if err != nil {
        return nil, err
    }

    return &ssov1.User{
        Id:     newUser.Id,
        Name:   newUser.Name,
        Phone:  newUser.Phone,
        AppId:  newUser.AppId,
        Description: ifNilReturnEmpty(newUser.Description),
        AvatarUrl: ifNilReturnEmpty(newUser.AvatarUrl),
    }, nil
}

func stringToPointer(s string) *string {
    if s == "" {
        return nil 
    }
    return &s
}

func (s *AuthService) GetUser(ctx context.Context, user_id int64) (*ssov1.User, error) {
    user, err := s.repository.GetUserById(ctx, user_id)
    if err != nil {
        return nil, err
    }
    
    return &ssov1.User{
        Id:          user.Id,
        Name:        user.Name,
        Phone:       user.Phone,
        AppId:       user.AppId,
        Balance:    int64(user.Balance),
        Description: ifNilReturnEmpty(user.Description),
        AvatarUrl:   ifNilReturnEmpty(user.AvatarUrl),
    }, nil
}

func ifNilReturnEmpty(s *string) string {
    if s == nil {
        return ""
    }
    return *s
}


func (s *AuthService) DeleteUser(ctx context.Context, user_id int64) (bool, error) {
    return s.repository.DeleteUser(ctx, user_id)
}

func (s *AuthService) ChangePassword(ctx context.Context, phone string, password string, app_id int32) (string, error) {
    const op = "authservice.ChangePassword"

    var wg sync.WaitGroup
    errChan := make(chan error, 2)
    userChan := make(chan *model.User, 1)
    appChan := make(chan *model.App, 1)

    wg.Add(1) 
    go func() {
        defer wg.Done()
        passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
            errChan <- err
            return
        }
        user, err := s.repository.UpdatePassword(ctx, phone, app_id, passHash)
        if err != nil {
            errChan <- err
            return
        }
        userChan <- &user
    }()

    wg.Add(1)
    go func() {
        defer wg.Done()

        app, err := s.repository.GetAppById(ctx, app_id)
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
                return "", fmt.Errorf("%s: %w", op, err)
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
    if err != nil {
        return "", fmt.Errorf("%s: %w", op, err)
    }

	return token, nil
}
