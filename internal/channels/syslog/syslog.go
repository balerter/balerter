package syslog

import (
	"fmt"
	"io"
	"log/syslog"
	"strings"

	syslogCfg "github.com/balerter/balerter/internal/config/channels/syslog"

	"go.uber.org/zap"
)

// Syslog represents a channel of type Syslog
type Syslog struct {
	name   string
	logger *zap.Logger
	w      io.Writer
	ignore bool
}

var (
	defaultPriority = "EMERG"
)

// New creates new Syslog channel
func New(cfg syslogCfg.Syslog, logger *zap.Logger) (*Syslog, error) {
	sl := &Syslog{
		name:   cfg.Name,
		logger: logger,
		ignore: cfg.Ignore,
	}

	var err error

	if cfg.Priority == "" {
		cfg.Priority = defaultPriority
	}

	pr, errParsePriority := parsePriority(cfg.Priority)
	if errParsePriority != nil {
		return nil, fmt.Errorf("error parse priority, %w", errParsePriority)
	}

	sl.w, err = syslog.Dial(cfg.Network, cfg.Address, pr, cfg.Tag)
	if err != nil {
		return nil, err
	}

	return sl, nil
}

// Name returns the channel name
func (sl *Syslog) Name() string {
	return sl.name
}

func (sl *Syslog) Ignore() bool {
	return sl.ignore
}

func parsePriority(s string) (syslog.Priority, error) {
	if s == "" {
		return syslog.LOG_EMERG, nil
	}

	parts := strings.Split(s, "|")

	priority, errGetSeverity := getSeverity(parts[0])
	if errGetSeverity != nil {
		return 0, errGetSeverity
	}

	if len(parts) == 2 { //nolint:gomnd // parts count
		f, errGetFacility := getFacility(parts[1])
		if errGetFacility != nil {
			return 0, errGetFacility
		}
		priority |= f
	}

	return priority, nil
}

var facilities = map[string]syslog.Priority{
	"KERN":     syslog.LOG_KERN,
	"USER":     syslog.LOG_USER,
	"MAIL":     syslog.LOG_MAIL,
	"DAEMON":   syslog.LOG_DAEMON,
	"AUTH":     syslog.LOG_AUTH,
	"SYSLOG":   syslog.LOG_SYSLOG,
	"LPR":      syslog.LOG_LPR,
	"NEWS":     syslog.LOG_NEWS,
	"UUCP":     syslog.LOG_UUCP,
	"CRON":     syslog.LOG_CRON,
	"AUTHPRIV": syslog.LOG_AUTHPRIV,
	"FTP":      syslog.LOG_FTP,
	"LOCAL0":   syslog.LOG_LOCAL0,
	"LOCAL1":   syslog.LOG_LOCAL1,
	"LOCAL2":   syslog.LOG_LOCAL2,
	"LOCAL3":   syslog.LOG_LOCAL3,
	"LOCAL4":   syslog.LOG_LOCAL4,
	"LOCAL5":   syslog.LOG_LOCAL5,
	"LOCAL6":   syslog.LOG_LOCAL6,
	"LOCAL7":   syslog.LOG_LOCAL7,
}

func getFacility(s string) (syslog.Priority, error) {
	v, ok := facilities[s]
	if !ok {
		return 0, fmt.Errorf("unexpected facility value %s", s)
	}
	return v, nil
}

var severities = map[string]syslog.Priority{
	"EMERG":   syslog.LOG_EMERG,
	"ALERT":   syslog.LOG_ALERT,
	"CRIT":    syslog.LOG_CRIT,
	"ERR":     syslog.LOG_ERR,
	"WARNING": syslog.LOG_WARNING,
	"NOTICE":  syslog.LOG_NOTICE,
	"INFO":    syslog.LOG_INFO,
	"DEBUG":   syslog.LOG_DEBUG,
}

func getSeverity(s string) (syslog.Priority, error) {
	v, ok := severities[s]
	if !ok {
		return 0, fmt.Errorf("unexpected severity value %s", s)
	}
	return v, nil
}
