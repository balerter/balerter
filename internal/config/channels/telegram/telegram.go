package telegram

import (
	"fmt"
	"strings"
	"time"
)

type ProxyConfig struct {
	Address string          `json:"address" yaml:"address"`
	Auth    ProxyAuthConfig `json:"auth" yaml:"auth"`
}

type ProxyAuthConfig struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

type Telegram struct {
	Name    string        `json:"name" yaml:"name"`
	Token   string        `json:"token" yaml:"token"`
	ChatID  int64         `json:"chatId" yaml:"chatId"`
	Proxy   ProxyConfig   `json:"proxy" yaml:"proxy"`
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
}

func (cfg Telegram) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}
	if strings.TrimSpace(cfg.Token) == "" {
		return fmt.Errorf("token must be not empty")
	}
	if cfg.ChatID == 0 {
		return fmt.Errorf("chat id must be not empty")
	}
	if cfg.Timeout < 0 {
		return fmt.Errorf("timeout must be greater than 0")
	}

	return nil
}
