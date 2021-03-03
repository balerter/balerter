package config

import (
	"flag"
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

func New() (*Config, *Flags, error) {
	cfg := &Config{
		Scripts:     scripts.Scripts{},
		DataSources: datasources.DataSources{},
		Channels:    channels.Channels{},
		Storages:    storages.Storages{},
		Global:      global.Global{},
	}

	var data []byte
	var err error

	flg := &Flags{}

	flag.StringVar(&flg.ConfigFilePath, "config", "config.yml", "configuration source. Currently supports only path to yaml file and 'stdin'.")
	flag.StringVar(&flg.LogLevel, "logLevel", "INFO", "log level. ERROR, INFO or DEBUG")
	flag.BoolVar(&flg.Debug, "debug", false, "debug mode")
	flag.BoolVar(&flg.Once, "once", false, "once run scripts and exit")
	flag.StringVar(&flg.Script, "script", "", "ignore all script sources and runs only one script. Meta-tag @ignore will be ignored")
	flag.BoolVar(&flg.AsJSON, "json", false, "output json format")
	flag.Parse()

	if flg.ConfigFilePath == "stdin" {
		data, err = ioutil.ReadAll(StdIn)
	} else {
		data, err = ioutil.ReadFile(flg.ConfigFilePath)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("error read config file, %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, nil, fmt.Errorf("error parse config file, %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, nil, fmt.Errorf("error config validation, %w", err)
	}

	return cfg, flg, nil
}

type Flags struct {
	ConfigFilePath string
	LogLevel       string
	Debug          bool
	Once           bool
	Script         string
	AsJSON         bool
}

type Config struct {
	Scripts     scripts.Scripts         `json:"scripts" yaml:"scripts"`
	DataSources datasources.DataSources `json:"datasources" yaml:"datasources"`
	Channels    channels.Channels       `json:"channels" yaml:"channels"`
	Storages    storages.Storages       `json:"storages" yaml:"storages"`
	Global      global.Global           `json:"global" yaml:"global"`
}

func (cfg Config) Validate() error {
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
