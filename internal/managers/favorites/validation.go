package favorites

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	client "api-tests-template/internal/client/http/favorites"
	advertisementModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/favorites/models"
)

// AddFavoriteExpectError проверяет отказ добавления объявления в избранное.
func AddFavoriteExpectError(t *testing.T, token, advertisementID string, expectedStatus int) advertisementModels.ErrorResponse {
	t.Helper()
	raw := client.HttpAddFavorite(t, token, models.AddFavoriteRequest{AdvertisementID: advertisementID})
	return readErrorResponse(t, raw, expectedStatus)
}

// AddFavoriteJSONExpectError проверяет ошибку добавления с некорректным JSON-телом.
func AddFavoriteJSONExpectError(t *testing.T, token, body string, expectedStatus int) advertisementModels.ErrorResponse {
	t.Helper()
	raw := client.HttpAddFavoriteJSON(t, token, body)
	return readErrorResponse(t, raw, expectedStatus)
}

// DeleteFavoriteExpectError проверяет отказ удаления записи избранного.
func DeleteFavoriteExpectError(t *testing.T, token, advertisementID string, expectedStatus int) advertisementModels.ErrorResponse {
	t.Helper()
	raw := client.HttpDeleteFavorite(t, token, advertisementID)
	return readErrorResponse(t, raw, expectedStatus)
}

// GetFavoritesExpectError проверяет отказ получения списка избранного.
func GetFavoritesExpectError(t *testing.T, token string, expectedStatus int) advertisementModels.ErrorResponse {
	t.Helper()
	raw := client.HttpGetFavorites(t, token)
	return readErrorResponse(t, raw, expectedStatus)
}

// readErrorResponse проверяет статус и читает сообщение ошибки избранного.
func readErrorResponse(t *testing.T, raw *http.Response, expectedStatus int) advertisementModels.ErrorResponse {
	t.Helper()
	body := readResponse(t, raw)
	require.Equal(t, expectedStatus, raw.StatusCode, "Неверный статус ошибки избранного: %s", body)
	var response advertisementModels.ErrorResponse
	require.NoError(t, json.Unmarshal(body, &response), "Не удалось разобрать ошибку избранного: %s", body)
	require.NotEmpty(t, response.Error, "Отсутствует код ошибки: %s", body)
	require.NotEmpty(t, response.Message, "Отсутствует сообщение ошибки: %s", body)
	return response
}
