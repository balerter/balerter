package email

import (
	"fmt"
	"strings"
)

// Email configures notifications via email.
type Email struct {
	Name     string `json:"name" yaml:"name"`
	From     string `json:"from" yaml:"from"`
	To       string `json:"to" yaml:"to"`
	Cc       string `json:"cc" yaml:"cc"`
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Secure   string `json:"secure" yaml:"secure"`
	Timeout  int    `json:"timeout" yaml:"timeout"`
}

// Validate checks the email configuration.
func (cfg *Email) Validate() error {
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
