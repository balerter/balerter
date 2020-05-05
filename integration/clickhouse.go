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
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

type ClickhouseTestSuite struct {
	suite.Suite

	container testcontainers.Container
	IP        string
	Port      nat.Port
	ctx       context.Context
}

func (suite *ClickhouseTestSuite) SetupSuite() {
	var err error
	suite.ctx = context.Background()

	wd, err := os.Getwd()
	if err != nil {
		suite.Failf("error get workdir, %s", err.Error())
		return
	}

	req := testcontainers.ContainerRequest{
		Image:        "yandex/clickhouse-server",
		ExposedPorts: []string{"9000/tcp", "8123:8123"},
		WaitingFor:   wait.ForLog("Include not found: clickhouse_compression"),
		//WaitingFor: wait.ForListeningPort(nat.Port("9000")),
		BindMounts: map[string]string{
			path.Join(wd, "assets/clickhouse/data.sql"): "/docker-entrypoint-initdb.d/data.sql",
		},
	}
	suite.container, err = testcontainers.GenericContainer(suite.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		suite.T().Fatalf("error start continer, %v", err)
	}

	suite.IP, err = suite.container.Host(suite.ctx)
	if err != nil {
		suite.T().Fatalf("error get host, %v", err)
		return
	}

	suite.Port, err = suite.container.MappedPort(suite.ctx, "9000")
	if err != nil {
		suite.T().Fatalf("error get port, %v", err)
		return
	}
}

func (suite *ClickhouseTestSuite) TearDownSuite() {
	if err := suite.container.Terminate(suite.ctx); err != nil {
		suite.T().Fatalf("error terminate container, %v", err)
	}
}

func (suite *ClickhouseTestSuite) TestGetData() {

	cfg := `datasources:
  clickhouse:
    - name: ch1
      host: {HOST}
      port: {PORT}
      username: default
      password: 
      database: default
global:
  luaModulesPath: ../modules/?/init.lua
`

	cfg = strings.Replace(cfg, "{HOST}", suite.IP, 1)
	cfg = strings.Replace(cfg, "{PORT}", suite.Port.Port(), 1)

	cmd := exec.Command("./balerter", "-config=stdin", "-once", "-script=./assets/clickhouse/script.lua")

	bufOut := bytes.NewBuffer([]byte{})
	cmd.Stdout = bufOut
	cmd.Stdin = bytes.NewBuffer([]byte(cfg))

	log.Printf("sleep 5")
	time.Sleep(time.Second * 5)

	err := cmd.Run()

	expectedOut := `{
    1 = {
        id = 1
        name = John
        balance = 10.1
    }
    2 = {
        id = 2
        name = Bill
        balance = -10.1
    }
    3 = {
        id = 3
        name = Mark
        balance = 0
    }
}

`

	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedOut, bufOut.String())
}
