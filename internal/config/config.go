package config

import "flag"

var (
	configSource = flag.String("config", "config.yml", "Configuration source. Currently support only path to yaml file.")
)

func New() *Config {
	cfg := &Config{}

	return cfg
}

type Config struct {
	Scripts Scripts `json:"scripts" yaml:"scripts"`
}

type Scripts struct {
	Sources ScriptsSources `json:"sources" yaml:"sources"`
}

type ScriptsSources struct {
	Folder []ScriptSourceFolder `json:"folder" yaml:"folder"`
}

type ScriptSourceFolder struct {
	Name string `json:"name" yaml:"name"`
	Path string `json:"path" yaml:"path"`
	Mask string `json:"mask" yaml:"mask"`
}
