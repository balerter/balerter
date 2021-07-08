package webhook

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// AuthBasicConfig is basic auth config
type AuthBasicConfig struct {
	// Login for basic auth
	Login string `json:"login" yaml:"login" hcl:"login"`
	// Password for basic auth
	Password string `json:"password" yaml:"password" hcl:"password"`
}

// AuthBearerConfig is bearer auth config
type AuthBearerConfig struct {
	// Token for bearer auth
	Token string `json:"token" yaml:"token" hcl:"token"`
}

// AuthCustomConfig is custom auth config
type AuthCustomConfig struct {
	// TODO (negasus): remove Headers in favor of cfg.Headers option
	// Headers is request headers
	Headers map[string]string `json:"headers" yaml:"headers" hcl:"headers,optional"`
	// QueryParams is query params
	QueryParams map[string]string `json:"queryParams" yaml:"queryParams" hcl:"queryParams,optional"`
}

// AuthConfig for requests with auth
type AuthConfig struct {
	AuthBasicConfig
	AuthBearerConfig
	AuthCustomConfig

	// Type of the auth
	Type string `json:"type" yaml:"type" hcl:"type"`
}

// consts
const (
	// AuthTypeNone use for disable auth
	AuthTypeNone = "none"
	// AuthTypeBasic use for basic auth
	AuthTypeBasic = "basic"
	// AuthTypeBearer use for bearer token auth
	AuthTypeBearer = "bearer"
	// AuthTypeCustom use for custom auth
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

// PayloadConfig for POST requests
type PayloadConfig struct {
	QueryParams map[string]string `json:"queryParams" yaml:"queryParams" hcl:"queryParams,optional"`
	Body        string            `json:"body" yaml:"body" hcl:"body,optional"`
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

// Settings for webhook config
type Settings struct {
	URL     string            `json:"url" yaml:"url"  hcl:"url"`
	Method  string            `json:"method" yaml:"method" hcl:"method"`
	Auth    AuthConfig        `json:"auth" yaml:"auth" hcl:"auth,block"`
	Payload PayloadConfig     `json:"payload" yaml:"payload" hcl:"payload,block"`
	Timeout int               `json:"timeout" yaml:"timeout" hcl:"timeout,optional"`
	Headers map[string]string `json:"headers" yaml:"headers" hcl:"headers,optional"`
}

// Webhook configures notifications via webhook.
type Webhook struct {
	Name     string   `json:"name" yaml:"name" hcl:"name,label"`
	Settings Settings `json:"settings" yaml:"settings" hcl:"settings,block"`
}

// Validate checks the webhook configuration.
func (cfg Webhook) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	return cfg.Settings.Validate()
}

// Validate settings config
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
