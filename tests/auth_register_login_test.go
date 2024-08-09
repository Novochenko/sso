package tests

import (
	"testing"
	"time"

	"github.com/Novochenko/protos/gen/go/sso"
	"github.com/Novochenko/sso/tests/suite"
	"github.com/brianvoe/gofakeit"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppID = ""
	appID      = "70ccbfb5-5116-4fdd-bc05-8a07adf40193"
	appSecret  = "test-secret"

	passDefaultLen = 10
)

func TestRegisteLogin_login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)
	email := gofakeit.Email()
	password := randomFakePassword()
	responseReg, err := st.AuthClient.Register(ctx, &sso.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, responseReg)

	respLogin, err := st.AuthClient.Login(ctx, &sso.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	token := respLogin.Token
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, responseReg.GetUserId(), claims["uid"].(string))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, claims["app_id"].(string))

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}
func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
