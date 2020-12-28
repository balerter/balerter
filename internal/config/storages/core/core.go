package core

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/storages/core/file"
	"github.com/balerter/balerter/internal/util"
)

type Core struct {
	File []*file.File `json:"file" yaml:"file"`
}

func (cfg Core) Validate() error {
	var names []string

	for _, c := range cfg.File {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return err
		}
	}
	if name := util.CheckUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for core storages 'file': %s", name)
	}

	return nil
}
