package config

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

type ScriptSourceFolder struct {
	UpdateInterval time.Duration `json:"update_interval" yaml:"update_interval"`
	Name           string        `json:"name" yaml:"name"`
	Path           string        `json:"path" yaml:"path"`
	Mask           string        `json:"mask" yaml:"mask"`
}

func (cfg *ScriptSourceFolder) SetDefaults() {
	cfg.Mask = "*.lua"
	cfg.UpdateInterval = time.Second * 60
}

func (cfg *ScriptSourceFolder) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}
	if strings.TrimSpace(cfg.Path) == "" {
		return fmt.Errorf("path must be not empty")
	}

	_, err := ioutil.ReadDir(cfg.Path)
	if err != nil {
		return fmt.Errorf("error read folder '%s', %w", cfg.Path, err)
	}

	return nil
}
