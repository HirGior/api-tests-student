package favorites

import (
	"net/http"
	"os"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/favorites"
)

func (s *TestSuite) TestFavoritesOwnAdvertisement() {
	owner := auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
	created := advertisements.CreateAdvertisement(s.T(), owner.Token, newAdvertisementRequest())
	before := favorites.GetFavorites(s.T(), owner.Token)
	response := favorites.AddFavoriteExpectError(s.T(), owner.Token, created.ID, http.StatusForbidden)
	s.checkFavoriteError(response, "forbidden", "Cannot add own advertisement to favorites")
	after := favorites.GetFavorites(s.T(), owner.Token)
	require.Equal(s.T(), before.Total, after.Total, "После отказа изменился размер избранного")
	require.Zero(s.T(), assertions.CountFavoritesByAdvertisementID(after.Items, created.ID), "Собственное объявление появилось в избранном")
	actual := advertisements.GetAdvertisement(s.T(), owner.Token, created.ID)
	require.Equal(s.T(), created.ID, actual.ID, "Собственное объявление должно сохраняться")
	s.T().Log("Собственное объявление не добавлено, само объявление доступно")
}
