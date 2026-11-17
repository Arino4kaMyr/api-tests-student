package auth

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/auth"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
)

// Login выполняет авторизацию пользователя и возвращает токен с данными пользователя, ожидает 200
func Login(t *testing.T, username, password string) models.LoginOkResponse {
	require.NotEmpty(t, username, "Отсутствует логин пользователя")
	require.NotEmpty(t, password, "Отсутствует пароль пользователя")

	request := models.LoginRequest{
		Email:    username,
		Password: password,
	}

	loginRawResponse := auth.HttpPostAuthLogin(t, request)
	defer loginRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(loginRawResponse.Body)
	require.NoError(t, errResponse)

	loginResponse := models.LoginOkResponse{}
	err := json.Unmarshal(respArray, &loginResponse)
	require.NoError(t, err)

	return loginResponse
}

// passwordPrefix гарантирует наличие цифры и заглавной буквы в пароле, которых требует API
const passwordPrefix = "Aa1"

// RegisterRandomUser регистрирует нового пользователя со случайным email и возвращает токен с его данными, ожидает 201
func RegisterRandomUser(t *testing.T) models.LoginOkResponse {
	t.Helper()
	request := models.RegisterRequest{
		Email:     "autotest_" + strings.ToLower(utils.RandomString(12)) + "@example.com",
		FirstName: "Autotest",
		Password:  passwordPrefix + utils.RandomString(16),
	}

	raw := auth.HttpPostAuthRegister(t, request)
	return response.DecodeJSON[models.LoginOkResponse](t, raw, http.StatusCreated)
}
