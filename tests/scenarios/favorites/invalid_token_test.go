package favorites

const (
	incorrectToken              = "incorrect_token"
	invalidAuthorizationMessage = "Invalid or expired token"
)

func (s *TestSuite) TestFavoritesWithInvalidToken() {
	s.checkUnauthorizedOperations(incorrectToken, invalidAuthorizationMessage)
}
