package folder

import (
	"github.com/balerter/balerter/internal/config/scripts/file"
	"github.com/balerter/balerter/internal/script/script"
	"io/ioutil"
	"path"
	"strings"
)

// Provider represents File script provider
type Provider struct {
	name     string
	filename string
}

// New creates new File script provider
func New(cfg file.File) *Provider {
	p := &Provider{
		name:     "file." + cfg.Name,
		filename: cfg.Filename,
	}

	return p
}

// Get scripts from the provider
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

	ss = append(ss, s)

	return ss, nil
}
