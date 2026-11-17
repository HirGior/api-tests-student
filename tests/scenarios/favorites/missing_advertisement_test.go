package favorites

import (
	"net/http"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/favorites"
)

func (s *TestSuite) TestFavoritesWithMissingAdvertisement() {
	owner, buyer := s.loginUsers()
	created := advertisements.CreateAdvertisement(s.T(), owner.Token, newAdvertisementRequest())
	advertisements.DeleteAdvertisement(s.T(), owner.Token, created.ID)
	s.T().Log("Подготовлен ID удалённого объявления")
	s.Run("добавление несуществующего объявления", func() {
		response := favorites.AddFavoriteExpectError(s.T(), buyer.Token, created.ID, http.StatusNotFound)
		s.checkFavoriteError(response, "not_found", "Advertisement not found")
	})
	s.Run("удаление несуществующего объявления", func() {
		response := favorites.DeleteFavoriteExpectError(s.T(), buyer.Token, created.ID, http.StatusNotFound)
		s.checkFavoriteError(response, "not_found", "Favorite entry not found")
	})
}
