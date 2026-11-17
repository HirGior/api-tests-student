package advertisementCreate

import (
	"net/http"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
)

const (
	incorrectToken              = "incorrect_token"
	invalidAuthorizationMessage = "Invalid or expired token"
)

func (s *TestSuite) TestCreateAdvertisementWithInvalidToken() {
	data := newAdvertisementRequest()
	s.T().Log("Отправляем корректное объявление с неверным токеном")
	response := advertisements.CreateAdvertisementExpectError(s.T(), incorrectToken, data, http.StatusUnauthorized)
	require.Equal(s.T(), unauthorizedError, response.Error, "Неверный код ошибки авторизации")
	require.Equal(s.T(), invalidAuthorizationMessage, response.Message, "Неверное сообщение о некорректном токене")
	s.T().Logf("Получен 401, сообщение проверено: %s", response.Message)
}
