package api

import (
	"context"
	"fmt"
	config "github.com/balerter/balerter/internal/config/channels/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"net"
	"time"
)

const (
	defaultHTTPClientTimeout = time.Second * 5
)

type API struct {
	api        *tgbotapi.BotAPI
}

func New(cfg *config.Telegram) (*API, error) {
	botAPI, err := tgbotapi.NewBotAPI(cfg.Token)

	if err != nil {
		return nil, fmt.Errorf("error create bot api, %w", err)
	}

	a := &API{
		api: botAPI,
	}

	return a, nil
}

type dialFunc func(network, addr string) (net.Conn, error)
type dialContextFunc func(ctx context.Context, network, addr string) (net.Conn, error)

func getDialContextFunc(d dialFunc) dialContextFunc {
	return func(_ context.Context, network, addr string) (net.Conn, error) {
		return d(network, addr)
	}
}
