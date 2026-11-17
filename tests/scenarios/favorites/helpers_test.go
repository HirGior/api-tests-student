package favorites

import (
	"net/http"
	"os"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	advertisementModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/utils"
)

// loginUsers авторизует владельца объявления и другого пользователя для работы с избранным.
func (s *TestSuite) loginUsers() (authModels.LoginOkResponse, authModels.LoginOkResponse) {
	require.NotEmpty(s.T(), os.Getenv("SECOND_TEST_LOGIN"), "Заполните SECOND_TEST_LOGIN в .env.override")
	require.NotEmpty(s.T(), os.Getenv("SECOND_TEST_PASSWORD"), "Заполните SECOND_TEST_PASSWORD в .env.override")
	s.T().Log("Подготовка: авторизация владельца и второго пользователя")
	owner := auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
	buyer := auth.Login(s.T(), os.Getenv("SECOND_TEST_LOGIN"), os.Getenv("SECOND_TEST_PASSWORD"))
	require.NotEmpty(s.T(), owner.User.Id, "Не получен ID владельца")
	require.NotEmpty(s.T(), buyer.User.Id, "Не получен ID второго пользователя")
	require.NotEqual(s.T(), owner.User.Id, buyer.User.Id, "Для теста нужны два разных пользователя")
	return owner, buyer
}

// newAdvertisementRequest возвращает корректные данные объявления для сценариев избранного.
func newAdvertisementRequest() advertisementModels.CreateAdvertisementRequest {
	return advertisementModels.CreateAdvertisementRequest{
		Title:       titlePrefix + utils.RandomString(16),
		Description: description,
		Price:       price,
		Quantity:    quantity,
		Photos:      []string{photoPath},
	}
}

// checkUnauthorizedOperations проверяет отказ трёх операций и сохранность избранного.
func (s *TestSuite) checkUnauthorizedOperations(requestToken, expectedMessage string) {
	owner, buyer := s.loginUsers()
	created := advertisements.CreateAdvertisement(s.T(), owner.Token, newAdvertisementRequest())
	existing := advertisements.CreateAdvertisement(s.T(), owner.Token, newAdvertisementRequest())
	favorites.AddFavorite(s.T(), buyer.Token, existing.ID)
	before := favorites.GetFavorites(s.T(), buyer.Token)
	s.Run("получение списка", func() {
		response := favorites.GetFavoritesExpectError(s.T(), requestToken, http.StatusUnauthorized)
		s.checkFavoriteError(response, unauthorizedError, expectedMessage)
	})
	s.Run("добавление", func() {
		response := favorites.AddFavoriteExpectError(s.T(), requestToken, created.ID, http.StatusUnauthorized)
		s.checkFavoriteError(response, unauthorizedError, expectedMessage)
	})
	s.Run("удаление", func() {
		response := favorites.DeleteFavoriteExpectError(s.T(), requestToken, existing.ID, http.StatusUnauthorized)
		s.checkFavoriteError(response, unauthorizedError, expectedMessage)
	})
	after := favorites.GetFavorites(s.T(), buyer.Token)
	require.Equal(s.T(), before.Total, after.Total, "Отклонённые запросы изменили размер избранного")
	require.Equal(s.T(), 1, assertions.CountFavoritesByAdvertisementID(after.Items, existing.ID), "Отклонённое удаление изменило избранное")
	require.Zero(s.T(), assertions.CountFavoritesByAdvertisementID(after.Items, created.ID), "Отклонённое добавление изменило избранное")
	s.T().Log("Избранное сохранилось после отклонённых запросов")
}

// checkFavoriteError сравнивает код и сообщение ожидаемой ошибки избранного.
func (s *TestSuite) checkFavoriteError(response advertisementModels.ErrorResponse, expectedError, expectedMessage string) {
	require.Equal(s.T(), expectedError, response.Error, "Неверный код ошибки избранного")
	require.Equal(s.T(), expectedMessage, response.Message, "Неверное сообщение ошибки избранного")
	s.T().Logf("Ожидаемая ошибка проверена: %s", response.Message)
}
