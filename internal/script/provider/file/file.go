package folder

import (
	"path"

	"github.com/balerter/balerter/internal/config"
	"github.com/balerter/balerter/internal/script/script"
)

type Provider struct {
	filename      string
	disableIgnore bool
	parser        *script.Parser
}

func New(parser *script.Parser, cfg config.ScriptSourceFile) *Provider {
	p := &Provider{
		filename:      cfg.Filename,
		disableIgnore: cfg.DisableIgnore,
		parser:        parser,
	}

	return p
}

func (p *Provider) Get() ([]*script.Script, error) {
	var ss []*script.Script
	s, err := p.parser.ParseFile(path.Join(p.filename))
	if err != nil {
		return nil, err
	}

	if p.disableIgnore || !s.Ignore {
		ss = append(ss, s)
	}

	return ss, nil
}
