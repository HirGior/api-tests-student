package myAdvertisement

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/managers/advertisements"
	advertisementModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/myAdvertisements"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

const (
	titlePrefix = "Моё тестовое объявление "
	description = "Тестовое объявление для проверки списка собственных объявлений."
	price       = 12500.50
	quantity    = 2
	photoPath   = "../../../testdata/bicycle.png"
)

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()
	s.T().Setenv("DEBUG", "false")

	base.Precondition("Авторизация пользователя с кредами из переменных окружения и получение его параметров")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// TestGetMyAdvertisementsPositive проверяет получение собственного тестового объявления.
func (s *TestSuite) TestGetMyAdvertisementsPositive() {
	data := advertisementModels.CreateAdvertisementRequest{
		Title:       titlePrefix + utils.RandomString(16),
		Description: description,
		Price:       price,
		Quantity:    quantity,
		Photos:      []string{photoPath},
	}
	created := advertisements.CreateAdvertisement(s.T(), s.loginData.Token, data)
	s.T().Logf("Подготовлено собственное объявление, ID=%s", created.ID)
	body := myAdvertisements.GetMyAdvertisements(s.T(), s.loginData.Token, http.StatusOK, 50, data.Title)
	items := gjson.Get(body, "items").Array()
	require.Len(s.T(), items, 1, "По уникальному названию должно находиться одно собственное объявление")
	item := items[0]
	require.Equal(s.T(), created.ID, item.Get("id").String(), "Вернулось другое объявление")
	require.Equal(s.T(), s.loginData.User.Id, item.Get("user_id").String(), "Объявление принадлежит другому пользователю")
	require.Equal(s.T(), data.Title, item.Get("title").String(), "Название отличается от отправленного")
	require.Equal(s.T(), data.Description, item.Get("description").String(), "Описание отличается от отправленного")
	require.NotEmpty(s.T(), item.Get("created_at").String(), "Не указано время создания")
	require.NotEmpty(s.T(), item.Get("updated_at").String(), "Не указано время обновления")
	require.NotEmpty(s.T(), item.Get("photos").Array(), "У объявления отсутствует фотография")
	s.T().Log("Собственное объявление найдено, поля и фотография проверены")
}

// TestGetMyAdvertisementsIncorrectToken проверяет ошибку получения списка с неверным токеном.
func (s *TestSuite) TestGetMyAdvertisementsIncorrectToken() {
	var advertisementsBody string
	s.Run("Получаем список собственных объявлений с неправильным авторизационным токеном", func() {
		advertisementsBody = myAdvertisements.GetMyAdvertisements(s.T(), "incorrect_token", http.StatusUnauthorized, 50)
	})

	s.Run("Проверяем, что показывается правильная ошибка", func() {
		require.Equal(s.T(), "unauthorized", gjson.Get(advertisementsBody, "error").String(),
			"Некорректное значение error параметра")
		require.Equal(s.T(), "Invalid or expired token", gjson.Get(advertisementsBody, "message").String(),
			"Некорректное значение message параметра")
	})
}
