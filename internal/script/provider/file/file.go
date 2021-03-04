package folder

import (
	"github.com/balerter/balerter/internal/config/scripts/file"
	"github.com/balerter/balerter/internal/script/script"
	"io/ioutil"
	"path"
	"strings"
)

type Provider struct {
	name     string
	filename string
}

func New(cfg file.File) *Provider {
	p := &Provider{
		name:     "file." + cfg.Name,
		filename: cfg.Filename,
	}

	return p
}

func (p *Provider) Get() ([]*script.Script, error) {
	ss := make([]*script.Script, 0)

	body, err := ioutil.ReadFile(path.Join(p.filename))
	if err != nil {
		return nil, err
	}

	_, fn := path.Split(p.filename)

	s := script.New()
	s.Name = p.name + "." + strings.TrimSuffix(fn, ".lua")
	s.Body = body

	if err := s.ParseMeta(); err != nil {
		return nil, err
	}

	if s.Ignore {
		return ss, nil
	}

	ss = append(ss, s)

	return ss, nil
}
