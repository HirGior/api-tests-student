package favorites

import (
	"net/http"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/favorites"
)

const (
	incorrectAdvertisementID        = "incorrect_id"
	invalidRequestError             = "invalid_request"
	missingAdvertisementIDMessage   = "Key: 'AddFavoriteRequest.AdvertisementID' Error:Field validation for 'AdvertisementID' failed on the 'required' tag"
	emptyAdvertisementIDMessage     = "invalid UUID length: 0"
	incorrectAdvertisementIDMessage = "invalid UUID length: 12"
	invalidDeleteIDMessage          = "Invalid advertisement ID"
)

func (s *TestSuite) TestFavoritesWithInvalidRequest() {
	_, buyer := s.loginUsers()
	before := favorites.GetFavorites(s.T(), buyer.Token)
	cases := []struct {
		name            string
		body            string
		expectedMessage string
	}{
		{"отсутствует ID", `{}`, missingAdvertisementIDMessage},
		{"пустой ID", `{"advertisement_id":""}`, emptyAdvertisementIDMessage},
		{"неверный формат ID", `{"advertisement_id":"incorrect_id"}`, incorrectAdvertisementIDMessage},
	}
	for _, tc := range cases {
		s.Run(tc.name, func() {
			response := favorites.AddFavoriteJSONExpectError(s.T(), buyer.Token, tc.body, http.StatusBadRequest)
			s.checkFavoriteError(response, invalidRequestError, tc.expectedMessage)
		})
	}
	s.Run("удаление с неверным форматом ID", func() {
		response := favorites.DeleteFavoriteExpectError(s.T(), buyer.Token, incorrectAdvertisementID, http.StatusBadRequest)
		s.checkFavoriteError(response, invalidRequestError, invalidDeleteIDMessage)
	})
	after := favorites.GetFavorites(s.T(), buyer.Token)
	require.Equal(s.T(), before.Total, after.Total, "Некорректные запросы изменили размер избранного")
}
