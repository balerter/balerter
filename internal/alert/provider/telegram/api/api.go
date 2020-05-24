package api

import (
	"fmt"
	"github.com/balerter/balerter/internal/config"
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

func New(cfg *config.ChannelTelegram) (*API, error) {
	a := &API{
		endpoint: apiEndpoint + cfg.Token + "/",
	}

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
			return nil, fmt.Errorf("error create proxy, %w", err)
		}

		tr = &http.Transport{
			Proxy:       nil,
			DialContext: nil,
			Dial:        d.Dial,
		}
	}

	a.httpClient = &http.Client{
		Transport:     tr,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       cfg.Timeout,
	}

	if a.httpClient.Timeout == 0 {
		a.httpClient.Timeout = defaultHTTPClientTimeout
	}

	return a, nil
}
