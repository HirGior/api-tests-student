package favorites

const (
	unauthorizedError            = "unauthorized"
	authorizationRequiredMessage = "Authorization header required"
)

func (s *TestSuite) TestFavoritesWithoutToken() {
	s.checkUnauthorizedOperations("", authorizationRequiredMessage)
}
