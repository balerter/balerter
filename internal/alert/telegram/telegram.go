package telegram

import (
	"github.com/balerter/balerter/internal/alert/alert"
	"github.com/balerter/balerter/internal/alert/message"
	"github.com/balerter/balerter/internal/config"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
	"golang.org/x/net/proxy"
	"log"
	"net/http"
	"time"
)

const (
	defaultHttpClientTimeout = time.Second * 10
)

type Telegram struct {
	name       string
	token      string
	chatID     int64
	logger     *zap.Logger
	httpClient *http.Client
	bot        *tgbotapi.BotAPI
}

func New(cfg config.ChannelTelegram, logger *zap.Logger) (*Telegram, error) {
	tg := &Telegram{
		name:   cfg.Name,
		chatID: cfg.ChatID,
		logger: logger,
	}

	var err error

	var tr *http.Transport

	if cfg.Proxy != nil {

		var proxyAuth *proxy.Auth

		if cfg.Proxy.Auth != nil {
			proxyAuth = &proxy.Auth{
				User:     cfg.Proxy.Auth.Username,
				Password: cfg.Proxy.Auth.Password,
			}
		}

		d, err := proxy.SOCKS5("tcp4", cfg.Proxy.Address, proxyAuth, nil)
		if err != nil {
			return nil, err
		}

		tr = &http.Transport{
			Proxy:       nil,
			DialContext: nil,
			Dial:        d.Dial,
		}
	}

	tg.httpClient = &http.Client{
		Transport:     tr,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       defaultHttpClientTimeout,
	}

	tg.bot, err = tgbotapi.NewBotAPIWithClient(cfg.Token, tg.httpClient)
	if err != nil {
		return nil, err
	}

	return tg, nil
}

func (tg *Telegram) Name() string {
	return tg.name
}

func (tg *Telegram) Send(level alert.Level, mes *message.Message) error {

	log.Printf("Authorized on account %s", tg.bot.Self.UserName)

	msg := tgbotapi.NewMessage(tg.chatID, mes.Text)
	_, err := tg.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
