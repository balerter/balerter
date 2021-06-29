package script

import (
	"crypto/sha1"
	"fmt"
	"github.com/robfig/cron/v3"
	"strings"
	"time"
)

const (
	// DefaultCronValue represents Default Cron Value
	DefaultCronValue = "0 * * * * *"
	// DefaultTimeout is the default timeout
	DefaultTimeout = time.Hour
)

// New creates new Script
func New() *Script {
	s := &Script{
		CronValue: DefaultCronValue,
		Timeout:   DefaultTimeout,
	}

	return s
}

var (
	// CronParser is default cron parser
	CronParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
)

// Script represents the Script
type Script struct {
	Name       string
	Body       []byte
	CronValue  string
	Timeout    time.Duration
	Ignore     bool
	Channels   []string
	IsTest     bool
	TestTarget string
}

// Hash returns the hash, based on script name and body
func (s *Script) Hash() string {
	return fmt.Sprintf("%x", sha1.Sum(append([]byte(s.Name+"@"), s.Body...)))
}

type parseMetaFunc func(l string, s *Script) error

var (
	metas = map[string]parseMetaFunc{
		"@cron":     parseMetaCron,
		"@ignore":   parseMetaIgnore,
		"@name":     parseMetaName,
		"@channels": parseMetaChannels,
		"@test":     parseMetaTest,
		"@timeout":  parseMetaTimeout,
	}
)

// ParseMeta parse script meta
func (s *Script) ParseMeta() error {
	lines := strings.Split(string(s.Body), "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		if !strings.HasPrefix(l, "--") {
			return nil
		}

		l = strings.TrimSpace(l[2:])

		for prefix, f := range metas {
			if strings.HasPrefix(l, prefix) {
				l = l[len(prefix):]
				if err := f(strings.TrimSpace(l), s); err != nil {
					return err
				}
				break
			}
		}
	}

	return nil
}

func parseMetaCron(l string, s *Script) error {
	_, err := CronParser.Parse(l)
	if err != nil {
		return fmt.Errorf("error parse cron value, %w", err)
	}

	s.CronValue = l

	return nil
}

func parseMetaTimeout(l string, s *Script) error {
	d, err := time.ParseDuration(strings.TrimSpace(l))
	if err != nil {
		return fmt.Errorf("error parse '%s' to time duration, %w", strings.TrimSpace(l), err)
	}

	s.Timeout = d

	return nil
}

func parseMetaIgnore(_ string, s *Script) error {
	s.Ignore = true

	return nil
}

func parseMetaName(l string, s *Script) error {
	if l == "" {
		return fmt.Errorf("name must be not empty")
	}

	s.Name = l

	return nil
}

func parseMetaTest(l string, s *Script) error {
	if l == "" {
		return fmt.Errorf("test must be not empty")
	}

	s.TestTarget = l
	s.IsTest = true

	return nil
}

func parseMetaChannels(l string, s *Script) error {
	if l == "" {
		return fmt.Errorf("channels must be not empty")
	}

	for _, channelName := range strings.Split(l, ",") {
		channelName = strings.TrimSpace(channelName)
		if channelName == "" {
			return fmt.Errorf("channel name must be not empty")
		}

		s.Channels = append(s.Channels, channelName)
	}

	return nil
}
