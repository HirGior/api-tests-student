package advertisementCreate

import (
	"net/http"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
)

const (
	unauthorizedError            = "unauthorized"
	authorizationRequiredMessage = "Authorization header required"
)

func (s *TestSuite) TestCreateAdvertisementWithoutToken() {
	data := newAdvertisementRequest()
	s.T().Log("Отправляем корректное объявление без токена")
	response := advertisements.CreateAdvertisementExpectError(s.T(), "", data, http.StatusUnauthorized)
	require.Equal(s.T(), unauthorizedError, response.Error, "Неверный код ошибки авторизации")
	require.Equal(s.T(), authorizationRequiredMessage, response.Message, "Неверное сообщение об отсутствующем токене")
	s.T().Logf("Получен 401, сообщение проверено: %s", response.Message)
}
