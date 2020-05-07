package integration

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestRunIntegrationTests(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
	suite.Run(t, new(ClickhouseTestSuite))
}
