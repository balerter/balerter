package telegram

import (
	"fmt"
	"github.com/balerter/balerter/internal/channels/telegram/api"
	"github.com/balerter/balerter/internal/config/channels/telegram"
	"go.uber.org/zap"
)

// APIer is an interface for Telegram API
type APIer interface {
	SendTextMessage(*api.TextMessage) error
	SendPhotoMessage(*api.PhotoMessage) error
}

// Telegram represent the channel of the type Telegram
type Telegram struct {
	name   string
	chatID int64
	logger *zap.Logger
	api    APIer
	ignore bool
}

// New creates new Telegram channel
func New(cfg telegram.Telegram, logger *zap.Logger) (*Telegram, error) {
	tg := &Telegram{
		name:   cfg.Name,
		chatID: cfg.ChatID,
		logger: logger,
		ignore: cfg.Ignore,
	}

	var err error

	tg.api, err = api.New(cfg)

	if err != nil {
		return nil, fmt.Errorf("error connect to bot API, %w", err)
	}

	return tg, nil
}

// Name returns the channel name
func (tg *Telegram) Name() string {
	return tg.name
}

func (tg *Telegram) Ignore() bool {
	return tg.ignore
}
