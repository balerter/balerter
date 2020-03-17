package syslog

import (
	"github.com/balerter/balerter/internal/config"
	"go.uber.org/zap"
	"log/syslog"
	"strings"
)

type Syslog struct {
	name   string
	logger *zap.Logger
	w      *syslog.Writer
}

func New(cfg config.ChannelSyslog, logger *zap.Logger) (*Syslog, error) {
	sl := &Syslog{
		name:   cfg.Name,
		logger: logger,
	}

	var err error

	if cfg.Priority == "" {
		cfg.Priority = "EMERG"
	}

	sl.w, err = syslog.Dial(cfg.Network, cfg.Address, parsePriority(cfg.Priority), cfg.Tag)
	if err != nil {
		return nil, err
	}

	return sl, nil
}

func (sl *Syslog) Name() string {
	return sl.name
}

func parsePriority(s string) syslog.Priority {
	if s == "" {
		return syslog.LOG_EMERG
	}

	parts := strings.Split(s, "|")

	priority := getSeverity(parts[0])

	if len(parts) == 2 {
		priority = priority | getFacility(parts[1])
	}

	return priority
}

func getFacility(s string) syslog.Priority {
	switch s {
	case "KERN":
		return syslog.LOG_KERN
	case "USER":
		return syslog.LOG_USER
	case "MAIL":
		return syslog.LOG_MAIL
	case "DAEMON":
		return syslog.LOG_DAEMON
	case "AUTH":
		return syslog.LOG_AUTH
	case "SYSLOG":
		return syslog.LOG_SYSLOG
	case "LPR":
		return syslog.LOG_LPR
	case "NEWS":
		return syslog.LOG_NEWS
	case "UUCP":
		return syslog.LOG_UUCP
	case "CRON":
		return syslog.LOG_CRON
	case "AUTHPRIV":
		return syslog.LOG_AUTHPRIV
	case "FTP":
		return syslog.LOG_FTP
	case "LOCAL0":
		return syslog.LOG_LOCAL0
	case "LOCAL1":
		return syslog.LOG_LOCAL1
	case "LOCAL2":
		return syslog.LOG_LOCAL2
	case "LOCAL3":
		return syslog.LOG_LOCAL3
	case "LOCAL4":
		return syslog.LOG_LOCAL4
	case "LOCAL5":
		return syslog.LOG_LOCAL5
	case "LOCAL6":
		return syslog.LOG_LOCAL6
	case "LOCAL7":
		return syslog.LOG_LOCAL7
	}

	panic("unexpected value")
}

func getSeverity(s string) syslog.Priority {
	switch s {
	case "EMERG":
		return syslog.LOG_EMERG
	case "ALERT":
		return syslog.LOG_ALERT
	case "CRIT":
		return syslog.LOG_CRIT
	case "ERR":
		return syslog.LOG_ERR
	case "WARNING":
		return syslog.LOG_WARNING
	case "NOTICE":
		return syslog.LOG_NOTICE
	case "INFO":
		return syslog.LOG_INFO
	case "DEBUG":
		return syslog.LOG_DEBUG
	}

	panic("unexpected value")
}
