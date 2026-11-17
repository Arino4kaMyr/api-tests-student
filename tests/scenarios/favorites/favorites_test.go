package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/apierrors"
	ads "api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	favs "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

const (
	malformedAdID           = "not-a-uuid"
	malformedAdIDPostError  = "invalid UUID length: 10"
	malformedAdIDDeleteText = "Invalid advertisement ID"
)

type TestSuite struct {
	suite.Suite
	owner authModels.LoginOkResponse // автор объявлений
	fan   authModels.LoginOkResponse // пользователь, добавляющий объявления в избранное
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Регистрация автора объявлений")
	s.owner = auth.RegisterRandomUser(s.T())

	base.Precondition("Регистрация пользователя, который будет добавлять объявления в избранное")
	s.fan = auth.RegisterRandomUser(s.T())
	require.NotEqual(s.T(), s.owner.User.Id, s.fan.User.Id, "Для тестов избранного нужны два разных пользователя")
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// createOwnerAdvertisement создаёт объявление автора и регистрирует его удаление после теста
func (s *TestSuite) createOwnerAdvertisement() models.Advertisement {
	created := ads.CreateAdvertisement(s.T(), s.owner.Token, ads.NewAdvertisementRequest())
	s.T().Cleanup(func() { ads.DeleteAdvertisementIfExists(s.T(), s.owner.Token, created.ID) })
	base.Precondition("Создано объявление автора " + created.ID)

	return created
}

// Positive: Добавленное объявление появляется в избранном, а после удаления из избранного пропадает
func (s *TestSuite) TestAddAndRemoveFavorite() {
	ad := s.createOwnerAdvertisement()

	s.Run("Добавляем объявление в избранное", func() {
		base.Step("POST /favorites для объявления " + ad.ID)
		favorite := favs.AddFavorite(s.T(), s.fan.Token, ad.ID)

		base.Check("В ответе вернулось добавленное объявление")
		require.NotEmpty(s.T(), favorite.FavoriteID, "favorite_id пустой")
		require.Equal(s.T(), ad.ID, favorite.Advertisement.ID, "В избранное добавлено другое объявление")
	})

	s.Run("Объявление есть в списке избранного", func() {
		base.Step("GET /favorites")
		list := favs.GetFavorites(s.T(), s.fan.Token)

		base.Check("Добавленное объявление присутствует в списке")
		assertions.ContainsFavoriteByAdID(s.T(), list.Items, ad.ID)
	})

	s.Run("Удаляем объявление из избранного", func() {
		base.Step("DELETE /favorites/" + ad.ID)
		favs.RemoveFavorite(s.T(), s.fan.Token, ad.ID)
	})

	s.Run("Объявления нет в списке избранного", func() {
		base.Step("GET /favorites")
		list := favs.GetFavorites(s.T(), s.fan.Token)

		base.Check("Удалённое объявление отсутствует в списке")
		assertions.NotContainsFavoriteByAdID(s.T(), list.Items, ad.ID)
	})
}

// Positive: Если удалить само объявление, оно исчезает из избранного пользователя
func (s *TestSuite) TestDeleteAdvertisementRemovesItFromFavorites() {
	ad := s.createOwnerAdvertisement()

	s.Run("Добавляем объявление в избранное", func() {
		base.Step("POST /favorites для объявления " + ad.ID)
		favs.AddFavorite(s.T(), s.fan.Token, ad.ID)

		base.Check("Объявление есть в списке избранного")
		assertions.ContainsFavoriteByAdID(s.T(), favs.GetFavorites(s.T(), s.fan.Token).Items, ad.ID)
	})

	s.Run("Автор удаляет объявление", func() {
		base.Step("DELETE /advertisement?id=" + ad.ID)
		ads.DeleteAdvertisement(s.T(), s.owner.Token, ad.ID)
	})

	s.Run("Объявление пропало из избранного", func() {
		base.Step("GET /favorites")
		list := favs.GetFavorites(s.T(), s.fan.Token)

		base.Check("Удалённое объявление отсутствует в списке избранного")
		assertions.NotContainsFavoriteByAdID(s.T(), list.Items, ad.ID)
	})
}

// Negative: Добавление в избранное несуществующего или некорректного ad_id
func (s *TestSuite) TestAddFavoriteInvalidAdID() {
	cases := []struct {
		name            string
		adID            string
		expectedStatus  int
		expectedError   string
		expectedMessage string
	}{
		{"несуществующее объявление", utils.RandomUUID(), http.StatusNotFound,
			apierrors.CodeNotFound, "Advertisement not found"},
		{"некорректный формат id", malformedAdID, http.StatusBadRequest,
			apierrors.CodeInvalidReq, malformedAdIDPostError},
		{"пустой id", "", http.StatusBadRequest, apierrors.CodeInvalidReq,
			"invalid UUID length: 0"},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			base.Step("POST /favorites с ad_id: " + tc.name)
			errResponse := favs.AddFavoriteExpectError(s.T(), s.fan.Token, tc.adID, tc.expectedStatus)

			base.Check("Возвращена ошибка: " + tc.expectedMessage)
			assertions.ErrorEquals(s.T(), errResponse, tc.expectedError, tc.expectedMessage)
		})
	}
}

// Negative: Удаление из избранного несуществующего, некорректного ad_id или объявления, не добавленного в избранное
func (s *TestSuite) TestRemoveFavoriteInvalidAdID() {
	notFavoriteAd := s.createOwnerAdvertisement()

	cases := []struct {
		name            string
		adID            string
		expectedStatus  int
		expectedError   string
		expectedMessage string
	}{
		{"несуществующее объявление", utils.RandomUUID(), http.StatusNotFound,
			apierrors.CodeNotFound, "Favorite entry not found"},
		{"существующее объявление не в избранном", notFavoriteAd.ID, http.StatusNotFound,
			apierrors.CodeNotFound, "Favorite entry not found"},
		{"некорректный формат id", malformedAdID, http.StatusBadRequest,
			apierrors.CodeInvalidReq, malformedAdIDDeleteText},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			base.Step("DELETE /favorites/{ad_id}: " + tc.name)
			errResponse := favs.RemoveFavoriteExpectError(s.T(), s.fan.Token, tc.adID, tc.expectedStatus)

			base.Check("Возвращена ошибка: " + tc.expectedMessage)
			assertions.ErrorEquals(s.T(), errResponse, tc.expectedError, tc.expectedMessage)
		})
	}
}

// Negative: Ручки избранного без авторизации или с неверным токеном возвращают 401
func (s *TestSuite) TestFavoritesUnauthorized() {
	ad := s.createOwnerAdvertisement()

	tokens := []struct {
		name            string
		token           string
		expectedMessage string
	}{
		{"без токена", "", apierrors.MsgAuthHeaderRequired},
		{"с неверным токеном", apierrors.InvalidToken, apierrors.MsgInvalidToken},
	}

	for _, tc := range tokens {
		s.Run("GET /favorites "+tc.name, func() {
			base.Step("GET /favorites " + tc.name)
			errResponse := favs.GetFavoritesExpectError(s.T(), tc.token, http.StatusUnauthorized)
			base.Check("Возвращена ошибка авторизации: " + tc.expectedMessage)
			assertions.ErrorEquals(s.T(), errResponse, apierrors.CodeUnauthorized, tc.expectedMessage)
		})
		s.Run("POST /favorites "+tc.name, func() {
			base.Step("POST /favorites " + tc.name)
			errResponse := favs.AddFavoriteExpectError(s.T(), tc.token, ad.ID, http.StatusUnauthorized)
			base.Check("Возвращена ошибка авторизации: " + tc.expectedMessage)
			assertions.ErrorEquals(s.T(), errResponse, apierrors.CodeUnauthorized, tc.expectedMessage)
		})
		s.Run("DELETE /favorites/{ad_id} "+tc.name, func() {
			base.Step("DELETE /favorites/{ad_id} " + tc.name)
			errResponse := favs.RemoveFavoriteExpectError(s.T(), tc.token, ad.ID, http.StatusUnauthorized)
			base.Check("Возвращена ошибка авторизации: " + tc.expectedMessage)
			assertions.ErrorEquals(s.T(), errResponse, apierrors.CodeUnauthorized, tc.expectedMessage)
		})
	}
}

// Negative: Повторное добавление того же объявления в избранное возвращает 409
func (s *TestSuite) TestAddFavoriteTwice() {
	ad := s.createOwnerAdvertisement()

	s.Run("Первое добавление проходит успешно", func() {
		base.Step("POST /favorites для объявления " + ad.ID)
		favs.AddFavorite(s.T(), s.fan.Token, ad.ID)
	})

	s.Run("Повторное добавление отклоняется", func() {
		base.Step("Повторный POST /favorites для объявления " + ad.ID)
		errResponse := favs.AddFavoriteExpectError(s.T(), s.fan.Token, ad.ID, http.StatusConflict)

		base.Check("Возвращена ошибка конфликта")
		assertions.ErrorEquals(s.T(), errResponse, apierrors.CodeAlreadyExists, "Advertisement is already in favorites")
	})

	s.Run("В избранном объявление встречается один раз", func() {
		list := favs.GetFavorites(s.T(), s.fan.Token)

		count := 0
		for _, item := range list.Items {
			if item.Advertisement.ID == ad.ID {
				count++
			}
		}
		require.Equal(s.T(), 1, count, "Объявление должно быть в избранном ровно один раз")
	})
}

// Negative: Добавление собственного объявления в избранное возвращает 403
func (s *TestSuite) TestAddOwnAdvertisementToFavorites() {
	ad := s.createOwnerAdvertisement()

	base.Step("POST /favorites автором объявления " + ad.ID)
	errResponse := favs.AddFavoriteExpectError(s.T(), s.owner.Token, ad.ID, http.StatusForbidden)

	base.Check("Возвращена ошибка доступа")
	assertions.ErrorEquals(s.T(), errResponse, apierrors.CodeForbidden, "Cannot add own advertisement to favorites")
}
