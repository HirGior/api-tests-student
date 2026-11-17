package advertisementSearch

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	advertisementModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"
)

const (
	beginningPrefix = "Старт"
	middlePrefix    = "Середина"
	endPrefix       = "Конец"
	description     = "Тестовое объявление для поиска по подстроке и регистру."
	price           = 12500.50
	quantity        = 2
	photoPath       = "../../../testdata/bicycle.png"
)

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

// SetupSuite загружает настройки и авторизует пользователя для теста поиска.
func (s *TestSuite) SetupSuite() {
	base.SetupSuite()
	s.T().Setenv("DEBUG", "false")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

// TearDownSuite завершает выполнение теста поиска.
func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

func (s *TestSuite) TestSearchByTitle_Substring() {
	marker := utils.RandomString(6)
	searchBeginning := beginningPrefix + marker
	searchMiddle := middlePrefix + marker
	searchEnd := endPrefix + marker
	data := advertisementModels.CreateAdvertisementRequest{
		Title:       searchBeginning + " " + searchMiddle + " " + searchEnd,
		Description: description,
		Price:       price,
		Quantity:    quantity,
		Photos:      []string{photoPath},
	}
	created := advertisements.CreateAdvertisement(s.T(), s.loginData.Token, data)
	s.T().Logf("Подготовлено объявление для поиска, ID=%s", created.ID)
	cases := []struct {
		name   string
		search string
	}{
		{"начало", searchBeginning},
		{"середина", searchMiddle},
		{"конец", searchEnd},
		{"нижний регистр", strings.ToLower(searchMiddle)},
		{"верхний регистр", strings.ToUpper(searchMiddle)},
		{"смешанный регистр", utils.MixCase(searchMiddle)},
	}
	for _, tc := range cases {
		s.Run(tc.name, func() {
			list := advertisements.GetAdvertisementsList(s.T(), s.loginData.Token, tc.search)
			require.NotEmpty(s.T(), list.Items, "Поиск по %q вернул пустой список", tc.search)
			assertions.ContainsAdvertisementByTitle(s.T(), list.Items, data.Title)
			found := false
			for _, item := range list.Items {
				if item.ID == created.ID {
					require.Equal(s.T(), data.Title, item.Title, "Название найденного объявления отличается")
					found = true
					break
				}
			}
			require.True(s.T(), found, "Поиск по %q не вернул созданное объявление %s", tc.search, created.ID)
			s.T().Log("Созданное объявление найдено по ID и названию")
		})
	}
}
