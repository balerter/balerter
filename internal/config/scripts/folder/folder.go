package folder

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// Folder script source config
type Folder struct {
	// Name of the script source
	Name string `json:"name" yaml:"name" hcl:"name,label"`
	// Path to the scripts folder
	Path string `json:"path" yaml:"path" hcl:"path"`
	// Mask for script matching
	Mask string `json:"mask" yaml:"mask" hcl:"mask"`
}

// Validate config
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
