package syslog

import (
	"fmt"
	"github.com/balerter/balerter/internal/util"
	"strings"
)

var (
	syslogSeverity = []string{"EMERG", "ALERT", "CRIT", "ERR", "WARNING", "NOTICE", "INFO", "DEBUG"}
	syslogFacility = []string{"KERN", "USER", "MAIL", "DAEMON", "AUTH", "SYSLOG",
		"LPR", "NEWS", "UUCP", "CRON", "AUTHPRIV", "FTP", "LOCAL0", "LOCAL1",
		"LOCAL2", "LOCAL3", "LOCAL4", "LOCAL5", "LOCAL6", "LOCAL7"}
)

// Syslog channel config
type Syslog struct {
	// Name of the channel
	Name string `json:"name" yaml:"name"`
	// Tag describe syslog tag
	Tag string `json:"tag" yaml:"tag"`
	// Network value
	Network string `json:"network" yaml:"network"`
	// Address value
	Address string `json:"address" yaml:"address"`
	// Priority value, Severity+Facility
	Priority string `json:"priority" yaml:"priority"`
}

// Validate config
func (cfg Syslog) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	if !strings.EqualFold(cfg.Network, "tcp") && !strings.EqualFold(cfg.Network, "udp") && cfg.Network != "" {
		return fmt.Errorf("corrent values for 'network': 'tcp', 'udp' or empty value")
	}

	if err := validatePriority(cfg.Priority); err != nil {
		return fmt.Errorf("error validate priority: %w", err)
	}

	return nil
}

func validatePriority(p string) error {
	if p == "" {
		return nil
	}

	parts := strings.Split(p, "|")
	if len(parts) > 2 {
		return fmt.Errorf("bad priority format")
	}

	if !util.InArray(parts[0], syslogSeverity) {
		return fmt.Errorf("bad priority format")
	}

	if len(parts) == 1 {
		return nil
	}

	if !util.InArray(parts[1], syslogFacility) {
		return fmt.Errorf("bad priority format")
	}

	return nil
}
