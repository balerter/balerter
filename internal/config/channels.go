package config

import "fmt"

type Channels struct {
	Slack    []ChannelSlack    `json:"slack" yaml:"slack"`
	Telegram []ChannelTelegram `json:"telegram" yaml:"telegram"`
	Syslog   []ChannelSyslog   `json:"syslog" yaml:"syslog"`
}

func (cfg Channels) SetDefaults() {
	for _, c := range cfg.Slack {
		c.SetDefaults()
	}
	for _, c := range cfg.Telegram {
		c.SetDefaults()
	}
	for _, c := range cfg.Syslog {
		c.SetDefaults()
	}
}

func (cfg Channels) Validate() error {
	for _, c := range cfg.Slack {
		if err := c.Validate(); err != nil {
			return fmt.Errorf("validate channel slack: %w", err)
		}
	}

	for _, c := range cfg.Telegram {
		if err := c.Validate(); err != nil {
			return fmt.Errorf("validate channel telegram: %w", err)
		}
	}

	for _, c := range cfg.Syslog {
		if err := c.Validate(); err != nil {
			return fmt.Errorf("validate channel syslog: %w", err)
		}
	}

	return nil
}
