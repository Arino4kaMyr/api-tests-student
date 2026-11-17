package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/helpers/response"
	commonModels "api-tests-template/internal/managers/common/models"
	"api-tests-template/internal/managers/favorites/models"
)

// AddFavorite добавляет объявление в избранное пользователя, ожидает 201
func AddFavorite(t *testing.T, token, adID string) models.FavoriteItem {
	t.Helper()
	raw := favorites.HttpPostFavorite(t, token, models.AddFavoriteRequest{AdvertisementID: adID})
	return response.DecodeJSON[models.FavoriteItem](t, raw, http.StatusCreated)
}

// AddFavoriteExpectError пытается добавить объявление в избранное и возвращает ошибку API с ожидаемым статус-кодом
func AddFavoriteExpectError(t *testing.T, token, adID string, expectedStatusCode int) commonModels.ErrorResponse {
	t.Helper()
	raw := favorites.HttpPostFavorite(t, token, models.AddFavoriteRequest{AdvertisementID: adID})
	return response.DecodeError(t, raw, expectedStatusCode)
}

// GetFavorites возвращает список избранного пользователя, ожидает 200
func GetFavorites(t *testing.T, token string) models.FavoritesList {
	t.Helper()
	raw := favorites.HttpGetFavorites(t, token)
	return response.DecodeJSON[models.FavoritesList](t, raw, http.StatusOK)
}

// GetFavoritesExpectError пытается получить список избранного и возвращает ошибку API с ожидаемым статус-кодом
func GetFavoritesExpectError(t *testing.T, token string, expectedStatusCode int) commonModels.ErrorResponse {
	t.Helper()
	raw := favorites.HttpGetFavorites(t, token)
	return response.DecodeError(t, raw, expectedStatusCode)
}

// RemoveFavorite удаляет объявление из избранного пользователя, ожидает 204
func RemoveFavorite(t *testing.T, token, adID string) {
	t.Helper()
	raw := favorites.HttpDeleteFavorite(t, token, adID)
	response.ExpectNoContent(t, raw, http.StatusNoContent)
}

// RemoveFavoriteExpectError пытается удалить объявление из избранного и возвращает ошибку API с ожидаемым статус-кодом
func RemoveFavoriteExpectError(t *testing.T, token, adID string, expectedStatusCode int) commonModels.ErrorResponse {
	t.Helper()
	raw := favorites.HttpDeleteFavorite(t, token, adID)
	return response.DecodeError(t, raw, expectedStatusCode)
}
