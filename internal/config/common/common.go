package common

// BasicAuth use for prometheus and loki configurations
type BasicAuth struct {
	// Username for basic auth
	Username string `json:"username" yaml:"username" hcl:"username"`
	// Password for basic auth
	Password string `json:"password" yaml:"password" hcl:"password"`
}
