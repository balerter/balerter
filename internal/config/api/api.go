package api

// API config
type API struct {
	// Address is address for handle API
	Address string `json:"address" yaml:"address" hcl:"address,optional"`
	// ServiceAddress is address for listen service handlers
	ServiceAddress string `json:"serviceAddress" yaml:"serviceAddress" hcl:"serviceAddress,optional"`
}

func (cfg API) Validate() error {
	return nil
}
