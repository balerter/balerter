package s3

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/storages/upload/s3"
	"github.com/balerter/balerter/internal/modules"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"net/http"
)

// Provider represents S3 upload storage provider
type Provider struct {
	name     string
	region   string
	endpoint string
	key      string
	secret   string
	bucket   string
	logger   *zap.Logger
}

// Methods returns module methods
func Methods() []string {
	return []string{
		"uploadPNG",
	}
}

// ModuleName returns the module name
func ModuleName(name string) string {
	return "s3." + name
}

// New creates new S3 upload storage module
func New(cfg s3.S3, logger *zap.Logger) (*Provider, error) {
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

func (p *Provider) CoreApiHandler(method string, parts []string, params map[string]string, body []byte) (any, int, error) {
	return nil, http.StatusNotImplemented, fmt.Errorf("not implemented")
}

// Name returns the module name
func (p *Provider) Name() string {
	return p.name
}

// Stop the module
func (p *Provider) Stop() error {
	return nil
}

func (p *Provider) GetLoaderJS(_ modules.Job) require.ModuleLoader {
	return func(runtime *goja.Runtime, object *goja.Object) {

	}
}

// GetLoader returns the lua loader
func (p *Provider) GetLoader(_ modules.Job) lua.LGFunction {
	return p.loader
}

func (p *Provider) loader(luaState *lua.LState) int {
	var exports = map[string]lua.LGFunction{
		"uploadPNG": p.uploadPNG,
	}

	mod := luaState.SetFuncs(luaState.NewTable(), exports)

	luaState.Push(mod)
	return 1
}
