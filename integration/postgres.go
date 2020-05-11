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
)

type PostgresTestSuite struct {
	suite.Suite

	postgresC    testcontainers.Container
	postgresIP   string
	postgresPort nat.Port
	ctx          context.Context
}

func (s *PostgresTestSuite) SetupSuite() {
	var err error
	s.ctx = context.Background()

	wd, err := os.Getwd()
	if err != nil {
		s.Failf("error get workdir, %s", err.Error())
		return
	}

	req := testcontainers.ContainerRequest{
		Image:        "postgres",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForLog("listening on IPv4 address \"0.0.0.0\", port 5432"),
		BindMounts: map[string]string{
			path.Join(wd, "assets/postgres/data.sql"): "/docker-entrypoint-initdb.d/data.sql",
		},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "secret",
			"POSTGRES_USER":     "user",
			"POSTGRES_DB":       "db",
		},
	}
	s.postgresC, err = testcontainers.GenericContainer(s.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		s.T().Fatalf("error start continer, %v", err)
	}

	s.postgresIP, err = s.postgresC.Host(s.ctx)
	if err != nil {
		s.T().Fatalf("error get host, %v", err)
		return
	}

	s.postgresPort, err = s.postgresC.MappedPort(s.ctx, "5432")
	if err != nil {
		s.T().Fatalf("error get port, %v", err)
		return
	}
}

func (s *PostgresTestSuite) TearDownSuite() {
	if err := s.postgresC.Terminate(s.ctx); err != nil {
		s.T().Fatalf("error terminate container, %v", err)
	}
}

func (s *PostgresTestSuite) TestGetData() {
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

	cfg = strings.Replace(cfg, "{HOST}", s.postgresIP, 1)
	cfg = strings.Replace(cfg, "{PORT}", s.postgresPort.Port(), 1)

	cmd := exec.Command("./balerter", "-config=stdin", "-once", "-script=./assets/postgres/script.lua")

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

	require.NoError(s.T(), err)
	assert.Equal(s.T(), expectedOut, bufOut.String())
}
