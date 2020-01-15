package folder

import (
	"github.com/balerter/balerter/internal/config"
	"github.com/balerter/balerter/internal/script/script"
	"io/ioutil"
	"path"
)

type Provider struct {
	path string
	mask string
}

func New(cfg config.ScriptSourceFolder) *Provider {
	p := &Provider{
		path: cfg.Path,
		mask: cfg.Mask,
	}

	return p
}

func (p *Provider) Get() ([]*script.Script, error) {
	ss := make([]*script.Script, 0)

	fi, err := ioutil.ReadDir(p.path)
	if err != nil {
		return nil, err
	}

	for _, fileInfo := range fi {

		// todo: check file mask

		body, err := ioutil.ReadFile(path.Join(p.path, fileInfo.Name()))
		if err != nil {
			return nil, err
		}

		s := &script.Script{
			Name:     fileInfo.Name(),
			Body:     body,
			Interval: script.DefaultInterval,
		}

		if err := s.ParseMeta(); err != nil {
			return nil, err
		}

		ss = append(ss, s)
	}

	return ss, nil
}
