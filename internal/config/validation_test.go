package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfig_Validate_NoError(t *testing.T) {
	cfg := &Config{}

	err := cfg.Validate()
	require.NoError(t, err)
}

func TestConfig_Validate_Error_Folder(t *testing.T) {
	cfg := &Config{Scripts: Scripts{Sources: ScriptsSources{UpdateInterval: 0, Folder: []ScriptSourceFolder{{Name: "", Path: ""}}}}}
	err := cfg.Validate()
	require.Error(t, err)
	assert.Equal(t, "name must be not empty", err.Error())

	cfg = &Config{Scripts: Scripts{Sources: ScriptsSources{UpdateInterval: 0, Folder: []ScriptSourceFolder{{Name: "name1", Path: ""}}}}}
	err = cfg.Validate()
	require.Error(t, err)
	assert.Equal(t, "path must be not empty", err.Error())

	cfg = &Config{Scripts: Scripts{Sources: ScriptsSources{UpdateInterval: 0, Folder: []ScriptSourceFolder{{Name: "name1", Path: "wrongpath"}}}}}
	err = cfg.Validate()
	require.Error(t, err)
	assert.Equal(t, "error read folder 'wrongpath', open wrongpath: no such file or directory", err.Error())
}

func TestConfig_Validate_Error_DSClickhouse(t *testing.T) {
	cfg := &Config{DataSources: DataSources{Clickhouse: []DataSourceClickhouse{{Name: "", Host: "", Port: 0, Username: ""}}}}
	err := cfg.Validate()
	require.Error(t, err)
	assert.Equal(t, "name must be not empty", err.Error())

	cfg = &Config{DataSources: DataSources{Clickhouse: []DataSourceClickhouse{{Name: "name1", Host: "", Port: 0, Username: ""}}}}
	err = cfg.Validate()
	require.Error(t, err)
	assert.Equal(t, "host must be defined", err.Error())

	cfg = &Config{DataSources: DataSources{Clickhouse: []DataSourceClickhouse{{Name: "name1", Host: "host", Port: 0, Username: ""}}}}
	err = cfg.Validate()
	require.Error(t, err)
	assert.Equal(t, "port must be defined", err.Error())

	cfg = &Config{DataSources: DataSources{Clickhouse: []DataSourceClickhouse{{Name: "name1", Host: "host", Port: 10, Username: ""}}}}
	err = cfg.Validate()
	require.Error(t, err)
	assert.Equal(t, "username must be defined", err.Error())
}

func TestConfig_Validate_Error_DSPrometheus(t *testing.T) {
	cfg := &Config{DataSources: DataSources{Prometheus: []DataSourcePrometheus{{Name: "", URL: ""}}}}
	err := cfg.Validate()
	require.Error(t, err)
	assert.Equal(t, "name must be not empty", err.Error())

	cfg = &Config{DataSources: DataSources{Prometheus: []DataSourcePrometheus{{Name: "name1", URL: ""}}}}
	err = cfg.Validate()
	require.Error(t, err)
	assert.Equal(t, "url must be not empty", err.Error())
}

func TestConfig_Validate_Error_CHSlack(t *testing.T) {
	cfg := &Config{Channels: Channels{Slack: []ChannelSlack{{Name: "", Token: "", Channel: ""}}}}
	err := cfg.Validate()
	require.Error(t, err)
	assert.Equal(t, "name must be not empty", err.Error())

	cfg = &Config{Channels: Channels{Slack: []ChannelSlack{{Name: "name1", Token: "", Channel: ""}}}}
	err = cfg.Validate()
	require.Error(t, err)
	assert.Equal(t, "token must be not empty", err.Error())

	cfg = &Config{Channels: Channels{Slack: []ChannelSlack{{Name: "name1", Token: "token", Channel: ""}}}}
	err = cfg.Validate()
	require.Error(t, err)
	assert.Equal(t, "channel must be not empty", err.Error())
}
