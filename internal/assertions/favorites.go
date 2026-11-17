package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/favorites/models"
)

// ContainsFavoriteByAdID проверяет, что в избранном есть объявление с указанным id
func ContainsFavoriteByAdID(t *testing.T, items []models.FavoriteItem, adID string) {
	t.Helper()
	require.Truef(t, hasFavorite(items, adID), "Объявление %s должно быть в избранном", adID)
}

// NotContainsFavoriteByAdID проверяет, что в избранном нет объявления с указанным id
func NotContainsFavoriteByAdID(t *testing.T, items []models.FavoriteItem, adID string) {
	t.Helper()
	require.Falsef(t, hasFavorite(items, adID), "Объявление %s не должно быть в избранном", adID)
}

// hasFavorite возвращает true, если в списке избранного есть объявление с указанным id
func hasFavorite(items []models.FavoriteItem, adID string) bool {
	for _, item := range items {
		if item.Advertisement.ID == adID {
			return true
		}
	}
	return false
}
