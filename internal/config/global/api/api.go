package api

type API struct {
	Address string `json:"address" yaml:"address"`
}

func (cfg API) Validate() error {
	return nil
}
