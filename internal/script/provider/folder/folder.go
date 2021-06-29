package folder

import (
	"github.com/balerter/balerter/internal/config/scripts/folder"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/balerter/balerter/internal/script/script"
)

// Provider represents Folder script provider
type Provider struct {
	name string
	path string
	mask string
}

// New creates new Folder script provider
func New(cfg folder.Folder) *Provider {
	p := &Provider{
		name: "folder." + cfg.Name,
		path: cfg.Path,
		mask: cfg.Mask,
	}

	if p.mask == "" {
		p.mask = "*.lua"
	}

	return p
}

// Get returns scripts from the provider
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
		s.Name = p.name + "." + strings.TrimSuffix(fn, ".lua")
		s.Body = body

		if err := s.ParseMeta(); err != nil {
			return nil, err
		}

		ss = append(ss, s)
	}

	return ss, nil
}
