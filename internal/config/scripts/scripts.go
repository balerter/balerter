package scripts

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/scripts/file"
	"github.com/balerter/balerter/internal/config/scripts/folder"
	"github.com/balerter/balerter/internal/util"
)

type Scripts struct {
	UpdateInterval int             `json:"updateInterval" yaml:"updateInterval" hcl:"updateInterval,optional"`
	Folder         []folder.Folder `json:"folder" yaml:"folder" hcl:"folder,block"`
	File           []file.File     `json:"file" yaml:"file" hcl:"file,block"`
}

func (cfg Scripts) Validate() error {
	if cfg.UpdateInterval < 0 {
		return fmt.Errorf("updateInterval must be not less than 0")
	}

	var names []string
	for _, c := range cfg.Folder {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return err
		}
	}
	if name := util.CheckUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for scritsource 'folder': %s", name)
	}

	names = names[:0]
	for _, c := range cfg.File {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return err
		}
	}
	if name := util.CheckUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for scritsource 'file': %s", name)
	}

	return nil
}
