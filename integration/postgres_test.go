package integration

import (
	"bytes"
	"context"
	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
)

type PostgresTestSuite struct {
	suite.Suite

	postgresC    testcontainers.Container
	postgresIP   string
	postgresPort nat.Port
	ctx          context.Context
}

func (suite *PostgresTestSuite) SetupSuite() {
	var err error
	suite.ctx = context.Background()

	wd, err := os.Getwd()
	if err != nil {
		suite.Failf("error get workdir, %s", err.Error())
		return
	}

	req := testcontainers.ContainerRequest{
		Image:        "postgres",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForLog("listening on IPv4 address \"0.0.0.0\", port 5432"),
		BindMounts: map[string]string{
			path.Join(wd, "assets/postgres-get-data/data.sql"): "/docker-entrypoint-initdb.d/data.sql",
		},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "secret",
			"POSTGRES_USER":     "user",
			"POSTGRES_DB":       "db",
		},
	}
	suite.postgresC, err = testcontainers.GenericContainer(suite.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		suite.T().Fatalf("error start continer, %v", err)
	}

	suite.postgresIP, err = suite.postgresC.Host(suite.ctx)
	if err != nil {
		suite.T().Fatalf("error get host, %v", err)
		return
	}

	suite.postgresPort, err = suite.postgresC.MappedPort(suite.ctx, "5432")
	if err != nil {
		suite.T().Fatalf("error get port, %v", err)
		return
	}
}

func (suite *PostgresTestSuite) TearDownSuite() {
	if err := suite.postgresC.Terminate(suite.ctx); err != nil {
		suite.T().Fatalf("error terminate container, %v", err)
	}
}

func (suite *PostgresTestSuite) TestGetData() {

	cfg := `datasources:
  postgres:
    - name: pg1
      host: {HOST}
      port: {PORT}
      username: user
      password: secret
      database: db
      sslMode: disable
global:
  luaModulesPath: ../modules/?/init.lua
`

	cfg = strings.Replace(cfg, "{HOST}", suite.postgresIP, 1)
	cfg = strings.Replace(cfg, "{PORT}", suite.postgresPort.Port(), 1)

	cmd := exec.Command("./balerter", "-config=stdin", "-once", "-script=./assets/postgres-get-data/script.lua")

	bufOut := bytes.NewBuffer([]byte{})
	cmd.Stdout = bufOut
	cmd.Stdin = bytes.NewBuffer([]byte(cfg))

	err := cmd.Run()

	expectedOut := `{
    1 = {
        id = 1
        name = John
        balance = 10.2
    }
    2 = {
        id = 2
        name = Mark
        balance = 12
    }
    3 = {
        id = 3
        name = Peter
        balance = -15.4
    }
}

`

	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedOut, bufOut.String())
}

func TestPostgresSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}
