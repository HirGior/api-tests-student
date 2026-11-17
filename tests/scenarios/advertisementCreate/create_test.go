package advertisementCreate

import (
	"net/http"
	"strings"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
)

const (
	titlePrefix = "Тестовый горный велосипед "
	description = "Велосипед с фотографией для автоматической проверки создания объявления."
	price       = 12500.50
	quantity    = 2
	photoPath   = "../../../testdata/bicycle.png"
)

func (s *TestSuite) TestCreateAdvertisementWithAllParameters() {
	login := s.login()
	data := newAdvertisementRequest()
	created := advertisements.CreateAdvertisement(s.T(), login.Token, data)
	s.T().Logf("Объявление создано, ID=%s", created.ID)
	actual := advertisements.GetAdvertisement(s.T(), login.Token, created.ID)
	require.Equal(s.T(), created.ID, actual.ID, "ID полученного объявления отличается")
	require.Equal(s.T(), data.Title, actual.Title, "Название отличается от отправленного")
	require.Equal(s.T(), data.Description, actual.Description, "Описание отличается от отправленного")
	require.Equal(s.T(), data.Price, actual.Price, "Цена отличается от отправленной")
	require.Equal(s.T(), data.Quantity, actual.Quantity, "Количество отличается от отправленного")
	require.Equal(s.T(), login.User.Id, actual.UserID, "У объявления неверный владелец")
	s.T().Log("Поля объявления совпадают с отправленными")
	list := advertisements.GetAdvertisementsList(s.T(), login.Token, data.Title)
	var found bool
	for _, item := range list.Items {
		if item.ID == created.ID {
			require.Equal(s.T(), data.Title, item.Title, "Название объявления в поиске отличается")
			found = true
			break
		}
	}
	require.True(s.T(), found, "В поиске по %q отсутствует созданное объявление %s", data.Title, created.ID)
	s.T().Log("Объявление найдено по полному названию")
	photos := advertisements.GetAdvertisementPhotos(s.T(), login.Token, created.ID)
	require.Len(s.T(), photos, len(data.Photos), "Количество фотографий отличается от загруженного")
	require.Len(s.T(), created.Photos, len(data.Photos), "В ответе создания отсутствует фотография")
	require.NotEmpty(s.T(), photos[0].ID, "У фотографии отсутствует ID")
	require.Equal(s.T(), created.Photos[0].ID, photos[0].ID, "Список содержит другую фотографию")
	require.Equal(s.T(), created.Photos[0].URL, photos[0].URL, "URL фотографии изменился")
	require.NotEmpty(s.T(), photos[0].URL, "У фотографии отсутствует URL")
	downloaded, contentType := advertisements.DownloadPhoto(s.T(), photos[0].URL)
	require.NotEmpty(s.T(), downloaded, "Сервер вернул пустую фотографию")
	require.True(s.T(), strings.HasPrefix(contentType, "image/"), "Ожидался тип изображения, получен %q", contentType)
	require.True(s.T(), strings.HasPrefix(http.DetectContentType(downloaded), "image/"), "Скачанный файл должен содержать изображение")
	s.T().Log("Фотография доступна и корректна")
}
