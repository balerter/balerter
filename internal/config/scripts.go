package config

type Scripts struct {
	Sources ScriptsSources `json:"sources" yaml:"sources"`
}

func (cfg Scripts) SetDefaults() {
	cfg.Sources.SetDefaults()
}

func (cfg Scripts) Validate() error {
	return cfg.Sources.Validate()
}
