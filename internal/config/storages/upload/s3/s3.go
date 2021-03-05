package s3

import (
	"fmt"
	"strings"
)

// S3 upload storage config
type S3 struct {
	// Name of the upload storage
	Name string `json:"name" yaml:"name" hcl:"name,label"`
	// Region value
	Region string `json:"region" yaml:"region" hcl:"region"`
	// Key value
	Key string `json:"key" yaml:"key" hcl:"key"`
	// Secret value
	Secret string `json:"secret" yaml:"secret" hcl:"secret"`
	// Endpoint value
	Endpoint string `json:"endpoint" yaml:"endpoint" hcl:"endpoint"`
	// Bucket value
	Bucket string `json:"bucket" yaml:"bucket" hcl:"bucket"`
}

// Validate config
func (cfg *S3) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	return nil
}
