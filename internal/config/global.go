package config

type Global struct {
	SendStartNotification []string       `json:"sendStartNotification" yaml:"sendStartNotification"`
	SendStopNotification  []string       `json:"sendStopNotification" yaml:"sendStopNotification"`
	API                   API            `json:"api" yaml:"api"`
	Storages              GlobalStorages `json:"storages" yaml:"storages"`
}

func (cfg Global) Validate() error {
	if err := cfg.API.Validate(); err != nil {
		return err
	}

	return cfg.Storages.Validate()
}
