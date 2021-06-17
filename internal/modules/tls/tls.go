package tls

import (
	"crypto/tls"
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"strings"
)

func ModuleName() string {
	return "tls"
}

func Methods() []string {
	return []string{
		"get",
	}
}

type TLS struct {
	dialFunc func(network, addr string, config *tls.Config) (*tls.Conn, error)
}

func New() *TLS {
	a := &TLS{
		dialFunc: tls.Dial,
	}

	return a
}

func (a *TLS) Name() string {
	return ModuleName()
}

func (a *TLS) GetLoader(_ *script.Script) lua.LGFunction {
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
