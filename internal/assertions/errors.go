package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/common/models"
)

// ErrorEquals проверяет код и сообщение ошибки в ответе API
func ErrorEquals(t *testing.T, actual models.ErrorResponse, expectedError, expectedMessage string) {
	t.Helper()
	require.Equal(t, expectedError, actual.Error, "Некорректное значение error параметра")
	require.Equal(t, expectedMessage, actual.Message, "Некорректное значение message параметра")
}
