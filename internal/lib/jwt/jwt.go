package jwt

import (
	"time"

	"github.com/ei-jobs/auth-service/internal/domain/model"
	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user *model.User, app *model.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.Id
	claims["phone"] = user.Phone
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.Id

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
