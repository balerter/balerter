package script

import (
	"crypto/sha1"
	"fmt"
	"strings"
	"time"
)

const (
	DefaultInterval = time.Second * 60
	DefaultTimeout  = time.Hour
)

func New() *Script {
	s := &Script{
		Interval: DefaultInterval,
		Timeout:  DefaultTimeout,
	}

	return s
}

type Script struct {
	Name       string
	Body       []byte
	Interval   time.Duration
	Timeout    time.Duration
	Ignore     bool
	Channels   []string
	IsTest     bool
	TestTarget string
}

func (s *Script) Hash() string {
	return fmt.Sprintf("%x", sha1.Sum(append([]byte(s.Name+"@"), s.Body...)))
}

var (
	metas = map[string]func(l string, s *Script) error{
		"@interval": parseMetaInterval,
		"@ignore":   parseMetaIgnore,
		"@name":     parseMetaName,
		"@channels": parseMetaChannels,
		"@test":     parseMetaTest,
		"@timeout":  parseMetaTimeout,
	}
)

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

func parseMetaInterval(l string, s *Script) error {
	d, err := time.ParseDuration(strings.TrimSpace(l))
	if err != nil {
		return err
	}

	s.Interval = d

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
