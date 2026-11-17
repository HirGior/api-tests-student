package advertisementCreate

import (
	"net/http"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
)

const (
	invalidRequestError        = "invalid_request"
	titleRequiredMessage       = "title is required"
	descriptionRequiredMessage = "description is required"
	photoRequiredMessage       = "at least 1 photo is required"
)

func (s *TestSuite) TestCreateAdvertisementWithoutRequiredFields() {
	login := s.login()
	cases := []struct {
		name            string
		omittedFields   []string
		expectedMessage string
	}{
		{"без всех полей", []string{"title", "description", "price", "quantity", "photos"}, titleRequiredMessage},
		{"без названия", []string{"title"}, titleRequiredMessage},
		{"без описания", []string{"description"}, descriptionRequiredMessage},
		{"без фотографии", []string{"photos"}, photoRequiredMessage},
	}
	for _, tc := range cases {
		s.Run(tc.name, func() {
			data := newAdvertisementRequest()
			response := advertisements.CreateAdvertisementExpectError(s.T(), login.Token, data, http.StatusBadRequest, tc.omittedFields...)
			require.Equal(s.T(), invalidRequestError, response.Error, "Неверный код ошибки валидации")
			require.Equal(s.T(), tc.expectedMessage, response.Message, "Неверное сообщение об отсутствующем поле")
			s.T().Logf("Получен 400, сообщение проверено: %s", response.Message)
		})
	}
}
