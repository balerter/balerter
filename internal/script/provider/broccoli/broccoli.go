package broccoli

import (
	"io"
	"os"

	"aletheia.icu/broccoli/fs"
	"github.com/balerter/balerter/bindata"
	"github.com/balerter/balerter/internal/script/provider"
	"github.com/balerter/balerter/internal/script/script"
)

type Provider struct {
	broccoli *fs.Broccoli
}

func New() *Provider {
	p := &Provider{
		broccoli: bindata.Broccoli(),
	}
	return p
}

func (p *Provider) Open(f string) (io.ReadCloser, error) {
	return p.broccoli.Open(f)
}

func (p *Provider) Get() ([]*script.Script, error) {
	var ss []*script.Script

	err := p.broccoli.Walk("./modules", func(currentPath string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			s, err := provider.ReadScript(p, currentPath)
			if err != nil {
				return err
			}

			if s.Ignore {
				return nil
			}

			ss = append(ss, s)
		}

		return nil
	})

	return ss, err
}
