package config

import (
	"fmt"
	"strings"
)

type GlobalStorages struct {
	KV    string `json:"kv" yaml:"kv"`
	Alert string `json:"alert" yaml:"alert"`
}

func (cfg *GlobalStorages) Validate() error {
	if strings.TrimSpace(cfg.KV) == "" {
		return fmt.Errorf("storages.kv must be not empty")
	}

	if strings.TrimSpace(cfg.Alert) == "" {
		return fmt.Errorf("storages.alert must be not empty")
	}

	return nil
}
