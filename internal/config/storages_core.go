package config

type StoragesCore struct {
	File []StorageCoreFile `json:"file" yaml:"file"`
}

func (cfg StoragesCore) SetDefaults() {
	for _, c := range cfg.File {
		c.SetDefaults()
	}
}

func (cfg StoragesCore) Validate() error {
	for _, c := range cfg.File {
		if err := c.Validate(); err != nil {
			return err
		}
	}

	return nil
}
