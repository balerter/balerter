package manager

import (
	"fmt"
	"github.com/balerter/balerter/internal/channels/alertmanager"
	alertmanagerreceiver "github.com/balerter/balerter/internal/channels/alertmanager_receiver"
	"github.com/balerter/balerter/internal/channels/twilio"
	"github.com/balerter/balerter/internal/channels/webhook"
	"github.com/balerter/balerter/internal/config/channels"

	"github.com/balerter/balerter/internal/channels/discord"
	"github.com/balerter/balerter/internal/channels/email"
	"github.com/balerter/balerter/internal/channels/notify"
	"github.com/balerter/balerter/internal/channels/slack"
	"github.com/balerter/balerter/internal/channels/syslog"
	"github.com/balerter/balerter/internal/channels/telegram"
	"github.com/balerter/balerter/internal/message"
	"go.uber.org/zap"
)

/*
Channels Manager is sending messages to channels

*/

type alertChannel interface {
	Name() string
	Send(*message.Message) error
}

// ChannelsManager represents the Alert manager struct
type ChannelsManager struct {
	logger   *zap.Logger
	channels map[string]alertChannel

	errs chan error
}

// New returns new Alert manager instance
func New(logger *zap.Logger) *ChannelsManager {
	m := &ChannelsManager{
		logger:   logger,
		channels: make(map[string]alertChannel),
		errs:     make(chan error),
	}

	go func() {
		for err := range m.errs {
			m.logger.Error("alert manager error", zap.Error(err))
		}
	}()

	return m
}

// Init the Alert manager
func (m *ChannelsManager) Init(cfg *channels.Channels) error {
	if cfg == nil {
		return nil
	}

	for idx := range cfg.Email {
		module, err := email.New(cfg.Email[idx], m.logger)
		if err != nil {
			return fmt.Errorf("error init email channel %s, %w", cfg.Email[idx].Name, err)
		}

		m.channels[module.Name()] = module
	}

	for idx := range cfg.Slack {
		module, err := slack.New(cfg.Slack[idx], m.logger)
		if err != nil {
			return fmt.Errorf("error init slack channel %s, %w", cfg.Slack[idx].Name, err)
		}

		m.channels[module.Name()] = module
	}

	for idx := range cfg.Telegram {
		module, err := telegram.New(cfg.Telegram[idx], m.logger)
		if err != nil {
			return fmt.Errorf("error init telegram channel %s, %w", cfg.Telegram[idx].Name, err)
		}

		m.channels[module.Name()] = module
	}

	for idx := range cfg.Syslog {
		module, err := syslog.New(cfg.Syslog[idx], m.logger)
		if err != nil {
			return fmt.Errorf("error init syslog channel %s, %w", cfg.Syslog[idx].Name, err)
		}

		m.channels[module.Name()] = module
	}

	for idx := range cfg.Notify {
		module, err := notify.New(cfg.Notify[idx], m.logger)
		if err != nil {
			return fmt.Errorf("error init syslog channel %s, %w", cfg.Notify[idx].Name, err)
		}

		m.channels[module.Name()] = module
	}

	for idx := range cfg.Discord {
		module, err := discord.New(cfg.Discord[idx], m.logger)
		if err != nil {
			return fmt.Errorf("error init discord channel %s, %w", cfg.Discord[idx].Name, err)
		}

		m.channels[module.Name()] = module
	}

	for idx := range cfg.Webhook {
		module, err := webhook.New(cfg.Webhook[idx], m.logger)
		if err != nil {
			return fmt.Errorf("error init webhook channel %s, %w", cfg.Webhook[idx].Name, err)
		}

		m.channels[module.Name()] = module
	}

	for idx := range cfg.Alertmanager {
		module, err := alertmanager.New(cfg.Alertmanager[idx], m.logger)
		if err != nil {
			return fmt.Errorf("error init alertmanager channel %s, %w", cfg.Alertmanager[idx].Name, err)
		}

		m.channels[module.Name()] = module
	}

	for idx := range cfg.AlertmanagerReceiver {
		module, err := alertmanagerreceiver.New(cfg.AlertmanagerReceiver[idx], m.logger)
		if err != nil {
			return fmt.Errorf("error init alertmanager_receiver channel %s, %w", cfg.AlertmanagerReceiver[idx].Name, err)
		}

		m.channels[module.Name()] = module
	}

	for idx := range cfg.TwilioVoice {
		module, err := twilio.New(cfg.TwilioVoice[idx], m.logger)
		if err != nil {
			return fmt.Errorf("error init twilio channel %s, %w", cfg.AlertmanagerReceiver[idx].Name, err)
		}

		m.channels[module.Name()] = module
	}

	return nil
}
