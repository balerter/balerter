package script

import (
	"crypto/sha1"
	"fmt"
	"strings"
	"time"
)

const (
	DefaultInterval time.Duration = time.Second * 30
)

type Script struct {
	Name     string
	Body     []byte
	Interval time.Duration
	Ignore   bool
}

func (s *Script) Hash() string {
	return fmt.Sprintf("%x", sha1.Sum(append([]byte(s.Name+"@"), s.Body...)))
}

var (
	metas = map[string]func(l string, s *Script) error{
		"@interval": parseMetaInterval,
		"@ignore":   parseMetaIgnore,
	}
)

func (s *Script) ParseMeta() error {
	lines := strings.Split(string(s.Body), "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if !strings.HasPrefix(l, "--") {
			return nil
		}

		l = strings.TrimSpace(l[2:])

		for prefix, f := range metas {
			if strings.HasPrefix(l, prefix) {
				l = l[len(prefix):]
				if err := f(l, s); err != nil {
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

func parseMetaIgnore(l string, s *Script) error {
	s.Ignore = true

	return nil
}
