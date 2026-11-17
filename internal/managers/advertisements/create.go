package advertisements

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	client "api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
)

// CreateAdvertisement создаёт объявление и сразу регистрирует его удаление после теста.
func CreateAdvertisement(t *testing.T, token string, data models.CreateAdvertisementRequest) models.Advertisement {
	t.Helper()
	raw := client.HttpCreateAdvertisement(t, token, data)
	body := readResponse(t, raw)
	var ad models.Advertisement
	err := json.Unmarshal(body, &ad)
	if ad.ID != "" {
		t.Cleanup(func() {
			cleanupAdvertisement(t, token, ad.ID)
			t.Log("Тестовое объявление удалено")
		})
	}
	require.Equal(t, http.StatusCreated, raw.StatusCode, "Создание объявления: %s", body)
	require.NoError(t, err, "Не удалось разобрать ответ создания объявления")
	require.NotEmpty(t, ad.ID, "Сервер не вернул ID созданного объявления")
	return ad
}

// GetAdvertisement возвращает объявление по ID и проверяет успешный ответ.
func GetAdvertisement(t *testing.T, token, id string) models.Advertisement {
	t.Helper()
	raw := client.HttpGetAdvertisement(t, token, id)
	body := readResponse(t, raw)
	require.Equal(t, http.StatusOK, raw.StatusCode, "Получение объявления %s: %s", id, body)
	var ad models.Advertisement
	require.NoError(t, json.Unmarshal(body, &ad), "Не удалось разобрать объявление")
	return ad
}

// DeleteAdvertisement удаляет объявление и проверяет отсутствие по тому же ID.
func DeleteAdvertisement(t *testing.T, token, id string) {
	t.Helper()
	raw := client.HttpDeleteAdvertisement(t, token, id)
	body := readResponse(t, raw)
	require.Equal(t, http.StatusNoContent, raw.StatusCode, "Удаление объявления %s: %s", id, body)
	check := client.HttpGetAdvertisement(t, token, id)
	checkBody := readResponse(t, check)
	require.Equal(t, http.StatusNotFound, check.StatusCode, "Объявление %s осталось после удаления: %s", id, checkBody)
}

// GetAdvertisementPhotos возвращает список фотографий объявления.
func GetAdvertisementPhotos(t *testing.T, token, id string) []models.Photo {
	t.Helper()
	raw := client.HttpGetAdvertisementPhotos(t, token, id)
	body := readResponse(t, raw)
	require.Equal(t, http.StatusOK, raw.StatusCode, "Получение фотографий %s: %s", id, body)
	var photos []models.Photo
	require.NoError(t, json.Unmarshal(body, &photos), "Не удалось разобрать список фотографий")
	return photos
}

// DownloadPhoto получает изображение по URL с ограничением времени ожидания.
func DownloadPhoto(t *testing.T, photoURL string) ([]byte, string) {
	t.Helper()
	response := client.HttpDownloadPhoto(t, photoURL)
	body := readResponse(t, response)
	require.Equal(t, http.StatusOK, response.StatusCode, "Фотография недоступна")
	return body, response.Header.Get("Content-Type")
}

// readResponse читает тело HTTP-ответа и закрывает его.
func readResponse(t *testing.T, response *http.Response) []byte {
	t.Helper()
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Не удалось прочитать HTTP-ответ")
	return body
}

// cleanupAdvertisement удаляет созданное тестом объявление, если оно ещё существует.
func cleanupAdvertisement(t *testing.T, token, id string) {
	t.Helper()
	check := client.HttpGetAdvertisement(t, token, id)
	body := readResponse(t, check)
	if check.StatusCode == http.StatusNotFound {
		return
	}
	require.Equal(t, http.StatusOK, check.StatusCode, "Не удалось проверить объявление %s перед очисткой: %s", id, body)
	DeleteAdvertisement(t, token, id)
}
