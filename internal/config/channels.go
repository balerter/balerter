package config

import "fmt"

type Channels struct {
	Email    []*ChannelEmail    `json:"email" yaml:"email"`
	Slack    []*ChannelSlack    `json:"slack" yaml:"slack"`
	Telegram []*ChannelTelegram `json:"telegram" yaml:"telegram"`
	Syslog   []*ChannelSyslog   `json:"syslog" yaml:"syslog"`
	Notify   []*ChannelNotify   `json:"notify" yaml:"notify"`
	Discord  []*ChannelDiscord  `json:"discord" yaml:"discord"`
	Webhook  []*ChannelWebhook  `json:"webhook" yaml:"webhook"`
	Alertmanager []*ChannelAlertmanager `json:"alertmanager" yaml:"alertmanager"`
}

func (cfg *Channels) Validate() error { //nolint:gocyclo // Validate calls only once on application start
	var names []string

	for _, c := range cfg.Email {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return fmt.Errorf("validate channel email: %w", err)
		}
	}
	if name := checkUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for channels 'email': %s", name)
	}

	names = names[:0]
	for _, c := range cfg.Slack {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return fmt.Errorf("validate channel slack: %w", err)
		}
	}
	if name := checkUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for channels 'slack': %s", name)
	}

	names = names[:0]
	for _, c := range cfg.Telegram {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return fmt.Errorf("validate channel telegram: %w", err)
		}
	}
	if name := checkUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for channels 'telegram': %s", name)
	}

	names = names[:0]
	for _, c := range cfg.Syslog {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return fmt.Errorf("validate channel syslog: %w", err)
		}
	}
	if name := checkUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for channels 'syslog': %s", name)
	}

	names = names[:0]
	for _, c := range cfg.Notify {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return fmt.Errorf("validate channel notify: %w", err)
		}
	}
	if name := checkUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for channels 'notify': %s", name)
	}

	names = names[:0]
	for _, c := range cfg.Discord {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return fmt.Errorf("validate channel discord: %w", err)
		}
	}
	if name := checkUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for channels 'discord': %s", name)
	}

	names = names[:0]
	for _, c := range cfg.Webhook {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return fmt.Errorf("validate channel webhook: %w", err)
		}
	}
	if name := checkUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for channels 'webhook': %s", name)
  }
  
	names = names[:0]
  for _, c := range cfg.Alertmanager {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return fmt.Errorf("validate channel alertmanager: %w", err)
		}
	}
	if name := checkUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for channels 'alertmanager': %s", name)
	}

	return nil
}
