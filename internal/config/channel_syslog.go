package config

import (
	"fmt"
	"strings"
)

func (cfg *ChannelSyslog) SetDefaults() {
	cfg.Severities.Alert = "EMERG"
	cfg.Severities.Warning = "EMERG"
	cfg.Severities.Success = "EMERG"
	cfg.Facilities.Alert = "SYSLOG"
	cfg.Facilities.Warning = "SYSLOG"
	cfg.Facilities.Success = "SYSLOG"
}

type ChannelSyslogSeverities struct {
	Alert   string `json:"alert" yaml:"alert"`
	Warning string `json:"warning" yaml:"warning"`
	Success string `json:"success" yaml:"success"`
}

type ChannelSyslogFacilities struct {
	Alert   string `json:"alert" yaml:"alert"`
	Warning string `json:"warning" yaml:"warning"`
	Success string `json:"success" yaml:"success"`
}

type ChannelSyslog struct {
	Name       string                  `json:"name" yaml:"name"`
	Tag        string                  `json:"tag" yaml:"tag"`
	Network    string                  `json:"network" yaml:"network"`
	Address    string                  `json:"address" yaml:"address"`
	Severities ChannelSyslogSeverities `json:"severities" yaml:"severities"`
	Facilities ChannelSyslogFacilities `json:"facilities" yaml:"facilities"`
}

func (cfg *ChannelSyslog) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	if strings.ToLower(cfg.Network) != "tcp" && strings.ToLower(cfg.Network) != "udp" && cfg.Network != "" {
		return fmt.Errorf("corrent values for 'network': 'tcp', 'udp' or empty value")
	}

	if strings.TrimSpace(cfg.Address) == "" {
		return fmt.Errorf("address must be not empty")
	}

	if err := cfg.Severities.Validate(); err != nil {
		return err
	}

	if err := cfg.Facilities.Validate(); err != nil {
		return err
	}

	return nil
}

func (cfg *ChannelSyslogSeverities) Validate() error {
	allowed := []string{"EMERG", "ALERT", "CRIT", "ERR", "WARNING", "NOTICE", "INFO", "DEBUG"}

	if !inArray(cfg.Alert, allowed) {
		return fmt.Errorf("Severities.Alert must be one of: %v", allowed)
	}

	if !inArray(cfg.Warning, allowed) {
		return fmt.Errorf("Severities.Warning must be one of: %v", allowed)
	}

	if !inArray(cfg.Success, allowed) {
		return fmt.Errorf("Severities.Success must be one of: %v", allowed)
	}

	return nil
}

func (cfg *ChannelSyslogFacilities) Validate() error {
	allowed := []string{"KERN", "USER", "MAIL", "DAEMON", "AUTH", "SYSLOG", "LPR", "NEWS", "UUCP", "CRON", "AUTHPRIV", "FTP", "LOCAL0", "LOCAL1", "LOCAL2", "LOCAL3", "LOCAL4", "LOCAL5", "LOCAL6", "LOCAL7"}

	if !inArray(cfg.Alert, allowed) {
		return fmt.Errorf("Facilities.Alert must be one of: %v", allowed)
	}

	if !inArray(cfg.Warning, allowed) {
		return fmt.Errorf("Facilities.Warning must be one of: %v", allowed)
	}

	if !inArray(cfg.Success, allowed) {
		return fmt.Errorf("Facilities.Success must be one of: %v", allowed)
	}

	return nil
}
