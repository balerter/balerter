package config

import (
	"flag"
	"fmt"
	"github.com/balerter/balerter/internal/config/api"
	"github.com/balerter/balerter/internal/config/channels"
	"github.com/balerter/balerter/internal/config/datasources"
	"github.com/balerter/balerter/internal/config/scripts"
	"github.com/balerter/balerter/internal/config/storages/core"
	"github.com/balerter/balerter/internal/config/storages/upload"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var StdIn io.Reader = os.Stdin

func New() (*Config, *Flags, error) {
	cfg := &Config{
		//Scripts:     &scripts.Scripts{},
		//DataSources: &datasources.DataSources{},
		//Channels:    &channels.Channels{},
		//Storages:    &storages.Storages{},
		//Global:      &global.Global{},
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

	if err := decodeCfg(flg.ConfigFilePath, data, cfg); err != nil {
		return nil, nil, fmt.Errorf("error parse config file, %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, nil, fmt.Errorf("error config validation, %w", err)
	}

	return cfg, flg, nil
}

func decodeCfg(filename string, data []byte, cfg *Config) error {
	if strings.HasSuffix(filename, ".yml") || strings.HasSuffix(filename, ".yaml") {
		return yaml.Unmarshal(data, cfg)
	}
	if strings.HasSuffix(filename, ".hcl") {
		return hclsimple.Decode(filename, data, nil, cfg)
	}

	return fmt.Errorf("unknown format")
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
	Scripts        *scripts.Scripts         `json:"scripts" yaml:"scripts" hcl:"scripts,block"`
	DataSources    *datasources.DataSources `json:"datasources" yaml:"datasources" hcl:"datasources,block"`
	Channels       *channels.Channels       `json:"channels" yaml:"channels" hcl:"channels,block"`
	StoragesUpload *upload.Upload           `json:"storagesUpload" yaml:"storagesUpload" hcl:"storagesUpload,block"`
	StoragesCore   *core.Core               `json:"storagesCore" yaml:"storagesCore" hcl:"storagesCore,block"`
	API            *api.API                 `json:"api" yaml:"api" hcl:"api,block"`

	LuaModulesPath string `json:"luaModulesPath" yaml:"luaModulesPath" hcl:"luaModulesPath,optional"`
	StorageAlert   string `json:"storageAlert" yaml:"storageAlert" hcl:"storageAlert,optional"`
	StorageKV      string `json:"storageKV" yaml:"storageKV" hcl:"storageKV,optional"`
}

func (cfg Config) Validate() error {
	if cfg.Scripts != nil {
		if err := cfg.Scripts.Validate(); err != nil {
			return fmt.Errorf("error Scripts validation, %w", err)
		}
	}
	if cfg.DataSources != nil {
		if err := cfg.DataSources.Validate(); err != nil {
			return fmt.Errorf("error DataSources validation, %w", err)
		}
	}
	if cfg.Channels != nil {
		if err := cfg.Channels.Validate(); err != nil {
			return fmt.Errorf("error Channels validation, %w", err)
		}
	}
	if cfg.StoragesUpload != nil {
		if err := cfg.StoragesUpload.Validate(); err != nil {
			return fmt.Errorf("error StoragesUpload validation, %w", err)
		}
	}
	if cfg.StoragesCore != nil {
		if err := cfg.StoragesCore.Validate(); err != nil {
			return fmt.Errorf("error StoragesCore validation, %w", err)
		}
	}
	if cfg.API != nil {
		if err := cfg.API.Validate(); err != nil {
			return fmt.Errorf("error api validation, %w", err)
		}
	}

	return nil
}
