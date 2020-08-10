package integration

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestRunIntegrationTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integration tests in short mode")
	}
	suite.Run(t, new(PostgresTestSuite))
	suite.Run(t, new(ClickhouseTestSuite))
}
