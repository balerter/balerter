package config

type Global struct {
	SendStartNotification []string        `json:"send_start_notification" yaml:"send_start_notification"`
	SendStopNotification  []string        `json:"send_stop_notification" yaml:"send_stop_notification"`
	API                   *API            `json:"api" yaml:"api"`
	Storages              *GlobalStorages `json:"storages" yaml:"storages"`
}

func (cfg Global) SetDefaults() {
	if cfg.API == nil {
		cfg.API = &API{}
	}
	cfg.API.SetDefaults()
	if cfg.Storages == nil {
		cfg.Storages = &GlobalStorages{}
	}
	cfg.Storages.SetDefaults()
}

func (cfg Global) Validate() error {
	if err := cfg.API.Validate(); err != nil {
		return err
	}

	return cfg.Storages.Validate()
}
