package config

import (
	"flag"
	"time"
)

const (
	defaultScriptsUpdateInterval time.Duration = time.Second * 60
)

var (
	configSource = flag.String("config", "config.yml", "Configuration source. Currently supports only path to yaml file.")
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
	Scripts Scripts `json:"scripts" yaml:"scripts"`
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
