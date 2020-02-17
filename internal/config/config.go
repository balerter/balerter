package config

import (
	"time"
)

const (
	defaultAPIAddress = "127.0.0.1:2000"
)

func New() *Config {
	cfg := &Config{}

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
	API                   API      `json:"api" yaml:"api"`
}

func (cfg *Global) SetDefaults() {
	cfg.API.SetDefaults()
}

type API struct {
	Address string `json:"address" yaml:"address"`
}

func (cfg *API) SetDefaults() {
	if cfg.Address == "" {
		cfg.Address = defaultAPIAddress
	}
}

type Channels struct {
	Slack    []ChannelSlack    `json:"slack" yaml:"slack"`
	Telegram []ChannelTelegram `json:"telegram" yaml:"telegram"`
}

type ProxyConfig struct {
	Address string           `json:"address" yaml:"address"`
	Auth    *ProxyAuthConfig `json:"auth" yaml:"auth"`
}

type ProxyAuthConfig struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

type ChannelTelegram struct {
	Name   string       `json:"name" yaml:"name"`
	Token  string       `json:"token" yaml:"token"`
	ChatID int64        `json:"chat_id" yaml:"chat_id"`
	Proxy  *ProxyConfig `json:"proxy"`
}

type ChannelSlack struct {
	Name    string `json:"name" yaml:"name"`
	Token   string `json:"token" yaml:"token"`
	Channel string `json:"channel" yaml:"channel"`
}

type DataSources struct {
	Clickhouse []DataSourceClickhouse `json:"clickhouse" yaml:"clickhouse"`
	Prometheus []DataSourcePrometheus `json:"prometheus" yaml:"prometheus"`
	Postgres   []DataSourcePostgres   `json:"postgres" yaml:"postgres"`
}

type DataSourcePostgres struct {
	Name        string `json:"name" yaml:"name"`
	Host        string `json:"host" yaml:"host"`
	Port        int    `json:"port" yaml:"port"`
	Username    string `json:"username" yaml:"username"`
	Password    string `json:"password" yaml:"password"`
	Database    string `json:"database" yaml:"database"`
	SSLMode     string `json:"ssl_mode" yaml:"ssl_mode"`
	SSLCertPath string `json:"ssl_cert_path" yaml:"ssl_cert_path"`
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
	SSLCertPath string `json:"ssl_cert_path" yaml:"ssl_cert_path"`
}

type Scripts struct {
	Sources ScriptsSources `json:"sources" yaml:"sources"`
}

type ScriptsSources struct {
	Folder []ScriptSourceFolder `json:"folder" yaml:"folder"`
}

type ScriptSourceFolder struct {
	UpdateInterval time.Duration `json:"update_interval" yaml:"update_interval"`
	Name           string        `json:"name" yaml:"name"`
	Path           string        `json:"path" yaml:"path"`
	Mask           string        `json:"mask" yaml:"mask"`
}
