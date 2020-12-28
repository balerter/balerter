package api

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/channels/telegram"
	"golang.org/x/net/proxy"
	"net/http"
	"time"
)

const (
	apiEndpoint              = "https://api.telegram.org/bot"
	defaultHTTPClientTimeout = time.Second * 5

	methodSendMessage = "sendMessage"
	methodSendPhoto   = "sendPhoto"
)

type API struct {
	endpoint   string
	httpClient *http.Client
}

func New(cfg *telegram.Telegram) (*API, error) {
	a := &API{
		endpoint: apiEndpoint + cfg.Token + "/",
	}

	a.httpClient = &http.Client{
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       cfg.Timeout,
	}

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
			return nil, fmt.Errorf("error create proxy, %w", err)
		}

		a.httpClient.Transport = &http.Transport{
			Proxy:       nil,
			DialContext: nil,
			Dial:        d.Dial,
		}
	}

	if a.httpClient.Timeout == 0 {
		a.httpClient.Timeout = defaultHTTPClientTimeout
	}

	return a, nil
}
