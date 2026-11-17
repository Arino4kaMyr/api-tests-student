package advertisements

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/helpers/multipart"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/advertisements/models"
	commonModels "api-tests-template/internal/managers/common/models"
	"api-tests-template/internal/utils"
)

const photoFileName = "photo.png"

// GetAdvertisementsList выполняет поиск объявлений по тексту и возвращает список, ожидает 200
func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	raw := advertisements.HttpGetAdvertisements(t, token, text)
	defer raw.Body.Close()
	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := models.AdvertisementsList{}
	err = json.Unmarshal(body, &resp)
	require.NoError(t, err)

	return resp
}

// CreateAdvertisement создаёт объявление и возвращает его, ожидает 201
func CreateAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) models.Advertisement {
	t.Helper()
	raw := advertisements.HttpPostAdvertisement(t, token, buildCreateForm(request))
	return response.DecodeJSON[models.Advertisement](t, raw, http.StatusCreated)
}

// CreateAdvertisementExpectError пытается создать объявление и возвращает ошибку API с ожидаемым статус-кодом
func CreateAdvertisementExpectError(t *testing.T, token string, request models.CreateAdvertisementRequest,
	expectedStatusCode int) commonModels.ErrorResponse {
	t.Helper()
	raw := advertisements.HttpPostAdvertisement(t, token, buildCreateForm(request))
	return response.DecodeError(t, raw, expectedStatusCode)
}

// GetAdvertisement возвращает объявление по id, ожидает 200
func GetAdvertisement(t *testing.T, id string) models.Advertisement {
	t.Helper()
	raw := advertisements.HttpGetAdvertisement(t, id)
	return response.DecodeJSON[models.Advertisement](t, raw, http.StatusOK)
}

// GetAdvertisementPhotos возвращает фотографии объявления, ожидает 200
func GetAdvertisementPhotos(t *testing.T, token, id string) []models.Photo {
	t.Helper()
	raw := advertisements.HttpGetAdvertisementPhotos(t, token, id)
	return response.DecodeJSON[[]models.Photo](t, raw, http.StatusOK)
}

// DeleteAdvertisement удаляет объявление, ожидает 204
func DeleteAdvertisement(t *testing.T, token, id string) {
	t.Helper()
	raw := advertisements.HttpDeleteAdvertisement(t, token, id)
	response.ExpectNoContent(t, raw, http.StatusNoContent)
}

// DeleteAdvertisementIfExists удаляет объявление, если оно ещё существует; для очистки данных после теста
func DeleteAdvertisementIfExists(t *testing.T, token, id string) {
	t.Helper()
	raw := advertisements.HttpGetAdvertisement(t, id)
	defer raw.Body.Close()

	if raw.StatusCode == http.StatusNotFound {
		return
	}
	DeleteAdvertisement(t, token, id)
}

// buildCreateForm собирает multipart-форму создания объявления из параметров запроса
func buildCreateForm(request models.CreateAdvertisementRequest) multipart.Form {
	builder := multipart.NewBuilder().
		Field("title", request.Title).
		Field("description", request.Description).
		Field("price", strconv.FormatFloat(request.Price, 'f', -1, 64)).
		Field("quantity", strconv.Itoa(request.Quantity))

	photo := multipart.GeneratePNG()
	for i := 0; i < request.PhotosCount; i++ {
		builder.File("photos", fmt.Sprintf("%d_%s", i, photoFileName), photo)
	}

	return builder.Build()
}

// NewAdvertisementRequest возвращает валидные параметры объявления со всеми полями и уникальным названием
func NewAdvertisementRequest() models.CreateAdvertisementRequest {
	return models.CreateAdvertisementRequest{
		Title:       "Автотест " + utils.RandomString(10),
		Description: "Описание тестового объявления " + utils.RandomString(10),
		Price:       1234.5,
		Quantity:    7,
		PhotosCount: 1,
	}
}
