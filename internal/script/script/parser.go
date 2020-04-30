package script

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"go.uber.org/zap"
)

type Parser struct {
	logger *zap.Logger
}

func NewParser(logger *zap.Logger) *Parser {
	return &Parser{
		logger: logger,
	}
}

func (p *Parser) ParseFile(filepath string) (*Script, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	name := strings.TrimSuffix(file.Name(), ".lua")
	body, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	s := &Script{
		Name: name,
		Body: body,
	}
	if err := s.parseMeta(); err != nil {
		return nil, fmt.Errorf("script '%s': parse meta: %w", name, err)
	}

	if s.Schedule == nil {
		s.Schedule = DefaultSchedule
		p.logger.Sugar().Warnf("script '%s': cron is not defined, using default: %s", s.Name, s.Schedule)
	}

	return s, nil
}

var (
	metas = map[string]func(meta, l string, s *Script) error{
		"@cron":     parseMetaCron,
		"@yearly":   parseMetaPredefinedSchedule,
		"@annually": parseMetaPredefinedSchedule,
		"@monthly":  parseMetaPredefinedSchedule,
		"@weekly":   parseMetaPredefinedSchedule,
		"@daily":    parseMetaPredefinedSchedule,
		"@midnight": parseMetaPredefinedSchedule,
		"@hourly":   parseMetaPredefinedSchedule,
		"@every":    parseMetaInterval,
		"@ignore":   parseMetaIgnore,
		"@name":     parseMetaName,
		"@channels": parseMetaChannels,
		"@test":     parseMetaTest,
	}
)

func (s *Script) parseMeta() error {
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
				if err := f(prefix, strings.TrimSpace(l), s); err != nil {
					return err
				}
				//break
			}
		}
	}

	return nil
}

func parseMetaCron(_, l string, s *Script) error {
	if s.Schedule != nil {
		return newCronAlreadyDefinedError(s.Schedule.String())
	}

	sched, err := NewSchedule(l)
	if err != nil {
		return err
	}

	s.Schedule = sched

	return nil
}

func parseMetaPredefinedSchedule(meta, _ string, s *Script) error {
	if s.Schedule != nil {
		return newCronAlreadyDefinedError(s.Schedule.String())
	}

	sched, err := NewSchedule(meta)
	if err != nil {
		return err
	}

	s.Schedule = sched

	return nil
}

func parseMetaInterval(meta, l string, s *Script) error {
	if s.Schedule != nil {
		return newCronAlreadyDefinedError(s.Schedule.String())
	}

	sched, err := NewSchedule(meta + " " + l)
	if err != nil {
		return err
	}

	s.Schedule = sched
	return nil
}

func parseMetaIgnore(_, _ string, s *Script) error {
	s.Ignore = true

	return nil
}

func parseMetaName(_, l string, s *Script) error {
	if l == "" {
		return fmt.Errorf("name must be not empty")
	}

	s.Name = l

	return nil
}

func parseMetaTest(_, l string, s *Script) error {
	if l == "" {
		return fmt.Errorf("test must be not empty")
	}

	s.TestTarget = l
	s.IsTest = true

	return nil
}

func parseMetaChannels(_, l string, s *Script) error {
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
