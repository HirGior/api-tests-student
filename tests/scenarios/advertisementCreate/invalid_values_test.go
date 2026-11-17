package advertisementCreate

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/utils"
)

const (
	maxTitleLength           = 50
	maxDescriptionLength     = 500
	maxPrice                 = 1000000
	maxQuantity              = 100
	negativePrice            = -1
	zeroQuantity             = 0
	negativeQuantity         = -1
	priceRangeMessage        = "price must be between 0 and 1000000 inclusive"
	quantityRangeMessage     = "quantity must be between 1 and 100 inclusive"
	descriptionLengthMessage = "description must be between 1 and 500 characters"
)

func (s *TestSuite) TestCreateAdvertisementWithInvalidValues() {
	login := s.login()
	title := titlePrefix + utils.RandomString(16)
	cases := []struct {
		name            string
		title           string
		description     string
		price           float64
		quantity        int
		expectedMessage string
	}{
		{"пустое название", "", description, price, quantity, titleRequiredMessage},
		{"слишком длинное название", strings.Repeat("a", maxTitleLength+1), description, price, quantity,
			fmt.Sprintf("title must not exceed %d characters provided %d", maxTitleLength, maxTitleLength+1)},
		{"отрицательная цена", title, description, negativePrice, quantity, priceRangeMessage},
		{"цена выше максимальной", title, description, maxPrice + 1, quantity, priceRangeMessage},
		{"пустое описание", title, "", price, quantity, descriptionRequiredMessage},
		{"слишком длинное описание", title, strings.Repeat("a", maxDescriptionLength+1), price, quantity, descriptionLengthMessage},
		{"нулевое количество", title, description, price, zeroQuantity, quantityRangeMessage},
		{"отрицательное количество", title, description, price, negativeQuantity, quantityRangeMessage},
		{"количество выше максимального", title, description, price, maxQuantity + 1, quantityRangeMessage},
	}
	for _, tc := range cases {
		s.Run(tc.name, func() {
			data := newAdvertisementRequest()
			data.Title = tc.title
			data.Description = tc.description
			data.Price = tc.price
			data.Quantity = tc.quantity
			response := advertisements.CreateAdvertisementExpectError(s.T(), login.Token, data, http.StatusBadRequest)
			require.Equal(s.T(), invalidRequestError, response.Error, "Неверный код ошибки валидации")
			require.Equal(s.T(), tc.expectedMessage, response.Message, "Неверное сообщение о некорректном значении")
			s.T().Logf("Получен 400, сообщение проверено: %s", response.Message)
		})
	}
}
