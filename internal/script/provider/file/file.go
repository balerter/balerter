package folder

import (
	"path"

	"github.com/balerter/balerter/internal/config"
	"github.com/balerter/balerter/internal/script/provider"
	"github.com/balerter/balerter/internal/script/script"
)

type Provider struct {
	filename      string
	disableIgnore bool
}

func New(cfg *config.ScriptSourceFile) *Provider {
	p := &Provider{
		filename:      cfg.Filename,
		disableIgnore: cfg.DisableIgnore,
	}

	return p
}

func (p *Provider) Get() ([]*script.Script, error) {
	ss := make([]*script.Script, 0)

	s, err := provider.ReadScript(provider.DefaultFs, path.Join(p.filename))
	if err != nil {
		return nil, err
	}

	if p.disableIgnore || !s.Ignore {
		ss = append(ss, s)
	}

	return ss, nil
}
