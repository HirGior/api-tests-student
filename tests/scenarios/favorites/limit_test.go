package favorites

import (
	"net/http"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
)

func (s *TestSuite) TestFavoritesLimit() {
	s.T().Setenv("API_SLOW_MODE", "true")
	owner, buyer := s.loginUsers()
	before := favorites.GetFavorites(s.T(), buyer.Token)
	test := s.T()
	test.Cleanup(func() {
		restored := favorites.GetFavorites(test, buyer.Token)
		require.Equal(test, before.Total, restored.Total, "После очистки изменилось исходное количество избранного")
		for _, item := range before.Items {
			require.Equal(test, 1, assertions.CountFavoritesByAdvertisementID(restored.Items, item.Advertisement.ID),
				"При очистке потеряна исходная запись избранного %s", item.FavoriteID)
		}
		test.Log("Избранное восстановлено до исходного состояния")
	})
	available := models.MaxFavorites - before.Total
	require.GreaterOrEqual(s.T(), available, 0, "В избранном уже больше 50 записей")
	s.T().Logf("Подготовка: заполняем %d свободных мест избранного", available)
	createdIDs := []string{}
	for i := 0; i < available; i++ {
		created := advertisements.CreateAdvertisement(s.T(), owner.Token, newAdvertisementRequest())
		createdIDs = append(createdIDs, created.ID)
		favorites.AddFavorite(s.T(), buyer.Token, created.ID)
	}
	full := favorites.GetFavorites(s.T(), buyer.Token)
	require.Equal(s.T(), models.MaxFavorites, full.Total, "Не удалось заполнить избранное до лимита")
	for _, id := range createdIDs {
		require.Equal(s.T(), 1, assertions.CountFavoritesByAdvertisementID(full.Items, id), "Тестовое объявление %s отсутствует в заполненном избранном", id)
	}
	extra := advertisements.CreateAdvertisement(s.T(), owner.Token, newAdvertisementRequest())
	response := favorites.AddFavoriteExpectError(s.T(), buyer.Token, extra.ID, http.StatusBadRequest)
	s.checkFavoriteError(response, "limit_reached", "Favorites limit of 50 reached")
	after := favorites.GetFavorites(s.T(), buyer.Token)
	require.Equal(s.T(), models.MaxFavorites, after.Total, "После превышения лимита изменилось количество записей")
	require.Zero(s.T(), assertions.CountFavoritesByAdvertisementID(after.Items, extra.ID), "Объявление сверх лимита добавлено в избранное")
	s.T().Log("Лимит 50 записей проверен, лишнее объявление отклонено")
}
