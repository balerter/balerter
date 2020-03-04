package config

import (
	"fmt"
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
	Storages    Storages    `json:"storages" yaml:"storages"`
	Global      Global      `json:"global" yaml:"global"`
}

func (cfg *Config) SetDefaults() {
	cfg.Scripts.SetDefaults()
	cfg.DataSources.SetDefaults()
	cfg.Channels.SetDefaults()
	cfg.Storages.SetDefaults()
	cfg.Global.SetDefaults()
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
