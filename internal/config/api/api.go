package api

type API struct {
	Address        string `json:"address" yaml:"address" hcl:"address,optional"`
	ServiceAddress string `json:"serviceAddress" yaml:"serviceAddress" hcl:"serviceAddress,optional"`
}

func (cfg API) Validate() error {
	return nil
}
