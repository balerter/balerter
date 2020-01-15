package config

import (
	"time"
)

const (
	defaultScriptsUpdateInterval time.Duration = time.Second * 60
)

func New() *Config {
	cfg := &Config{
		Scripts: Scripts{
			Sources: ScriptsSources{
				UpdateInterval: defaultScriptsUpdateInterval,
			},
		},
	}

	return cfg
}

type Config struct {
	Scripts     Scripts     `json:"scripts" yaml:"scripts"`
	DataSources DataSources `json:"datasources" yaml:"datasources"`
	Channels    Channels    `json:"channels" yaml:"channels"`
	Global      Global      `json:"global" yaml:"global"`
}

type Global struct {
	SendStartNotification []string `json:"send_start_notification" yaml:"send_start_notification"`
	SendStopNotification  []string `json:"send_stop_notification" yaml:"send_stop_notification"`
}

type Channels struct {
	Slack []ChannelSlack `json:"slack" yaml:"slack"`
}

type ChannelSlack struct {
	Name string `json:"name" yaml:"name"`
	URL  string `json:"url" yaml:"url"`
}

type DataSources struct {
	Clickhouse []DataSourceClickhouse `json:"clickhouse" yaml:"clickhouse"`
	Prometheus []DataSourcePrometheus `json:"prometheus" yaml:"prometheus"`
}

type DataSourcePrometheus struct {
	Name      string    `json:"name" yaml:"name"`
	URL       string    `json:"url" yaml:"url"`
	BasicAuth BasicAuth `json:"basic_auth" yaml:"basic_auth"`
}

type BasicAuth struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

type DataSourceClickhouse struct {
	Name        string `json:"name" yaml:"name"`
	Host        string `json:"host" yaml:"host"`
	Port        int    `json:"port" yaml:"port"`
	Username    string `json:"username" yaml:"username"`
	Password    string `json:"password" yaml:"password"`
	Database    string `json:"database" yaml:"database"`
	SSLMode     string `json:"ssl_mode" yaml:"ssl_mode"`
	SSLCertPath string `json:"ssl_cert_path" yaml:"ssl_cert_path"`
}

type Scripts struct {
	Sources ScriptsSources `json:"sources" yaml:"sources"`
}

type ScriptsSources struct {
	UpdateInterval time.Duration        `json:"update_interval" yaml:"update_interval"`
	Folder         []ScriptSourceFolder `json:"folder" yaml:"folder"`
}

type ScriptSourceFolder struct {
	Name string `json:"name" yaml:"name"`
	Path string `json:"path" yaml:"path"`
	Mask string `json:"mask" yaml:"mask"`
}
