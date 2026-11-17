package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements/models"
)

func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	t.Helper()
	for _, item := range items {
		if item.Title == title {
			return
		}
	}
	require.FailNow(t, "Объявление не найдено", "Название: %q", title)
}

// ContainsAdvertisementByID проверяет, что среди объявлений есть объявление с указанным id
func ContainsAdvertisementByID(t *testing.T, items []models.Advertisement, id string) {
	t.Helper()
	for _, item := range items {
		if item.ID == id {
			return
		}
	}
	require.FailNow(t, "Объявление не найдено", "ID: %q", id)
}

// AdvertisementMatchesRequest проверяет, что параметры объявления совпадают с переданными при создании
func AdvertisementMatchesRequest(t *testing.T, request models.CreateAdvertisementRequest, actual models.Advertisement) {
	t.Helper()
	require.Equal(t, request.Title, actual.Title, "title не совпадает с переданным при создании")
	require.Equal(t, request.Description, actual.Description, "description не совпадает с переданным при создании")
	require.InDelta(t, request.Price, actual.Price, 0.001, "price не совпадает с переданным при создании")
	require.Equal(t, request.Quantity, actual.Quantity, "quantity не совпадает с переданным при создании")
	require.Len(t, actual.Photos, request.PhotosCount, "Количество фотографий не совпадает с переданным при создании")
}
