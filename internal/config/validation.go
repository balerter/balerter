package config

import (
	"fmt"
	"io/ioutil"
)

func (cfg *Config) Validate() error {
	for _, c := range cfg.Scripts.Sources.Folder {
		if err := c.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (cfg *ScriptSourceFolder) Validate() error {
	_, err := ioutil.ReadDir(cfg.Path)
	if err != nil {
		return fmt.Errorf("error read folder '%s', %w", cfg.Path, err)
	}

	return nil
}
