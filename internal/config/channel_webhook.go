package config

import (
	"fmt"
	"net/url"
	"strings"
)

type AuthBasicConfig struct {
	Login    string `json:"login" yaml:"login"`
	Password string `json:"password" yaml:"password"`
}

type AuthBearerConfig struct {
	Token string `json:"token" yaml:"token"`
}

type AuthCustomConfig struct {
	Headers     map[string]string `json:"headers" yaml:"headers"`
	QueryParams map[string]string `json:"query_params" yaml:"query_params"`
}

type AuthConfig struct {
	AuthBasicConfig
	AuthBearerConfig
	AuthCustomConfig

	Type string `json:"type" yaml:"type"`
}

// Validate checks the authorization configuration.
func (cfg AuthConfig) Validate() error {
	switch authType := strings.ToLower(strings.TrimSpace(cfg.Type)); authType {
	case "", "none":
		return nil
	case "basic":
		if strings.TrimSpace(cfg.AuthBasicConfig.Login) == "" {
			return fmt.Errorf("login must be not empty")
		}
		if strings.TrimSpace(cfg.AuthBasicConfig.Password) == "" {
			return fmt.Errorf("password must be not empty")
		}
		return nil
	case "bearer":
		if strings.TrimSpace(cfg.AuthBearerConfig.Token) == "" {
			return fmt.Errorf("token must be not empty")
		}
		return nil
	case "custom":
		if len(cfg.AuthCustomConfig.Headers)+len(cfg.AuthCustomConfig.QueryParams) == 0 {
			return fmt.Errorf("headers and query_params must be not empty")
		}
		return nil
	default:
		return fmt.Errorf("type must be set to none, basic, bearer or custom")
	}
}

type PayloadConfig struct {
	QueryParams map[string]string `json:"query_params" yaml:"query_params"`
	Body        string            `json:"body" yaml:"body"`
}

// Validate checks the payload configuration.
func (cfg PayloadConfig) Validate(method string) error {
	switch method {
	case "post":
		if strings.TrimSpace(cfg.Body) == "" {
			return fmt.Errorf("body must be not empty")
		}
		return nil
	case "get":
		if len(cfg.QueryParams) == 0 {
			return fmt.Errorf("query_params must be not empty")
		}
		return nil

	default:
		return fmt.Errorf("method must be set to post or get")
	}
}

// ChannelWebhook configures notifications via webhook.
type ChannelWebhook struct {
	Name    string        `json:"name" yaml:"name"`
	URL     string        `json:"url" yaml:"url"`
	Method  string        `json:"method" yaml:"method"`
	Auth    AuthConfig    `json:"auth" yaml:"auth"`
	Payload PayloadConfig `json:"payload" yaml:"payload"`
}

// Validate checks the webhool configuration.
func (cfg ChannelWebhook) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	addr := strings.TrimSpace(cfg.URL)
	if addr == "" {
		return fmt.Errorf("url must be not empty")
	}
	if _, err := url.ParseRequestURI(addr); err != nil {
		return fmt.Errorf("error validate url: %w", err)
	}

	method := strings.ToLower(strings.TrimSpace(cfg.Method))
	if method == "" {
		method = "post"
	}
	if method != "post" && method != "get" {
		return fmt.Errorf("method must be set to post or get")
	}
	if err := cfg.Auth.Validate(); err != nil {
		return fmt.Errorf("error validate auth: %w", err)
	}
	if err := cfg.Payload.Validate(method); err != nil {
		return fmt.Errorf("error validate payload: %w", err)
	}
	return nil
}
