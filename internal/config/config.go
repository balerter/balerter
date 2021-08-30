package config

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/balerter/balerter/internal/config/secrets/env"
	"github.com/balerter/balerter/internal/config/secrets/vault"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/balerter/balerter/internal/config/api"
	"github.com/balerter/balerter/internal/config/channels"
	"github.com/balerter/balerter/internal/config/datasources"
	"github.com/balerter/balerter/internal/config/scripts"
	"github.com/balerter/balerter/internal/config/storages/core"
	"github.com/balerter/balerter/internal/config/storages/upload"
	"github.com/balerter/balerter/internal/config/system"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"gopkg.in/yaml.v2"
)

// StdIn is default stdin reader
var StdIn io.Reader = os.Stdin

// New creates new config instance
func New(fs *flag.FlagSet, args []string) (*Config, *Flags, error) {
	cfg := &Config{}

	var data []byte
	var err error

	flg := &Flags{}

	fs.StringVar(&flg.ConfigFilePath, "config", "config.yml", "configuration source. Currently supports only path to yaml file and 'stdin'.")
	fs.StringVar(&flg.LogLevel, "logLevel", "INFO", "log level. ERROR, INFO or DEBUG")
	fs.BoolVar(&flg.Debug, "debug", false, "debug mode")
	fs.BoolVar(&flg.Once, "once", false, "once run scripts and exit")
	fs.StringVar(&flg.Script, "script", "", "ignore all script sources and runs only one script. Meta-tag @ignore will be ignored")
	fs.BoolVar(&flg.AsJSON, "json", false, "output json format")
	err = fs.Parse(args)
	if err != nil {
		return nil, nil, fmt.Errorf("error parse falgs, %w", err)
	}

	if flg.ConfigFilePath == "stdin" {
		data, err = ioutil.ReadAll(StdIn)
	} else {
		data, err = ioutil.ReadFile(flg.ConfigFilePath)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("error read config file, %w", err)
	}

	var errDecodeSecrets error

	data, errDecodeSecrets = decodeSecrets(data)
	if errDecodeSecrets != nil {
		return nil, nil, fmt.Errorf("error decode secrets, %w", errDecodeSecrets)
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

// Flags represents CLI flags
type Flags struct {
	// ConfigFilePath for CLI flag '-script'
	ConfigFilePath string
	// LogLevel for CLI flag '-logLevel'
	LogLevel string
	// Debug for CLI flag '-debug'
	Debug bool
	// Once for CLI flag '-once'
	Once bool
	// Script for CLI flag '-script'
	Script string
	// for CLI flag '-asJson' for test tool
	AsJSON bool
}

// Config represent balerter configuration
type Config struct {
	// Scripts section for define script sources
	Scripts *scripts.Scripts `json:"scripts" yaml:"scripts" hcl:"scripts,block"`
	// DataSources section for define data sources
	DataSources *datasources.DataSources `json:"datasources" yaml:"datasources" hcl:"datasources,block"`
	// Channels section for define channels
	Channels *channels.Channels `json:"channels" yaml:"channels" hcl:"channels,block"`
	// StoragesUpload section for define upload storages
	StoragesUpload *upload.Upload `json:"storagesUpload" yaml:"storagesUpload" hcl:"storagesUpload,block"`
	// StoragesCore section for define core storages
	StoragesCore *core.Core `json:"storagesCore" yaml:"storagesCore" hcl:"storagesCore,block"`
	// API section for define API settings
	API *api.API `json:"api" yaml:"api" hcl:"api,block"`

	// LuaModulesPath for path to lua modules
	LuaModulesPath string `json:"luaModulesPath" yaml:"luaModulesPath" hcl:"luaModulesPath,optional"`
	// StorageAlert item
	StorageAlert string `json:"storageAlert" yaml:"storageAlert" hcl:"storageAlert,optional"`
	// StorageKV item
	StorageKV string `json:"storageKV" yaml:"storageKV" hcl:"storageKV,optional"`

	System *system.System `json:"system" yaml:"system" hcl:"system,block"`
}

// Validate the config
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
	if cfg.System != nil {
		if err := cfg.System.Validate(); err != nil {
			return fmt.Errorf("error system validation, %w", err)
		}
	}

	return nil
}

var (
	reSecrets = regexp.MustCompile(`{secret:(vault|env):((?U).+)}`)
)

func decodeSecrets(data []byte) ([]byte, error) {
	secrets := reSecrets.FindAllSubmatch(data, -1)

	for _, secret := range secrets {
		if len(secret) != 3 {
			return nil, fmt.Errorf("unexpected secret RE submatch len = %d", len(secret))
		}

		var v []byte
		var err error

		switch string(secret[1]) {
		case "env":
			v, err = env.DecodeSecret(secret[2])
			if err != nil {
				return nil, err
			}
		case "vault":
			v, err = vault.DecodeSecret(secret[2])
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unsupported secret engine = %s", secret[1])
		}

		data = bytes.Replace(data, secret[0], v, 1)
	}

	return data, nil
}
