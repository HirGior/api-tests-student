package favorites

import (
	"net/http"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/favorites"
)

func (s *TestSuite) TestFavoritesDuplicateAdvertisement() {
	owner, buyer := s.loginUsers()
	created := advertisements.CreateAdvertisement(s.T(), owner.Token, newAdvertisementRequest())
	added := favorites.AddFavorite(s.T(), buyer.Token, created.ID)
	before := favorites.GetFavorites(s.T(), buyer.Token)
	response := favorites.AddFavoriteExpectError(s.T(), buyer.Token, created.ID, http.StatusConflict)
	s.checkFavoriteError(response, "already_exists", "Advertisement is already in favorites")
	after := favorites.GetFavorites(s.T(), buyer.Token)
	require.Equal(s.T(), before.Total, after.Total, "Повторное добавление изменило размер избранного")
	require.Equal(s.T(), 1, assertions.CountFavoritesByAdvertisementID(after.Items, created.ID), "В избранном появились дубликаты")
	for _, item := range after.Items {
		if item.Advertisement.ID == created.ID {
			require.Equal(s.T(), added.FavoriteID, item.FavoriteID, "При повторном добавлении исходная запись изменилась")
		}
	}
	s.T().Log("Повторное добавление отклонено, исходная запись сохранена")
}
