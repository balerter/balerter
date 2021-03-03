package folder

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Folder struct {
	Name string `json:"name" yaml:"name"`
	Path string `json:"path" yaml:"path"`
	Mask string `json:"mask" yaml:"mask"`
}

func (cfg Folder) Validate() error {
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
