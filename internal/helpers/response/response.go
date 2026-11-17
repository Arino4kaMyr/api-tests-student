package response

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/common/models"
)

// DecodeJSON проверяет статус-код ответа и разбирает тело ответа в структуру T
func DecodeJSON[T any](t *testing.T, raw *http.Response, expectedStatusCode int) T {
	t.Helper()
	body := readBody(t, raw, expectedStatusCode)

	var result T
	require.NoErrorf(t, json.Unmarshal(body, &result), "Не удалось разобрать тело ответа: %s", body)

	return result
}

// DecodeError проверяет статус-код ответа и разбирает тело ответа с ошибкой
func DecodeError(t *testing.T, raw *http.Response, expectedStatusCode int) models.ErrorResponse {
	t.Helper()
	return DecodeJSON[models.ErrorResponse](t, raw, expectedStatusCode)
}

// ExpectNoContent проверяет статус-код ответа без тела
func ExpectNoContent(t *testing.T, raw *http.Response, expectedStatusCode int) {
	t.Helper()
	readBody(t, raw, expectedStatusCode)
}

// readBody проверяет статус-код, закрывает тело ответа и возвращает его содержимое
func readBody(t *testing.T, raw *http.Response, expectedStatusCode int) []byte {
	t.Helper()
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Не удалось прочитать тело ответа")
	require.Equalf(t, expectedStatusCode, raw.StatusCode,
		"HTTP status code должен быть %d, тело ответа: %s", expectedStatusCode, body)

	return body
}
