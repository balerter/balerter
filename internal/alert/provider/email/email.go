package email

import (
	"errors"
	"net/smtp"
	"os"
	"strings"

	"github.com/balerter/balerter/internal/config"
	"github.com/balerter/balerter/internal/types"
	"go.uber.org/zap"
)

// Email implements a Provider for email notifications.
type Email struct {
	conf     *config.ChannelEmail
	hostname string
	logger   *zap.Logger
	name     string
}

func New(cfg config.ChannelEmail, logger *zap.Logger) (*Email, error) {
	h, err := os.Hostname()
	// Use localhost if os.Hostname() fails
	if err != nil {
		h = "localhost.localdomain"
	}
	return &Email{conf: &cfg, logger: logger, hostname: h}, nil
}

func (e *Email) Name() string {
	return e.name
}

// auth resolves a string of authentication mechanisms.
func (e *Email) auth(mechs string) (smtp.Auth, error) {
	username := e.conf.AuthUsername

	// If no username is set, keep going without authentication.
	if e.conf.AuthUsername == "" {
		e.logger.Debug("email", zap.String("auth", "auth_username is not configured. Attempting to send email without authenticating"))
		return nil, nil
	}

	err := &types.MultiError{}
	for _, mech := range strings.Split(mechs, " ") {
		switch mech {
		case "CRAM-MD5":
			secret := string(e.conf.AuthSecret)
			if secret == "" {
				err.Add(errors.New("missing secret for CRAM-MD5 auth mechanism"))
				continue
			}
			return smtp.CRAMMD5Auth(username, secret), nil

		case "PLAIN":
			password := string(e.conf.AuthPassword)
			if password == "" {
				err.Add(errors.New("missing password for PLAIN auth mechanism"))
				continue
			}
			identity := e.conf.AuthIdentity

			return smtp.PlainAuth(identity, username, password, e.conf.ServerName), nil
		case "LOGIN":
			password := string(e.conf.AuthPassword)
			if password == "" {
				err.Add(errors.New("missing password for LOGIN auth mechanism"))
				continue
			}
			return LoginAuth(username, password), nil
		}
	}
	if err.Len() == 0 {
		err.Add(errors.New("unknown auth mechanism: " + mechs))
	}
	return nil, err
}

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

// Next is needed for for AUTH LOGIN.
func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch strings.ToLower(string(fromServer)) {
		case "username:":
			return []byte(a.username), nil
		case "password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("unexpected server challenge")
		}
	}
	return nil, nil
}
