package folder

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/balerter/balerter/internal/config"
	"github.com/balerter/balerter/internal/script/script"
)

type Provider struct {
	path string
	mask string
}

func New(cfg *config.ScriptSourceFolder) *Provider {
	p := &Provider{
		path: cfg.Path,
		mask: cfg.Mask,
	}

	if p.mask == "" {
		p.mask = "*.lua"
	}

	return p
}

func (p *Provider) Get() ([]*script.Script, error) {
	ss := make([]*script.Script, 0)

	mask := path.Join(p.path, p.mask)

	matches, err := filepath.Glob(mask)
	if err != nil {
		return nil, err
	}

	for _, filename := range matches {
		body, err := ioutil.ReadFile(path.Join(filename))
		if err != nil {
			return nil, err
		}

		_, fn := path.Split(filename)

		s := script.New()
		s.Name = strings.TrimSuffix(fn, ".lua")
		s.Body = body

		if err := s.ParseMeta(); err != nil {
			return nil, err
		}

		if s.Ignore {
			continue
		}

		ss = append(ss, s)
	}

	return ss, nil
}
