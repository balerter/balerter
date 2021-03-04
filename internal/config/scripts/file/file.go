package file

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type File struct {
	Name     string `json:"name" yaml:"name" hcl:"name"`
	Filename string `json:"filename" yaml:"filename" hcl:"filename"`
}

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
