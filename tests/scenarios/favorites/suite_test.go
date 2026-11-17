package favorites

import (
	"testing"

	"github.com/stretchr/testify/suite"

	base "api-tests-template/tests"
)

type TestSuite struct {
	suite.Suite
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

// SetupSuite загружает настройки для сценариев избранного.
func (s *TestSuite) SetupSuite() {
	base.SetupSuite()
	s.T().Setenv("DEBUG", "false")
}

// TearDownSuite завершает выполнение сценариев избранного.
func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}
