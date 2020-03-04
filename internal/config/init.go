package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func (cfg *Config) Init(configSource string) error {
	body, err := ioutil.ReadFile(configSource)
	if err != nil {
		return err
	}

	cfg.SetDefaults()

	err = yaml.Unmarshal(body, cfg)
	if err != nil {
		return err
	}

	return cfg.Validate()
}
