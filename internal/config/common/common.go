package common

type BasicAuth struct {
	Username string `json:"username" yaml:"username" hcl:"username"`
	Password string `json:"password" yaml:"password" hcl:"password"`
}
