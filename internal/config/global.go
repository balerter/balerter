package config

type Global struct {
	SendStartNotification []string `json:"send_start_notification" yaml:"send_start_notification"`
	SendStopNotification  []string `json:"send_stop_notification" yaml:"send_stop_notification"`
	API                   API      `json:"api" yaml:"api"`
}

func (cfg *Global) SetDefaults() {
	cfg.API.SetDefaults()
}

func (cfg *Global) Validate() error {
	return cfg.API.Validate()
}
