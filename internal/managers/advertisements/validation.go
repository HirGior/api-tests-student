package advertisements

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	client "api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
)

// CreateAdvertisementExpectError проверяет отказ создания и возвращает ошибку API.
func CreateAdvertisementExpectError(t *testing.T, token string, data models.CreateAdvertisementRequest, expectedStatus int, omittedFields ...string) models.ErrorResponse {
	t.Helper()
	raw := client.HttpCreateAdvertisement(t, token, data, omittedFields...)
	body := readResponse(t, raw)
	var ad models.Advertisement
	if json.Unmarshal(body, &ad) == nil && ad.ID != "" {
		t.Cleanup(func() {
			cleanupToken := token
			if expectedStatus == http.StatusUnauthorized {
				login := auth.Login(t, os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
				cleanupToken = login.Token
			}
			DeleteAdvertisement(t, cleanupToken, ad.ID)
			t.Log("Неожиданно созданное объявление удалено")
		})
	}
	require.Equal(t, expectedStatus, raw.StatusCode, "Неверный статус отказа создания: %s", body)
	var response models.ErrorResponse
	require.NoError(t, json.Unmarshal(body, &response), "Не удалось разобрать ошибку: %s", body)
	require.Empty(t, ad.ID, "При некорректном запросе создано объявление %s", ad.ID)
	return response
}
