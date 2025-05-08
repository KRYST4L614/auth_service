package tests

import (
	"github.com/KRYST4L614/auth_service/tests/suite"
	ssov1 "github.com/KRYST4L614/auth_service_protos/gen/go/sso"
	"github.com/brianvoe/gofakeit"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	passDefaultLen    = 10
	appId             = 1
	appSecret         = "test-secret"
	tokenDeltaSeconds = 10
)

func TestLogin_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, true, passDefaultLen)

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)
	require.NotEmpty(t, respReg.GetUserId())

	loginTime := time.Now()

	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appId,
	})
	require.NoError(t, err)

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appId, int(claims["app_id"].(float64)))
	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTl).Unix(), claims["exp"].(float64), tokenDeltaSeconds)
}
