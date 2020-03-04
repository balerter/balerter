package config

import (
	"fmt"
	"strings"
)

type StorageS3 struct {
	Name     string `json:"name" yaml:"name"`
	Region   string `json:"region" yaml:"region"`
	Key      string `json:"key" yaml:"key"`
	Secret   string `json:"secret" yaml:"secret"`
	Endpoint string `json:"endpoint" yaml:"endpoint"`
	Bucket   string `json:"bucket" yaml:"bucket"`
}

func (cfg StorageS3) SetDefaults() {

}

func (cfg StorageS3) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	return nil
}
