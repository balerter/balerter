package channels

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/channels/alertmanager"
	"github.com/balerter/balerter/internal/config/channels/alertmanagerreceiver"
	"github.com/balerter/balerter/internal/config/channels/discord"
	"github.com/balerter/balerter/internal/config/channels/email"
	"github.com/balerter/balerter/internal/config/channels/notify"
	"github.com/balerter/balerter/internal/config/channels/slack"
	"github.com/balerter/balerter/internal/config/channels/syslog"
	"github.com/balerter/balerter/internal/config/channels/telegram"
	"github.com/balerter/balerter/internal/config/channels/webhook"
	"github.com/balerter/balerter/internal/util"
)

// Channels config for define channels
type Channels struct {
	// Email channel
	Email []email.Email `json:"email" yaml:"email"`
	// Slack channel
	Slack []slack.Slack `json:"slack" yaml:"slack"`
	// Telegram channel
	Telegram []telegram.Telegram `json:"telegram" yaml:"telegram" hcl:"telegram,block"`
	// Syslog channel
	Syslog []syslog.Syslog `json:"syslog" yaml:"syslog"`
	// Notify channel
	Notify []notify.Notify `json:"notify" yaml:"notify"`
	// Discord channel
	Discord []discord.Discord `json:"discord" yaml:"discord"`
	// Webhook channel
	Webhook []webhook.Webhook `json:"webhook" yaml:"webhook"`
	// Alertmanager channel
	Alertmanager []alertmanager.Alertmanager `json:"alertmanager" yaml:"alertmanager"`
	// AlertmanagerReceiver channel
	AlertmanagerReceiver []alertmanagerreceiver.AlertmanagerReceiver `json:"alertmanager_receiver" yaml:"alertmanager_receiver"`
}

// Validate config
func (cfg Channels) Validate() error {
	var names []string

	for _, c := range cfg.Email {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return fmt.Errorf("validate channel email: %w", err)
		}
	}
	if name := util.CheckUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for channels 'email': %s", name)
	}

	names = names[:0]
	for _, c := range cfg.Slack {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return fmt.Errorf("validate channel slack: %w", err)
		}
	}
	if name := util.CheckUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for channels 'slack': %s", name)
	}

	names = names[:0]
	for _, c := range cfg.Telegram {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return fmt.Errorf("validate channel telegram: %w", err)
		}
	}
	if name := util.CheckUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for channels 'telegram': %s", name)
	}

	names = names[:0]
	for _, c := range cfg.Syslog {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return fmt.Errorf("validate channel syslog: %w", err)
		}
	}
	if name := util.CheckUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for channels 'syslog': %s", name)
	}

	names = names[:0]
	for _, c := range cfg.Notify {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return fmt.Errorf("validate channel notify: %w", err)
		}
	}
	if name := util.CheckUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for channels 'notify': %s", name)
	}

	names = names[:0]
	for _, c := range cfg.Discord {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return fmt.Errorf("validate channel discord: %w", err)
		}
	}
	if name := util.CheckUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for channels 'discord': %s", name)
	}

	names = names[:0]
	for _, c := range cfg.Webhook {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return fmt.Errorf("validate channel webhook: %w", err)
		}
	}
	if name := util.CheckUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for channels 'webhook': %s", name)
	}

	names = names[:0]
	for _, c := range cfg.Alertmanager {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return fmt.Errorf("validate channel alertmanager: %w", err)
		}
	}
	if name := util.CheckUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for channels 'alertmanager': %s", name)
	}

	names = names[:0]
	for _, c := range cfg.AlertmanagerReceiver {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return fmt.Errorf("validate channel alertmanager_receiver: %w", err)
		}
	}
	if name := util.CheckUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for channels 'alertmanager_receiver': %s", name)
	}

	return nil
}
