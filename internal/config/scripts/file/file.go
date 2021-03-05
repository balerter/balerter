package file

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// File script source config
type File struct {
	// Name of the script source
	Name string `json:"name" yaml:"name" hcl:"name,label"`
	// Filename of the script
	Filename string `json:"filename" yaml:"filename" hcl:"filename"`
}

// Validate config
func (cfg File) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}
	if strings.TrimSpace(cfg.Filename) == "" {
		return fmt.Errorf("filename must be not empty")
	}

	_, err := ioutil.ReadFile(cfg.Filename)
	if err != nil {
		return fmt.Errorf("error read file '%s', %w", cfg.Filename, err)
	}

	return nil
}
