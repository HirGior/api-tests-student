package favorites

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	client "api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/managers/favorites/models"
)

// AddFavorite добавляет объявление и регистрирует очистку записи после теста.
func AddFavorite(t *testing.T, token, advertisementID string) models.Favorite {
	t.Helper()
	raw := client.HttpAddFavorite(t, token, models.AddFavoriteRequest{AdvertisementID: advertisementID})
	body := readResponse(t, raw)
	if raw.StatusCode == http.StatusCreated {
		t.Cleanup(func() {
			response := client.HttpDeleteFavorite(t, token, advertisementID)
			cleanupBody := readResponse(t, response)
			require.Contains(t, []int{http.StatusNoContent, http.StatusNotFound}, response.StatusCode,
				"Не удалось очистить избранное для объявления %s: %s", advertisementID, cleanupBody)
		})
	}
	require.Equal(t, http.StatusCreated, raw.StatusCode, "Добавление в избранное: %s", body)
	var favorite models.Favorite
	require.NoError(t, json.Unmarshal(body, &favorite), "Не удалось разобрать запись избранного")
	require.NotEmpty(t, favorite.FavoriteID, "Сервер не вернул ID записи избранного")
	return favorite
}

// GetFavorites возвращает избранное пользователя и проверяет полноту страницы.
func GetFavorites(t *testing.T, token string) models.FavoritesList {
	t.Helper()
	raw := client.HttpGetFavorites(t, token)
	body := readResponse(t, raw)
	require.Equal(t, http.StatusOK, raw.StatusCode, "Получение избранного: %s", body)
	var list models.FavoritesList
	require.NoError(t, json.Unmarshal(body, &list), "Не удалось разобрать список избранного")
	require.LessOrEqual(t, list.TotalPages, 1, "Все записи избранного должны помещаться на одной странице")
	require.Len(t, list.Items, list.Total, "Ответ содержит неполный список избранного")
	return list
}

// DeleteFavorite удаляет запись избранного с ожидаемым статусом 204.
func DeleteFavorite(t *testing.T, token, advertisementID string) {
	t.Helper()
	raw := client.HttpDeleteFavorite(t, token, advertisementID)
	body := readResponse(t, raw)
	require.Equal(t, http.StatusNoContent, raw.StatusCode, "Удаление из избранного %s: %s", advertisementID, body)
}

// readResponse читает и закрывает тело ответа API избранного.
func readResponse(t *testing.T, response *http.Response) []byte {
	t.Helper()
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Не удалось прочитать ответ API избранного")
	return body
}
