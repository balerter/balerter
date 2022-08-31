package cloud

import (
	"fmt"
)

type Cloud struct {
	Token   string `json:"token" yaml:"token" hcl:"token"`
	Server  string `json:"server" yaml:"server" hcl:"server,optional"`
	Timeout int    `json:"timeout" yaml:"timeout" hcl:"timeout,optional"`
}

func (s *Cloud) Validate() error {
	if s.Token == "" {
		return fmt.Errorf("token must be not empty")
	}
	return nil
}
