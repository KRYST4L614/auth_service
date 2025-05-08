package jwt

import (
	"github.com/KRYST4L614/auth_service/internal/domain/entity"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(user entity.User, app entity.App, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":    user.ID,
		"email":  user.Email,
		"exp":    time.Now().Add(duration).Unix(),
		"app_id": app.ID,
	})

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
