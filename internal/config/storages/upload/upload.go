package upload

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/storages/upload/s3"
	"github.com/balerter/balerter/internal/util"
)

// Upload storage config
type Upload struct {
	S3 []s3.S3 `json:"s3" yaml:"s3" hcl:"s3,block"`
}

// Validate config
func (cfg Upload) Validate() error {
	var names []string

	for _, c := range cfg.S3 {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return err
		}
	}
	if name := util.CheckUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for upload storages 's3': %s", name)
	}

	return nil
}
