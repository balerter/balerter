package telegram

import (
	"fmt"
	"github.com/balerter/balerter/internal/channels/telegram/api"
	"github.com/balerter/balerter/internal/config/channels/telegram"
	"go.uber.org/zap"
)

type APIer interface {
	SendTextMessage(*api.TextMessage) error
	SendPhotoMessage(*api.PhotoMessage) error
}

type Telegram struct {
	name   string
	chatID int64
	logger *zap.Logger
	api    APIer
}

func New(cfg *telegram.Telegram, logger *zap.Logger) (*Telegram, error) {
	tg := &Telegram{
		name:   cfg.Name,
		chatID: cfg.ChatID,
		logger: logger,
	}

	var err error

	tg.api, err = api.New(cfg)

	if err != nil {
		return nil, fmt.Errorf("error connect to bot API, %w", err)
	}

	return tg, nil
}

func (tg *Telegram) Name() string {
	return tg.name
}
