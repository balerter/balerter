package storages

import (
	"github.com/balerter/balerter/internal/config/storages/core"
	"github.com/balerter/balerter/internal/config/storages/upload"
)

type Storages struct {
	Upload upload.Upload `json:"upload" yaml:"upload"`
	Core   core.Core     `json:"core" yaml:"core"`
}

func (cfg *Storages) Validate() error {
	if err := cfg.Upload.Validate(); err != nil {
		return err
	}
	if err := cfg.Core.Validate(); err != nil {
		return err
	}

	return nil
}
