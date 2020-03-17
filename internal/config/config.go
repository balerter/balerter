package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

func New(configSource string) (*Config, error) {
	cfg := &Config{
		viper: viper.New(),
	}

	cfg.viper.SetConfigName(configSource)
	cfg.viper.SetConfigType("yaml")
	cfg.viper.AddConfigPath(".")
	err := cfg.viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error read config file, %w", err)
	}

	cfg.SetDefaults()

	err = cfg.viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error unmarchal config, %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("error config validation")
	}

	return cfg, nil
}

type Config struct {
	viper *viper.Viper

	Scripts     Scripts     `json:"scripts" yaml:"scripts"`
	DataSources DataSources `json:"datasources" yaml:"datasources"`
	Channels    Channels    `json:"channels" yaml:"channels"`
	Storages    Storages    `json:"storages" yaml:"storages"`
	Global      Global      `json:"global" yaml:"global"`
}

func (cfg *Config) SetDefaults() {
	cfg.viper.SetDefault("global.storages.alert", "memory")
	cfg.viper.SetDefault("global.storages.kv", "memory")
	cfg.viper.SetDefault("global.api.address", "127.0.0.1:2000")
	cfg.viper.SetDefault("scripts.updateInterval", "1m")
}

func (cfg *Config) Validate() error {
	if err := cfg.Scripts.Validate(); err != nil {
		return fmt.Errorf("error Scripts validation, %w", err)
	}

	if err := cfg.DataSources.Validate(); err != nil {
		return fmt.Errorf("error DataSources validation, %w", err)
	}

	if err := cfg.Channels.Validate(); err != nil {
		return fmt.Errorf("error Channels validation, %w", err)
	}

	if err := cfg.Storages.Validate(); err != nil {
		return fmt.Errorf("error Storages validation, %w", err)
	}

	if err := cfg.Global.Validate(); err != nil {
		return fmt.Errorf("error global validation, %w", err)
	}

	return nil
}

type BasicAuth struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

func inArray(value string, arr []string) bool {
	for _, v := range arr {
		if value == v {
			return true
		}
	}

	return false
}

// check a slice for unique values,
// If founded non unique elements, returns a conflict element name. Else returns empty string
func checkUnique(data []string) string {
	m := map[string]struct{}{}

	for _, item := range data {
		item = strings.ToLower(item)
		if _, ok := m[item]; ok {
			return item
		}
		m[item] = struct{}{}
	}

	return ""
}
