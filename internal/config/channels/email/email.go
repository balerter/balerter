package email

import (
	"fmt"
	"strings"
)

// Email configures notifications via email.
type Email struct {
	// Name of the channel
	Name string `json:"name" yaml:"name" hcl:"name,label"`
	// From field
	From string `json:"from" yaml:"from" hcl:"from"`
	// To field
	To string `json:"to" yaml:"to" hcl:"to"`
	// Cc field
	Cc string `json:"cc" yaml:"cc" hcl:"cc,optional"`
	// Host value
	Host string `json:"host" yaml:"host" hcl:"host"`
	// Port value
	Port string `json:"port" yaml:"port" hcl:"port"`
	// Username value
	Username string `json:"username" yaml:"username" hcl:"username"`
	// Password value
	Password string `json:"password" yaml:"password" hcl:"password"`
	// Secure value
	Secure string `json:"secure" yaml:"secure" hcl:"secure"`
	// Timeout value
	Timeout int  `json:"timeout" yaml:"timeout" hcl:"timeout,optional"`
	Ignore  bool `json:"ignore" yaml:"ignore" hcl:"ignore,optional"`
}

// Validate checks the email configuration.
func (cfg Email) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}
	if strings.TrimSpace(cfg.From) == "" {
		return fmt.Errorf("from must be not empty")
	}
	if strings.TrimSpace(cfg.To) == "" {
		return fmt.Errorf("to must be not empty")
	}
	if strings.TrimSpace(cfg.Host) == "" {
		return fmt.Errorf("host must be not empty")
	}
	if strings.TrimSpace(cfg.Port) == "" {
		return fmt.Errorf("port must be not empty")
	}
	s := strings.TrimSpace(strings.ToLower(cfg.Secure))
	if s != "none" && s != "ssl" && s != "tls" && s != "" {
		return fmt.Errorf("secure must be set to none, ssl or tls")
	}

	return nil
}
