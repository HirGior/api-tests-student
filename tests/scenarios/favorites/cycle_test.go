package favorites

import (
	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/favorites"
)

const (
	titlePrefix = "Велосипед для избранного "
	description = "Тестовое объявление для проверки добавления и удаления из избранного."
	price       = 15000.50
	quantity    = 2
	photoPath   = "../../../testdata/bicycle.png"
)

func (s *TestSuite) TestFavoritesAddAndRemove() {
	owner, buyer := s.loginUsers()
	data := newAdvertisementRequest()
	created := advertisements.CreateAdvertisement(s.T(), owner.Token, data)
	s.T().Logf("Объявление создано, ID=%s", created.ID)

	added := favorites.AddFavorite(s.T(), buyer.Token, created.ID)
	require.Equal(s.T(), created.ID, added.Advertisement.ID, "В избранное добавлено другое объявление")
	require.Equal(s.T(), data.Title, added.Advertisement.Title, "Название в ответе добавления отличается")
	require.NotEmpty(s.T(), added.AddedAt, "Не указано время добавления в избранное")
	list := favorites.GetFavorites(s.T(), buyer.Token)
	require.Equal(s.T(), 1, assertions.CountFavoritesByAdvertisementID(list.Items, created.ID),
		"Созданное объявление должно присутствовать в избранном ровно один раз")
	for _, item := range list.Items {
		if item.Advertisement.ID == created.ID {
			require.Equal(s.T(), added.FavoriteID, item.FavoriteID, "ID записи избранного отличается")
			require.Equal(s.T(), data.Title, item.Advertisement.Title, "Название в списке избранного отличается")
		}
	}
	s.T().Log("Объявление добавлено в избранное и найдено в списке")

	favorites.DeleteFavorite(s.T(), buyer.Token, created.ID)
	list = favorites.GetFavorites(s.T(), buyer.Token)
	require.Zero(s.T(), assertions.CountFavoritesByAdvertisementID(list.Items, created.ID),
		"После удаления объявление осталось в избранном")
	s.T().Log("Объявление удалено из избранного")
	actual := advertisements.GetAdvertisement(s.T(), owner.Token, created.ID)
	require.Equal(s.T(), created.ID, actual.ID, "Само объявление должно сохраняться после удаления из избранного")
	require.Equal(s.T(), owner.User.Id, actual.UserID, "Владелец объявления изменился")
	s.T().Log("Само объявление осталось доступным владельцу")
}
