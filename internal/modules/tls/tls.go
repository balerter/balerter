package tls

import (
	"crypto/tls"
	"strings"

	"github.com/balerter/balerter/internal/modules"

	lua "github.com/yuin/gopher-lua"
)

// ModuleName returns the module name
func ModuleName() string {
	return "tls"
}

// Methods returns module methods
func Methods() []string {
	return []string{
		"get",
	}
}

// TLS represents the TLS core module
type TLS struct {
	dialFunc func(network, addr string, config *tls.Config) (*tls.Conn, error)
}

// New creates new TLS core module
func New() *TLS {
	a := &TLS{
		dialFunc: tls.Dial,
	}

	return a
}

// Name returns the module name
func (a *TLS) Name() string {
	return ModuleName()
}

// GetLoader returns the lua module
func (a *TLS) GetLoader(_ modules.Job) lua.LGFunction {
	return func() lua.LGFunction {
		return func(luaState *lua.LState) int {
			var exports = map[string]lua.LGFunction{
				"get": a.get,
			}

			mod := luaState.SetFuncs(luaState.NewTable(), exports)

			luaState.Push(mod)
			return 1
		}
	}()
}

// Stop the module
func (a *TLS) Stop() error {
	return nil
}

func (a *TLS) get(luaState *lua.LState) int {
	hostnameVar := luaState.Get(1)
	if hostnameVar.Type() != lua.LTString {
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("parameter must be a string"))
		return 2
	}

	conn, err := a.dialFunc("tcp", hostnameVar.String()+":443", &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("error dial to host " + hostnameVar.String() + ", " + err.Error()))
		return 2
	}
	defer conn.Close()

	t := &lua.LTable{}
	for _, cert := range conn.ConnectionState().PeerCertificates {
		tc := &lua.LTable{}
		tc.RawSetString("issuer", lua.LString(cert.Issuer.String()))
		tc.RawSetString("expiry", lua.LNumber(cert.NotAfter.Unix()))
		tc.RawSetString("dns names", lua.LString(strings.Join(cert.DNSNames, ",")))
		tc.RawSetString("email addressed", lua.LString(strings.Join(cert.EmailAddresses, ",")))
		t.Append(tc)
	}

	luaState.Push(t)
	return 1
}
