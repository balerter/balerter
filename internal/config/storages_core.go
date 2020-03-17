package config

import "fmt"

type StoragesCore struct {
	File []StorageCoreFile `json:"file" yaml:"file"`
}

func (cfg StoragesCore) Validate() error {
	var names []string

	for _, c := range cfg.File {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return err
		}
	}
	if name := checkUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for core storages 'file': %s", name)
	}

	return nil
}
