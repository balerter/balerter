package config

import "fmt"

type ScriptsSources struct {
	Folder []ScriptSourceFolder `json:"folder" yaml:"folder"`
}

func (cfg ScriptsSources) Validate() error {
	var names []string
	for _, c := range cfg.Folder {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return err
		}
	}

	if name := checkUnique(names); name != "" {
		return fmt.Errorf("found duplicated name: %s", name)
	}

	return nil
}
