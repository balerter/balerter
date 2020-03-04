package config

import (
	"fmt"
	"strings"
)

var (
	syslogSeverity = []string{"EMERG", "ALERT", "CRIT", "ERR", "WARNING", "NOTICE", "INFO", "DEBUG"}
	syslogFacility = []string{"KERN", "USER", "MAIL", "DAEMON", "AUTH", "SYSLOG", "LPR", "NEWS", "UUCP", "CRON", "AUTHPRIV", "FTP", "LOCAL0", "LOCAL1", "LOCAL2", "LOCAL3", "LOCAL4", "LOCAL5", "LOCAL6", "LOCAL7"}
)

type ChannelSyslogPriority struct {
	Alert   string `json:"alert" yaml:"alert"`
	Warning string `json:"warning" yaml:"warning"`
	Success string `json:"success" yaml:"success"`
}

type ChannelSyslog struct {
	Name     string                `json:"name" yaml:"name"`
	Tag      string                `json:"tag" yaml:"tag"`
	Network  string                `json:"network" yaml:"network"`
	Address  string                `json:"address" yaml:"address"`
	Priority ChannelSyslogPriority `json:"priority" yaml:"priority"`
}

func (cfg ChannelSyslog) SetDefaults() {

}

func (cfg ChannelSyslog) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	if strings.ToLower(cfg.Network) != "tcp" && strings.ToLower(cfg.Network) != "udp" && cfg.Network != "" {
		return fmt.Errorf("corrent values for 'network': 'tcp', 'udp' or empty value")
	}

	if strings.TrimSpace(cfg.Address) == "" {
		return fmt.Errorf("address must be not empty")
	}

	if err := cfg.Priority.Validate(); err != nil {
		return err
	}

	return nil
}

func (cfg ChannelSyslogPriority) Validate() error {

	if err := validatePriority(cfg.Alert); err != nil {
		return fmt.Errorf("error validate priority: %w", err)
	}

	if err := validatePriority(cfg.Warning); err != nil {
		return fmt.Errorf("error validate priority: %w", err)
	}

	if err := validatePriority(cfg.Success); err != nil {
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

	if !inArray(parts[0], syslogSeverity) {
		return fmt.Errorf("bad priority format")
	}

	if len(parts) == 1 {
		return nil
	}

	if !inArray(parts[1], syslogFacility) {
		return fmt.Errorf("bad priority format")
	}

	return nil
}
