package sources

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/scripts/sources/file"
	"github.com/balerter/balerter/internal/config/scripts/sources/folder"
	"github.com/balerter/balerter/internal/util"
)

type Sources struct {
	Folder []*folder.Folder `json:"folder" yaml:"folder"`
	File   []*file.File     `json:"file" yaml:"file"`
}

func (cfg Sources) Validate() error {
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
