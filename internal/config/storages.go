package config

type Storages struct {
	Upload StoragesUpload `json:"upload" yaml:"upload"`
	Core   StoragesCore   `json:"core" yaml:"core"`
}

func (cfg Storages) SetDefaults() {
	cfg.Upload.SetDefaults()
	cfg.Core.SetDefaults()
}

func (cfg Storages) Validate() error {
	if err := cfg.Upload.Validate(); err != nil {
		return err
	}

	if err := cfg.Core.Validate(); err != nil {
		return err
	}

	return nil
}
