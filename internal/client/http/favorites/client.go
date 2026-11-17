package favorites

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/favorites/models"
)

// HttpAddFavorite добавляет объявление в избранное пользователя.
func HttpAddFavorite(t *testing.T, token string, data models.AddFavoriteRequest) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesPath).
		JSON(data).Expect(t).End().Response
}

// HttpGetFavorites получает первую страницу избранного с максимальным размером.
func HttpGetFavorites(t *testing.T, token string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.FavoritesPath).
		Query("page", "1").Query("page_size", strconv.Itoa(models.MaxFavorites)).
		Expect(t).End().Response
}

// HttpDeleteFavorite удаляет объявление из избранного по ID объявления.
func HttpDeleteFavorite(t *testing.T, token, advertisementID string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(fmt.Sprintf(path.FavoritePath, advertisementID)).Expect(t).End().Response
}

// HttpAddFavoriteJSON отправляет JSON-тело для проверки валидации избранного.
func HttpAddFavoriteJSON(t *testing.T, token, body string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesPath).
		ContentType("application/json").Body(body).Expect(t).End().Response
}
