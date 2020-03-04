package config

type API struct {
	Address string `json:"address" yaml:"address"`
	Metrics bool   `json:"metrics" yaml:"metrics"`
}

func (cfg *API) SetDefaults() {
	if cfg.Address == "" {
		cfg.Address = defaultAPIAddress
	}
}

func (cfg *API) Validate() error {
	return nil
}
