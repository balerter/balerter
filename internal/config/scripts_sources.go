package config

type ScriptsSources struct {
	Folder []ScriptSourceFolder `json:"folder" yaml:"folder"`
}

func (cfg ScriptsSources) SetDefaults() {
	for _, c := range cfg.Folder {
		c.SetDefaults()
	}
}

func (cfg ScriptsSources) Validate() error {
	for _, c := range cfg.Folder {
		if err := c.Validate(); err != nil {
			return err
		}
	}

	return nil
}
