package createAdvertisement

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/advertisement"
	"api-tests-template/internal/constants/apierrors"
	ads "api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"

	base "api-tests-template/tests"
)

type TestSuite struct {
	suite.Suite
	loginData authModels.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Регистрация пользователя для тестов создания объявления")
	s.loginData = auth.RegisterRandomUser(s.T())
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// Positive: Добавление объявления со всеми параметрами — оно доступно по id, находится поиском и имеет загруженное фото
func (s *TestSuite) TestCreateAdvertisementPositive() {
	request := ads.NewAdvertisementRequest()

	base.Precondition(fmt.Sprintf("Объявление %q будет удалено после теста", request.Title))
	var created models.Advertisement
	t := s.T() // cleanup регистрируется на тесте целиком, а не на подтесте, иначе объявление удалится сразу после создания
	t.Cleanup(func() {
		if created.ID != "" {
			ads.DeleteAdvertisementIfExists(t, s.loginData.Token, created.ID)
		}
	})
	s.Run("Создаём объявление со всеми параметрами", func() {
		base.Step("POST /advertisement со всеми параметрами и одной фотографией")
		created = ads.CreateAdvertisement(s.T(), s.loginData.Token, request)

		base.Check("В ответе возвращены переданные параметры и автор объявления")
		require.NotEmpty(s.T(), created.ID, "id созданного объявления пустой")
		require.Equal(s.T(), s.loginData.User.Id, created.UserID, "Автор объявления не совпадает с авторизованным пользователем")
		assertions.AdvertisementMatchesRequest(s.T(), request, created)
	})

	s.Run("Объявление доступно по id с теми же параметрами", func() {
		base.Step("GET /advertisement?id=" + created.ID)
		actual := ads.GetAdvertisement(s.T(), created.ID)

		base.Check("Параметры объявления совпадают с переданными при создании")
		require.Equal(s.T(), created.ID, actual.ID, "id объявления не совпадает")
		require.Equal(s.T(), s.loginData.User.Id, actual.UserID, "user_id не совпадает")
		assertions.AdvertisementMatchesRequest(s.T(), request, actual)
	})

	s.Run("Объявление находится поиском по полному названию", func() {
		base.Step("GET /advertisements?search=" + request.Title)
		list := ads.GetAdvertisementsList(s.T(), s.loginData.Token, request.Title)

		base.Check("В результатах поиска есть созданное объявление")
		assertions.ContainsAdvertisementByID(s.T(), list.Items, created.ID)
	})

	s.Run("Фотография загружена на сервер и доступна", func() {
		base.Step("GET /advertisements/" + created.ID + "/photos")
		photos := ads.GetAdvertisementPhotos(s.T(), s.loginData.Token, created.ID)

		base.Check("Возвращена ровно одна фотография, совпадающая с фотографией из ответа на создание")
		require.Len(s.T(), photos, request.PhotosCount, "Количество загруженных фотографий не совпадает")
		require.Equal(s.T(), created.Photos[0].ID, photos[0].ID, "id фотографии не совпадает")
		require.NotEmpty(s.T(), photos[0].URL, "url фотографии пустой")
		require.Equal(s.T(), created.Photos[0].URL, photos[0].URL, "url фотографии не совпадает")
	})
}

// Negative: Создание объявления с некорректными параметрами возвращает 400 и понятное сообщение об ошибке
func (s *TestSuite) TestCreateAdvertisementInvalidParams() {
	cases := []struct {
		name            string
		modify          func(request *models.CreateAdvertisementRequest)
		expectedMessage string
	}{
		{"пустой title", func(r *models.CreateAdvertisementRequest) { r.Title = "" },
			"title is required"},
		{"слишком длинный title", func(r *models.CreateAdvertisementRequest) {
			r.Title = strings.Repeat("a", advertisement.MaxTitleLength+1)
		}, fmt.Sprintf("title must not exceed %d characters provided %d",
			advertisement.MaxTitleLength, advertisement.MaxTitleLength+1)},
		{"пустой description", func(r *models.CreateAdvertisementRequest) { r.Description = "" },
			"description is required"},
		{"слишком длинный description", func(r *models.CreateAdvertisementRequest) {
			r.Description = strings.Repeat("a", advertisement.MaxDescriptionLength+1)
		}, fmt.Sprintf("description must be between 1 and %d characters", advertisement.MaxDescriptionLength)},
		{"отрицательная цена", func(r *models.CreateAdvertisementRequest) { r.Price = -1 },
			fmt.Sprintf("price must be between 0 and %d inclusive", advertisement.MaxPrice)},
		{"цена выше максимальной", func(r *models.CreateAdvertisementRequest) { r.Price = advertisement.MaxPrice + 1 },
			fmt.Sprintf("price must be between 0 and %d inclusive", advertisement.MaxPrice)},
		{"количество меньше минимального", func(r *models.CreateAdvertisementRequest) { r.Quantity = advertisement.MinQuantity - 1 },
			fmt.Sprintf("quantity must be between %d and %d inclusive", advertisement.MinQuantity, advertisement.MaxQuantity)},
		{"количество больше максимального", func(r *models.CreateAdvertisementRequest) { r.Quantity = advertisement.MaxQuantity + 1 },
			fmt.Sprintf("quantity must be between %d and %d inclusive", advertisement.MinQuantity, advertisement.MaxQuantity)},
		{"без фотографий", func(r *models.CreateAdvertisementRequest) { r.PhotosCount = advertisement.MinPhotos - 1 },
			fmt.Sprintf("at least %d photo is required", advertisement.MinPhotos)},
		{"фотографий больше максимума", func(r *models.CreateAdvertisementRequest) { r.PhotosCount = advertisement.MaxPhotos + 1 },
			fmt.Sprintf("maximum %d photos allowed", advertisement.MaxPhotos)},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			request := ads.NewAdvertisementRequest()
			tc.modify(&request)

			base.Step(fmt.Sprintf("POST /advertisement с некорректным параметром: %s", tc.name))
			errResponse := ads.CreateAdvertisementExpectError(s.T(), s.loginData.Token, request, http.StatusBadRequest)

			base.Check("Возвращена ошибка валидации: " + tc.expectedMessage)
			assertions.ErrorEquals(s.T(), errResponse, apierrors.CodeInvalidReq, tc.expectedMessage)
		})
	}
}

// Negative: Создание объявления без авторизации или с неверным токеном возвращает 401
func (s *TestSuite) TestCreateAdvertisementUnauthorized() {
	cases := []struct {
		name            string
		token           string
		expectedMessage string
	}{
		{"без токена", "", apierrors.MsgAuthHeaderRequired},
		{"с неверным токеном", apierrors.InvalidToken, apierrors.MsgInvalidToken},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			base.Step("POST /advertisement " + tc.name)
			errResponse := ads.CreateAdvertisementExpectError(s.T(), tc.token, ads.NewAdvertisementRequest(),
				http.StatusUnauthorized)

			base.Check("Возвращена ошибка авторизации: " + tc.expectedMessage)
			assertions.ErrorEquals(s.T(), errResponse, apierrors.CodeUnauthorized, tc.expectedMessage)
		})
	}
}
