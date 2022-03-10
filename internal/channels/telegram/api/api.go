package api

import (
	"context"
	"fmt"
	"github.com/balerter/balerter/internal/config/channels/telegram"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"time"
)

const (
	apiEndpoint              = "https://api.telegram.org/bot"
	defaultHTTPClientTimeout = time.Second * 5
)

// API represents the Telegram API
type API struct {
	endpoint   string
	httpClient *http.Client
}

// New creates new Telegram API
func New(cfg telegram.Telegram) (*API, error) {
	a := &API{
		endpoint: apiEndpoint + cfg.Token + "/",
	}

	a.httpClient = &http.Client{
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Millisecond * time.Duration(cfg.Timeout),
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
			DialContext: getDialContextFunc(d.Dial),
		}
	}

	if a.httpClient.Timeout == 0 {
		a.httpClient.Timeout = defaultHTTPClientTimeout
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
