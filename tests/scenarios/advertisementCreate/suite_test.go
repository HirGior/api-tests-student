package advertisementCreate

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"
)

type TestSuite struct {
	suite.Suite
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

// SetupSuite загружает общие настройки тестов создания объявлений.
func (s *TestSuite) SetupSuite() {
	base.SetupSuite()
	s.T().Setenv("DEBUG", "false")
}

// login авторизует пользователя перед сценарием, которому нужен действующий токен.
func (s *TestSuite) login() authModels.LoginOkResponse {
	s.T().Log("Подготовка данных и авторизация")
	return auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

// newAdvertisementRequest возвращает корректные данные с уникальным названием.
func newAdvertisementRequest() models.CreateAdvertisementRequest {
	return models.CreateAdvertisementRequest{
		Title:       titlePrefix + utils.RandomString(16),
		Description: description,
		Price:       price,
		Quantity:    quantity,
		Photos:      []string{photoPath},
	}
}

// TearDownSuite завершает выполнение набора тестов.
func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}
