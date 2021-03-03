package telegram

import (
	"fmt"
	"strings"
)

type ProxyConfig struct {
	Address string           `json:"address" yaml:"address" hcl:"address"`
	Auth    *ProxyAuthConfig `json:"auth" yaml:"auth" hcl:"auth,block"`
}

type ProxyAuthConfig struct {
	Username string `json:"username" yaml:"username" hcl:"username"`
	Password string `json:"password" yaml:"password" hcl:"password"`
}

type Telegram struct {
	Name    string       `json:"name" yaml:"name" hcl:"name,label"`
	Token   string       `json:"token" yaml:"token" hcl:"token"`
	ChatID  int64        `json:"chatId" yaml:"chatId" hcl:"chatId"`
	Proxy   *ProxyConfig `json:"proxy" yaml:"proxy" hcl:"proxy,block"`
	Timeout int          `json:"timeout" yaml:"timeout" hcl:"timeout,optional"`
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
