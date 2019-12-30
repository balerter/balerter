package config

import "flag"

var (
	configSource = flag.String("config", "config.yml", "Configuration source. Currently support only path to yaml file.")
)

type Config struct {
}

func New() *Config {
	cfg := &Config{}

	return cfg
}
