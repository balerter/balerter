package config

type GlobalStorages struct {
	KV    string `json:"kv" yaml:"kv"`
	Alert string `json:"alert" yaml:"alert"`
}

func (cfg *GlobalStorages) Validate() error {
	return nil
}
