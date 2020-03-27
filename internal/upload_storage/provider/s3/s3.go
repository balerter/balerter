package s3

import (
	"github.com/balerter/balerter/internal/config"
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

type Provider struct {
	name     string
	region   string
	endpoint string
	key      string
	secret   string
	bucket   string
	logger   *zap.Logger
}

func Methods() []string {
	return []string{
		"uploadPNG",
	}
}

func ModuleName(name string) string {
	return "s3." + name
}

func New(cfg config.StorageUploadS3, logger *zap.Logger) (*Provider, error) {
	p := &Provider{
		name:     ModuleName(cfg.Name),
		region:   cfg.Region,
		endpoint: cfg.Endpoint,
		key:      cfg.Key,
		secret:   cfg.Secret,
		bucket:   cfg.Bucket,
		logger:   logger,
	}

	return p, nil
}

func (p *Provider) Name() string {
	return p.name
}

func (p *Provider) Stop() error {
	return nil
}

func (p *Provider) GetLoader(_ *script.Script) lua.LGFunction {
	return p.loader
}

func (p *Provider) loader(L *lua.LState) int {
	var exports = map[string]lua.LGFunction{
		"uploadPNG": p.uploadPNG,
	}

	mod := L.SetFuncs(L.NewTable(), exports)
	// register other stuff
	//L.SetField(mod, "name", lua.LString("value"))

	// returns the module
	L.Push(mod)
	return 1
}
