package manager

import (
	"fmt"
	"github.com/balerter/balerter/internal/alert/provider/alertmanager"
	alertmanagerreceiver "github.com/balerter/balerter/internal/alert/provider/alertmanager_receiver"
	"github.com/balerter/balerter/internal/alert/provider/webhook"

	"github.com/balerter/balerter/internal/alert/message"
	"github.com/balerter/balerter/internal/alert/provider/discord"
	"github.com/balerter/balerter/internal/alert/provider/email"
	"github.com/balerter/balerter/internal/alert/provider/notify"
	"github.com/balerter/balerter/internal/alert/provider/slack"
	"github.com/balerter/balerter/internal/alert/provider/syslog"
	"github.com/balerter/balerter/internal/alert/provider/telegram"
	"github.com/balerter/balerter/internal/config"
	coreStorage "github.com/balerter/balerter/internal/corestorage"
	"go.uber.org/zap"
)

type alertChannel interface {
	Name() string
	Send(*message.Message) error
}

type sendMessageFunc func(level, alertName, text string, channels, fields []string, image string)

// Manager represents the Alert manager struct
type Manager struct {
	logger   *zap.Logger
	channels map[string]alertChannel

	sendMessageFunc sendMessageFunc
	storage         coreStorage.CoreStorage
}

// New returns new Alert manager instance
func New(engine coreStorage.CoreStorage, logger *zap.Logger) *Manager {
	m := &Manager{
		logger:   logger,
		storage:  engine,
		channels: make(map[string]alertChannel),
	}

	m.sendMessageFunc = m.Send

	return m
}

// Init the Alert manager
func (m *Manager) Init(cfg *config.Channels) error { //nolint:gocyclo,funlen // Validate calls only once on application start
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

	return nil
}
