package advertisementCreate

import (
	"strings"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/utils"
)

const (
	minTitleLength       = 1
	minDescriptionLength = 1
	minPrice             = 0
	minQuantity          = 1
)

func (s *TestSuite) TestCreateAdvertisementWithBoundaryValues() {
	login := s.login()
	title := titlePrefix + utils.RandomString(16)
	cases := []struct {
		name        string
		title       string
		description string
		price       float64
		quantity    int
	}{
		{"название 1 символ", utils.RandomString(minTitleLength), description, price, quantity},
		{"название 50 символов", utils.RandomString(maxTitleLength), description, price, quantity},
		{"описание 1 символ", title, strings.Repeat("a", minDescriptionLength), price, quantity},
		{"описание 500 символов", title, strings.Repeat("a", maxDescriptionLength), price, quantity},
		{"цена 0", title, description, minPrice, quantity},
		{"цена 1000000", title, description, maxPrice, quantity},
		{"количество 1", title, description, price, minQuantity},
		{"количество 100", title, description, price, maxQuantity},
	}
	for _, tc := range cases {
		s.Run(tc.name, func() {
			data := newAdvertisementRequest()
			data.Title = tc.title
			data.Description = tc.description
			data.Price = tc.price
			data.Quantity = tc.quantity
			created := advertisements.CreateAdvertisement(s.T(), login.Token, data)
			actual := advertisements.GetAdvertisement(s.T(), login.Token, created.ID)
			require.Equal(s.T(), created.ID, actual.ID, "ID объявления отличается")
			require.Equal(s.T(), data.Title, actual.Title, "Название отличается от отправленного")
			require.Equal(s.T(), data.Description, actual.Description, "Описание отличается от отправленного")
			require.Equal(s.T(), data.Price, actual.Price, "Цена отличается от отправленной")
			require.Equal(s.T(), data.Quantity, actual.Quantity, "Количество отличается от отправленного")
			s.T().Log("Объявление с граничным значением создано, поля проверены")
		})
	}
}
