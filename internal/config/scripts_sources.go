package config

import "fmt"

type ScriptsSources struct {
	Folder []*ScriptSourceFolder `json:"folder" yaml:"folder"`
	File   []*ScriptSourceFile   `json:"file" yaml:"file"`
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
		return fmt.Errorf("found duplicated name for scritsource 'folder': %s", name)
	}

	names = names[:0]
	for _, c := range cfg.File {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return err
		}
	}
	if name := checkUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for scritsource 'file': %s", name)
	}

	return nil
}
