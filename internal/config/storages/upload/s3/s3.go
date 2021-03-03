package s3

import (
	"fmt"
	"strings"
)

type S3 struct {
	Name     string `json:"name" yaml:"name" hcl:"name,label"`
	Region   string `json:"region" yaml:"region" hcl:"region"`
	Key      string `json:"key" yaml:"key" hcl:"key"`
	Secret   string `json:"secret" yaml:"secret" hcl:"secret"`
	Endpoint string `json:"endpoint" yaml:"endpoint" hcl:"endpoint"`
	Bucket   string `json:"bucket" yaml:"bucket" hcl:"bucket"`
}

func (cfg *S3) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	return nil
}
