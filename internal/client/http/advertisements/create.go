package advertisements

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisements/models"
)

// HttpCreateAdvertisement отправляет multipart-запрос, исключая поля из omittedFields.
func HttpCreateAdvertisement(t *testing.T, token string, data models.CreateAdvertisementRequest, omittedFields ...string) *http.Response {
	fields := map[string]string{
		"title":       data.Title,
		"description": data.Description,
		"price":       strconv.FormatFloat(data.Price, 'f', -1, 64),
		"quantity":    strconv.Itoa(data.Quantity),
	}
	for _, field := range omittedFields {
		delete(fields, field)
	}
	request := apiRunner.GetRunner().Auth(token).Create().Post(path.AdvertisementPath).
		MultipartFormData("", []string{}...)
	for name, value := range fields {
		request.MultipartFormData(name, value)
	}
	if !slices.Contains(omittedFields, "photos") {
		request.MultipartFile("photos", data.Photos...)
	}
	return request.Expect(t).End().Response
}

// HttpGetAdvertisement получает объявление по его идентификатору.
func HttpGetAdvertisement(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.AdvertisementPath).
		Query("id", id).Expect(t).End().Response
}

// HttpDeleteAdvertisement удаляет объявление текущего пользователя.
func HttpDeleteAdvertisement(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.AdvertisementPath).
		Query("id", id).Expect(t).End().Response
}

// HttpGetAdvertisementPhotos получает фотографии объявления.
func HttpGetAdvertisementPhotos(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(fmt.Sprintf(path.AdvertisementPhotosPath, id)).Expect(t).End().Response
}

// HttpDownloadPhoto скачивает файл фотографии с ограничением времени ожидания.
func HttpDownloadPhoto(t *testing.T, photoURL string) *http.Response {
	t.Helper()
	downloadClient := &http.Client{Timeout: 30 * time.Second}
	response, err := downloadClient.Get(photoURL)
	require.NoError(t, err, "Не удалось скачать фотографию")
	return response
}
