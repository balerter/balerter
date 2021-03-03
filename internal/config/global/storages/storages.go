package storages

type Storages struct {
	KV    string `json:"kv" yaml:"kv"`
	Alert string `json:"alert" yaml:"alert"`
}

func (cfg Storages) Validate() error {
	return nil
}
