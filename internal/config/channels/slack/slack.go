package slack

import (
	"fmt"
	"strings"
)

type Slack struct {
	Name    string `json:"name" yaml:"name"`
	Token   string `json:"token" yaml:"token"`
	Channel string `json:"channel" yaml:"channel"`
}

func (cfg Slack) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}
	if strings.TrimSpace(cfg.Token) == "" {
		return fmt.Errorf("token must be not empty")
	}
	if strings.TrimSpace(cfg.Channel) == "" {
		return fmt.Errorf("channel must be not empty")
	}

	return nil
}
