package config

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/channels"
	"github.com/balerter/balerter/internal/config/datasources"
	"github.com/balerter/balerter/internal/config/global"
	"github.com/balerter/balerter/internal/config/scripts"
	"github.com/balerter/balerter/internal/config/storages"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
)

var StdIn io.Reader = os.Stdin

func New(configSource string) (*Config, error) {
	cfg := &Config{
		Scripts:     &scripts.Scripts{},
		DataSources: &datasources.DataSources{},
		Channels:    &channels.Channels{},
		Storages:    &storages.Storages{},
		Global:      &global.Global{},
	}

	var data []byte
	var err error

	if configSource == "stdin" {
		data, err = ioutil.ReadAll(StdIn)
	} else {
		data, err = ioutil.ReadFile(configSource)
	}

	if err != nil {
		return nil, fmt.Errorf("error read config file, %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("error parse config file, %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("error config validation, %w", err)
	}

	return cfg, nil
}

type Config struct {
	Scripts     *scripts.Scripts         `json:"scripts" yaml:"scripts"`
	DataSources *datasources.DataSources `json:"datasources" yaml:"datasources"`
	Channels    *channels.Channels       `json:"channels" yaml:"channels"`
	Storages    *storages.Storages       `json:"storages" yaml:"storages"`
	Global      *global.Global           `json:"global" yaml:"global"`
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
