package config

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
)

var StdIn io.Reader = os.Stdin

func New(configSource string) (*Config, error) {
	cfg := &Config{}

	loader := aconfig.LoaderFor(cfg, aconfig.Config{
		SkipEnv:   true,
		SkipFlags: true,

		FailOnFileNotFound: true,
		Files:              []string{configSource},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
		},
	})

	if err := loader.Load(); err != nil {
		return nil, fmt.Errorf("error read config file, %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("error config validation, %w", err)
	}

	return cfg, nil
}

type Config struct {
	Scripts     *Scripts     `json:"scripts" yaml:"scripts"`
	DataSources *DataSources `json:"datasources" yaml:"datasources"`
	Channels    *Channels    `json:"channels" yaml:"channels"`
	Storages    *Storages    `json:"storages" yaml:"storages"`
	Global      *Global      `json:"global" yaml:"global"`
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
