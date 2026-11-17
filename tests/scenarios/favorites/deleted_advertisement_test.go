package favorites

import (
	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/favorites"
)

func (s *TestSuite) TestDeletedAdvertisementDisappearsFromFavorites() {
	owner, buyer := s.loginUsers()
	data := newAdvertisementRequest()
	created := advertisements.CreateAdvertisement(s.T(), owner.Token, data)
	s.T().Logf("Объявление создано, ID=%s", created.ID)

	added := favorites.AddFavorite(s.T(), buyer.Token, created.ID)
	require.Equal(s.T(), created.ID, added.Advertisement.ID, "В избранное добавлено другое объявление")
	list := favorites.GetFavorites(s.T(), buyer.Token)
	require.Equal(s.T(), 1, assertions.CountFavoritesByAdvertisementID(list.Items, created.ID),
		"Перед удалением объявление должно присутствовать в избранном ровно один раз")
	s.T().Log("Объявление добавлено в избранное и найдено в списке")

	advertisements.DeleteAdvertisement(s.T(), owner.Token, created.ID)
	s.T().Log("Владелец удалил объявление, получение по ID возвращает 404")
	list = favorites.GetFavorites(s.T(), buyer.Token)
	require.Zero(s.T(), assertions.CountFavoritesByAdvertisementID(list.Items, created.ID),
		"После удаления самого объявления оно осталось в избранном пользователя")
	s.T().Log("Удалённое объявление исчезло из избранного")
}
