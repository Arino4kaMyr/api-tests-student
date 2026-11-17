package advertisements

import (
	"fmt"
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/helpers/multipart"
)

// HttpGetAdvertisements выполняет GET-запрос поиска объявлений по тексту, ожидает 200
func HttpGetAdvertisements(t *testing.T, token, textForSearch string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Expect(t).Status(http.StatusOK).End().Response
}

// HttpPostAdvertisement выполняет POST-запрос создания объявления с фотографиями (multipart/form-data)
func HttpPostAdvertisement(t *testing.T, token string, form multipart.Form) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Post(path.AdvertisementPath).
		Header("Content-Type", form.ContentType).
		Body(form.Body).
		Expect(t).
		End().Response
}

// HttpGetAdvertisement выполняет GET-запрос получения объявления по id
func HttpGetAdvertisement(t *testing.T, id string) *http.Response {
	return apiRunner.GetRunner().Create().
		Get(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}

// HttpDeleteAdvertisement выполняет DELETE-запрос удаления объявления по id
func HttpDeleteAdvertisement(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}

// HttpGetAdvertisementPhotos выполняет GET-запрос получения фотографий объявления
func HttpGetAdvertisementPhotos(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(fmt.Sprintf(path.AdvertisementPhotosFmt, id)).
		Expect(t).
		End().Response
}
