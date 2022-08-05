package api

type CoreAPI struct {
	Address   string `json:"address" yaml:"address" hcl:"address,optional"`
	AuthToken string `json:"authToken" yaml:"authToken" hcl:"authToken,optional"`
}

// API config
type API struct {
	// Address is address for handle API
	Address string `json:"address" yaml:"address" hcl:"address,optional"`
	// ServiceAddress is address for listen service handlers
	ServiceAddress string `json:"serviceAddress" yaml:"serviceAddress" hcl:"serviceAddress,optional"`

	CoreApi *CoreAPI `json:"coreApi" yaml:"coreApi" hcl:"coreApi,block"`
}

// Validate config
func (cfg API) Validate() error {
	return nil
}
