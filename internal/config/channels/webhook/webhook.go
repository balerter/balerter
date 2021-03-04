package webhook

import (
	"fmt"
	"net/http"
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
	// TODO (negasus): remove Headers in favor of cfg.Headers option
	Headers     map[string]string `json:"headers" yaml:"headers"`
	QueryParams map[string]string `json:"queryParams" yaml:"queryParams"`
}

type AuthConfig struct {
	AuthBasicConfig
	AuthBearerConfig
	AuthCustomConfig

	Type string `json:"type" yaml:"type"`
}

// consts
const (
	AuthTypeNone   = "none"
	AuthTypeBasic  = "basic"
	AuthTypeBearer = "bearer"
	AuthTypeCustom = "custom"
)

// Validate checks the authorization configuration.
func (cfg AuthConfig) Validate() error {
	switch authType := strings.ToLower(strings.TrimSpace(cfg.Type)); authType {
	case "":
		cfg.Type = AuthTypeNone
		return nil
	case AuthTypeNone:
		return nil
	case AuthTypeBasic:
		if strings.TrimSpace(cfg.AuthBasicConfig.Login) == "" {
			return fmt.Errorf("login must be not empty")
		}
		if strings.TrimSpace(cfg.AuthBasicConfig.Password) == "" {
			return fmt.Errorf("password must be not empty")
		}
		return nil
	case AuthTypeBearer:
		if strings.TrimSpace(cfg.AuthBearerConfig.Token) == "" {
			return fmt.Errorf("token must be not empty")
		}
		return nil
	case AuthTypeCustom:
		return nil
	default:
		return fmt.Errorf("type must be set to none, basic, bearer or custom")
	}
}

type PayloadConfig struct {
	QueryParams map[string]string `json:"queryParams" yaml:"queryParams"`
	Body        string            `json:"body" yaml:"body"`
}

// Validate checks the payload configuration.
func (cfg PayloadConfig) Validate(method string) error {
	switch method {
	case http.MethodPost:
		if strings.TrimSpace(cfg.Body) == "" {
			return fmt.Errorf("body must be not empty")
		}
		return nil
	case http.MethodGet:
		if len(cfg.QueryParams) == 0 {
			return fmt.Errorf("queryParams must be not empty")
		}
		return nil

	default:
		return fmt.Errorf("method must be set to post or get")
	}
}

type Settings struct {
	URL     string            `json:"url" yaml:"url"`
	Method  string            `json:"method" yaml:"method"`
	Auth    AuthConfig        `json:"auth" yaml:"auth"`
	Payload PayloadConfig     `json:"payload" yaml:"payload"`
	Timeout int               `json:"timeout" yaml:"timeout"`
	Headers map[string]string `json:"headers" yaml:"headers"`
}

// Webhook configures notifications via webhook.
type Webhook struct {
	Name     string   `json:"name" yaml:"name"`
	Settings Settings `json:"settings" yaml:"settings"`
}

// Validate checks the webhook configuration.
func (cfg Webhook) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	return cfg.Settings.Validate()
}

func (cfg Settings) Validate() error {
	addr := strings.TrimSpace(cfg.URL)
	if addr == "" {
		return fmt.Errorf("url must be not empty")
	}
	if _, err := url.ParseRequestURI(addr); err != nil {
		return fmt.Errorf("error validate url: %w", err)
	}

	// TODO (negasus): mutation in the Validate method looks not good
	cfg.Method = strings.ToUpper(strings.TrimSpace(cfg.Method))
	if cfg.Method == "" {
		cfg.Method = http.MethodPost
	}

	if cfg.Timeout < 0 {
		return fmt.Errorf("timeout must be greater or equals 0")
	}

	if err := cfg.Auth.Validate(); err != nil {
		return fmt.Errorf("error validate auth: %w", err)
	}
	if err := cfg.Payload.Validate(cfg.Method); err != nil {
		return fmt.Errorf("error validate payload: %w", err)
	}
	return nil
}
